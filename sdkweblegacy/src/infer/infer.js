/**
 * @fileoverview Runs an i6 model (e.g. "relevance-fashion-fbt") in the browser and returns the result.
 *
 * Flow: get the model's sqlite + tflite (URLs from the model signed-get API, or from the caller in tests),
 * run the named query `model-${modelName}-input` (its result is the model inputs, `{name: {data, shape}}`),
 * run the tflite model, and run the named query `model-${modelName}-output` with the model outputs.
 * Its result is returned as is. Nothing here is model specific.
 *
 * Set `window.I6_DEBUG = true` to see what happens in the console, and to get the model tensors in the "i6:infer-tensors" event.
 */

import { I6Error } from "../errors.js";
import { downloadFile, isPartitions, partitionsKey } from "../internal/assets/assets.js";
import { fs } from "../internal/fs/fs.js";
import { sqlitez } from "../internal/sqlitez/sqlitez.js";
import { describeModel, loadModel, runModel } from "./litert.js";
import { debugLogger } from "../internal/debug/debug.js";

/**
 * @typedef {import("../services/services.js").Services} Services
 * @typedef {import("../internal/fs/fs.js").FileRef} FileRef
 * @typedef {import("../invoker/invoker.js").ServiceOptions} ServiceOptions
 */

/**
 * Options of an inference run.
 * @typedef {Object} InferOptions
 * @property {string} modelName - Model name (e.g. "relevance-fashion-fbt"); names the cached assets and the named queries.
 * @property {Object<string, string>} partitions - Partitions of the model, name -> value (e.g. `{ version: "v1", store_id: "s1" }`). Any partitions.
 * @property {string} dataset - Dataset that holds the model files.
 * @property {*} params - Model input, any JSON (e.g. `{ skus: [...] }`). Bound as `?1` of the input query.
 * @property {string} sessionId - Id of the user session; sent in the recommendation events.
 * @property {string} userId - Id of the user; sent in the recommendation events.
 * @property {number} [limit] - Max number of results asked from the output query (its `?2`). Default 20.
 * @property {string} [sqliteUrl] - Tests only, together with `tfliteUrl`: URL of the sqlite (must answer HEAD with an etag). Default: the model API.
 * @property {string} [tfliteUrl] - Tests only, together with `sqliteUrl`: URL of the tflite.
 * @property {string} [litertWasm] - URL of the folder with the LiteRT wasm files (ending in "/"). Defaults to the jsdelivr copy.
 * @property {ServiceOptions} [serviceOpts] - Given to the API calls (the model API and the `recraw` event).
 */

/**
 * Options of `Infer.selectResults`.
 * @typedef {Object} SelectResultsOptions
 * @property {string} dataset - Dataset that receives the `recraw` and `recshow` events.
 * @property {string} modelName - Model that made the recommendation (e.g. "relevance-fashion-fbt").
 * @property {string} modelVersion - Version of that model.
 * @property {string} storeId - Store of the event.
 * @property {string} sessionId - Id of the user session; sent in the recommendation events.
 * @property {string} userId - Id of the user; sent in the recommendation events.
 * @property {*} req - What the caller asked for; given as is to `selectResults`.
 * @property {*} results - What `run` returned (its `rec_id` is the id of the recommendation); given as is to `selectResults`.
 * @property {string} recId - Id of the recommendation, the `rec_id` that `run` returned. Same as in the `recraw` event that `run` sent; sent in the `recshow` event.
 * @property {(req: *, results: *) => Promise<*>} [selectResults] - Picks which of the recommended skus are shown: returns
 *   `{ results: [{ sku_id, rating?, enriched? }] }`. Only `sku_id` is sent (in the `recshow` event); `enriched` is any data of
 *   your own for the sku (e.g. name, price, image), returned to you as is. Default: `defaultSelectResults`.
 * @property {ServiceOptions} [serviceOpts] - Given to the API call of the `recshow` event.
 */

/** Default for `InferOptions.limit`. */
const DEFAULT_LIMIT = 20

/** Names of the model files in the model API response. */
const FILES = { sqlite: "sqlite.bin", tflite: "tflite.bin" }

/** Prints only when `window.I6_DEBUG` is true. */
const debug = debugLogger("infer")

/**
 * Inference client.
 */
class Infer {
  /**
   * Constructs an Infer client instance.
   * @param {Services} services - The services manager instance.
   */
  constructor(services) {
    /**
     * @private
     * @type {Services}
     */
    this._services = services;
  }

  /**
   * Runs the model for the given params and returns what the output query gives.
   * @param {InferOptions} opts - Inference options.
   * @returns {Promise<*>} The `ret` of the named query `model-${modelName}-output`, plus a generated `rec_id` (uuid). Sends the `recraw` event with that id.
   * @throws {I6Error} With code "invalid_argument" if opts is malformed, has no `sessionId` and `userId` or the model inputs do not fit the model,
   *   "asset" if a download fails, "model" if loading or running the model fails, "query" if the sqlite fails.
   */
  async run(opts) {
    validate(opts)
    const { modelName, params, sessionId, userId } = opts
    const ids = { sessionId, userId }
    const queryName = (kind) => `model-${modelName}-${kind}`

    // 1. Get the model files
    const urls = await this._assetUrls(opts)
    const cacheName = (kind) => `infer-${modelName}-${partitionsKey(opts.partitions)}-${kind}`
    const [sqliteRef, tfliteRef] = await Promise.all([
      downloadFile(cacheName("sqlite"), urls.sqlite),
      downloadFile(cacheName("tflite"), urls.tflite),
    ])

    const conn = await sqlitez.open({ path: await fs.resolve(sqliteRef) })
    try {
      // 2. The input query gives the model inputs: {name: {data, shape}}. The model compiles meanwhile.
      const [inputs, compiled] = await Promise.all([
        conn.namedJsonQuery({ name: queryName("input"), params: [params] }),
        loadModel({ ref: tfliteRef, wasmUrl: opts.litertWasm }),
      ])
      if (!inputs || typeof inputs !== "object") throw new I6Error(`infer: ${queryName("input")} returned no model inputs`, { code: "query" })
      debug("input query returned", Object.keys(inputs), "model", describeModel(compiled))

      // 3. Run the model
      const outputs = await runModel(compiled, inputs)
      debug("output returned (before namedJsonQuery)", outputs)
      announceTensors({ modelName, partitions: opts.partitions, model: describeModel(compiled), inputs, outputs })

      // 4. The sqlite turns the outputs into the result
      const result = await conn.namedJsonQuery({ name: queryName("output"), params: outputQueryParams(outputs, opts) })
      debug("output query returned (after namedJsonQuery)", result)
      // The rec_id ties this run to the recshow event that selectResults sends later
      const recId = crypto.randomUUID()
      const results = { ...result, rec_id: recId }
      const eventOpts = { dataset: opts.dataset, modelName, modelVersion: opts.partitions.version, storeId: opts.partitions.store_id, serviceOpts: opts.serviceOpts }
      await this.#sendRecEvent(eventOpts, ids, recId, "recraw", results)
      return results
    } finally {
      await conn.close()
    }
  }

  /**
   * Asks `selectResults` which of the recommended skus are shown and sends their sku ids as the `recshow` event.
   * @param {SelectResultsOptions} opts - Select options.
   * @returns {Promise<*>} What `selectResults` returned, as is (the default sets the `rating` of each sku; see `defaultSelectResults`).
   * @throws {I6Error} With code "invalid_argument" if `selectResults` is not a function, `recId` is missing, `opts` has no `sessionId` and
   *   `userId`, or `selectResults` returned an item with no `sku_id` or a sku not in the results; "event" if an event is not accepted.
   */
  async selectResults(opts) {
    const selectResults = opts?.selectResults ?? defaultSelectResults
    if (typeof selectResults !== "function") {
      throw new I6Error("infer: opts.selectResults must be a function", { code: "invalid_argument" })
    }

    const { req, results, recId, sessionId, userId } = opts

    if (!recId) throw new I6Error("infer: opts.recId is required", { code: "invalid_argument" })
    if (!sessionId || !userId) {
      throw new I6Error("infer: opts needs sessionId and userId", { code: "invalid_argument" })
    }
    const ids = { sessionId, userId }

    const selected = await selectResults(req, results)
    if (!Array.isArray(selected?.results)) {
      throw new I6Error("infer: selectResults must return { results: [...] }", { code: "invalid_argument" })
    }

    // The event carries only the sku ids
    const shown = {
      skus: selected.results.map((item, i) => {
        const skuId = item?.sku_id
        if (!skuId) throw new I6Error(`infer: selectResults returned an item with no sku_id at results[${i}]`, { code: "invalid_argument" })
        return { sku_id: skuId }
      })
    }

    await this.#sendRecEvent(opts, ids, recId, "recshow", shown)

    return selected
  }

  /**
   * Sends one recommendation event.
   * @private
   * @param {{dataset: string, modelName: string, modelVersion: string, storeId: string, serviceOpts?: ServiceOptions}} opts - Event options.
   * @param {{sessionId: string, userId: string}} ids - Session and user of the sdk.
   * @param {string} recId - Id of the recommendation, the same in the recraw and recshow events.
   * @param {"recraw"|"recshow"} kind - Which event.
   * @param {*} payload - What the event carries, any JSON.
   * @returns {Promise<void>}
   * @throws {I6Error} With code "event" if the event is not accepted.
   */
  async #sendRecEvent(opts, ids, recId, kind, payload) {
    await this._services.apis().ingestzStreamRelevanceFashionEvent({
      params: { dataset: opts.dataset },
      payload: {
        data: [{
          id: crypto.randomUUID(),
          store_id: opts.storeId,
          session_id: ids.sessionId,
          user_id: ids.userId,
          dataset: opts.dataset,
          created_at: new Date().toISOString().split('.')[0] + 'Z',
          [kind]: {
            rec_id: recId,
            model_id: opts.modelName,
            model_version: opts.modelVersion,
            payload: JSON.stringify(payload),
          },
        }],
      },
    }, opts.serviceOpts)
  }

  /**
   * Finds where the model files are: the signed URLs of the model API, unless opts has both URLs (tests).
   * @private
   * @param {InferOptions} opts - Inference options.
   * @returns {Promise<{sqlite: string, tflite: string}>} Download URLs of the files.
   * @throws {I6Error} With code "asset" if the API gave no URL for a file.
   */
  async _assetUrls(opts) {
    if (opts.sqliteUrl) return { sqlite: opts.sqliteUrl, tflite: opts.tfliteUrl }

    const { modelName, partitions, dataset } = opts
    const resp = await this._services.apis().dszModelGet({
      params: { dataset, model_name: modelName },
      payload: { model_version: opts.partitions.version, model_partitions: partitions },
    }, opts.serviceOpts)

    const byName = Object.fromEntries((resp.payload?.urls ?? []).map((file) => [file.name, file.url]))
    const missing = Object.values(FILES).find((name) => !byName[name])
    if (missing) {
      throw new I6Error(`infer: the model api gave no "${missing}" file: ${resp.error?.message ?? JSON.stringify(Object.keys(byName))}`, { code: "asset", status: resp.error?.status, cause: resp.error })
    }
    return { sqlite: byName[FILES.sqlite], tflite: byName[FILES.tflite] }
  }
}

/**
 * With `window.I6_DEBUG` set, tells the page what went into and came out of the model, as the "i6:infer-tensors" event
 * (`detail`: modelName, partitions, model, inputs, outputs). The demo shows it on screen. Outputs are copied to plain arrays.
 * @param {{modelName: string, partitions: Object<string, string>, model: Object, inputs: Object, outputs: Object<string, ArrayLike<number>>}} tensors
 */
function announceTensors({ outputs, ...rest }) {
  if (!globalThis.I6_DEBUG || typeof globalThis.dispatchEvent !== "function") return
  const plain = Object.fromEntries(Object.entries(outputs).map(([name, values]) => [name, Array.from(values)]))
  globalThis.dispatchEvent(new CustomEvent("i6:infer-tensors", { detail: { ...rest, outputs: plain } }))
}

/**
 * @param {InferOptions} opts - Options given to `Infer.run`.
 * @throws {I6Error} With code "invalid_argument" if a required option is missing.
 */
function validate(opts) {
  for (const key of ["modelName", "dataset", "sessionId", "userId"]) {
    if (!opts || !opts[key]) throw new I6Error(`infer: opts.${key} is required`, { code: "invalid_argument" })
  }

  if (!isPartitions(opts.partitions)) {
    throw new I6Error("infer: opts.partitions must be an object of string values", { code: "invalid_argument" })
  }
  if (opts.params === undefined) throw new I6Error("infer: opts.params is required", { code: "invalid_argument" })
  if (opts.limit !== undefined && !(Number.isInteger(opts.limit) && opts.limit > 0)) {
    throw new I6Error("infer: opts.limit must be a positive integer", { code: "invalid_argument" })
  }
  if (!opts.sqliteUrl !== !opts.tfliteUrl) {
    throw new I6Error("infer: opts.sqliteUrl and opts.tfliteUrl go together", { code: "invalid_argument" })
  }
}

/**
 * The values bound to the output query: `?1` the model outputs as `{name: {values: [...]}}` (plain arrays, which JSON
 * can hold), `?2` the limit.
 * @param {Record<string, ArrayLike<number>>} outputs - Model outputs (typed arrays).
 * @param {InferOptions} opts - Inference options.
 * @returns {Array<*>} The params of `namedJsonQuery`.
 */
function outputQueryParams(outputs, opts) {
  const named = Object.entries(outputs).map(([name, values]) => [name, { values: Array.from(values) }])
  return [Object.fromEntries(named), opts.limit ?? DEFAULT_LIMIT]
}

/**
 * Default of `SelectResultsOptions.selectResults`: one sku per result, skipping the skus already in the cart, with no
 * duplicates.
 *
 * Goes through `results.results` in the order the model returned them and, for each result, takes the first of its skus
 * that is not in `req.input.skus` (the cart) and was not already picked, with the `rating` of that result. A result with
 * no such sku gives nothing. All picks are returned (no limit; slice them if you show fewer). No `enriched` is set: pass
 * your own `selectResults` to attach your catalog data or to choose differently.
 *
 * @example
 * // cart: [{ sku_id: "1" }]
 * // results: [{ skus: [{ sku_id: "1" }, { sku_id: "2" }], rating: 0.9 }, { skus: [{ sku_id: "2" }, { sku_id: "3" }], rating: 0.7 }, { skus: [{ sku_id: "3" }], rating: 0.5 }]
 * // returns: { results: [{ sku_id: "2", rating: 0.9 }, { sku_id: "3", rating: 0.7 }] }
 * @param {*} req - What the caller asked for; its `input.skus` (the cart) is required.
 * @param {*} results - What `run` returned: `{ results: [{ skus: [{ sku_id }], rating }] }`.
 * @returns {Promise<{results: Array<{sku_id: string, rating: number}>}>}
 */
async function defaultSelectResults(req, results) {
  const taken = new Set(req.input.skus.map((sku) => sku.sku_id))
  const shown = []

  for (const result of results.results) {
    for (const { sku_id } of result.skus) {
      if (!taken.has(sku_id)) {
        taken.add(sku_id)
        shown.push({ sku_id, rating: result.rating })
        break
      }
    }
  }

  return {
    results: shown
  }
}

export { Infer };
