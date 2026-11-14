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
 * Options of an inference update/preload.
 * @typedef {Object} InferUpdateOptions
 * @property {string} modelName - Model name (e.g. "relevance-fashion-fbt").
 * @property {Object<string, string>} partitions - Partitions of the model, name -> value (e.g. `{ version: "v1", store_id: "s1" }`).
 * @property {string} dataset - Dataset that holds the model files.
 * @property {string} [sqliteUrl] - Tests only: URL of the sqlite (must answer HEAD with an etag). Default: the model API.
 * @property {ServiceOptions} [serviceOptions] - Given to the API call that locates the model files.
 * @property {ServiceOptions} [serviceOpts] - Alias for `serviceOptions`.
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
 * @property {ServiceOptions} [serviceOptions] - Given to the API calls (the model API and the `recraw` event).
 * @property {ServiceOptions} [serviceOpts] - Alias for `serviceOptions`.
 */

/**
 * Options of `selectResults`.
 * @typedef {Object} SelectResultsOptions
 * @property {*} req - What the caller asked for (the cart).
 * @property {*} results - What `run` returned.
 * @property {string} recId - Id of the recommendation run, returned by `run`.
 * @property {string} sessionId - Id of the user session.
 * @property {string} userId - Id of the user.
 * @property {string} dataset - Dataset of the model.
 * @property {string} modelName - Model name.
 * @property {string} modelVersion - Model version.
 * @property {string} [storeId] - Store partition.
 * @property {ServiceOptions} [serviceOptions] - Given to the `recshow` event call.
 * @property {ServiceOptions} [serviceOpts] - Alias for `serviceOptions`.
 * @property {(req: *, results: *) => Promise<*>} [selectResults] - Custom selection function.
 */

/** The two files every model needs, by their role. */
const FILES = { sqlite: "sqlite.bin" }

/** Default limit when opts has none. */
const DEFAULT_LIMIT = 20

/** Seconds before checking for model file updates in the background (5 minutes). */
const UPDATE_INTERVAL = 5 * 60

/** Prints only when `window.I6_DEBUG` is true. */
const debug = debugLogger("infer")

/**
 * Infer client.
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
    /**
     * In-flight preload/update promises by cache name.
     * @private
     * @type {Map<string, Promise<FileRef>>}
     */
    this._preloads = new Map();
  }

  /**
   * Preloads or revalidates the model assets in the background.
   * Safe to call repeatedly as warmup; deduplicates in-flight calls and skips (noop) if checked within 5 minutes.
   * @param {InferUpdateOptions} opts - Preload options.
   * @returns {Promise<FileRef>} The cached sqlite file reference.
   */
  async update(opts) {
    validateUpdate(opts)

    const cacheName = `infer-${opts.modelName}-${partitionsKey(opts.partitions)}-sqlite`
    if (this._preloads.has(cacheName)) {
      debug("reusing in-flight preload", { cacheName })
      return this._preloads.get(cacheName)
    }

    const cached = await fs.get({ name: cacheName })
    if (cached) {
      const lastChecked = new Date(cached.checkedAt || cached.version).getTime()
      const age = (Date.now() - lastChecked) / 1000
      if (age < UPDATE_INTERVAL) {
        debug(`skipping update for ${cacheName} (checked ${Math.round(age)}s ago < ${UPDATE_INTERVAL}s)`)
        return cached
      }
    }

    const promise = (async () => {
      const urls = await this._assetUrls(opts)
      debug("model files located for preload", { sqlite: withoutQuery(urls.sqlite) })
      return downloadFile(cacheName, urls.sqlite, UPDATE_INTERVAL)
    })().finally(() => {
      this._preloads.delete(cacheName)
    })

    // Background detach safety
    promise.catch(() => {})

    this._preloads.set(cacheName, promise)
    return promise
  }

  /**
   * Runs the model in the browser and returns its recommendation result.
   * @param {InferOptions} opts - Inference options.
   * @returns {Promise<*>} The results of the model.
   */
  async run(opts) {
    validate(opts)

    const startedAt = performance.now()
    const elapsed = () => `[${((performance.now() - startedAt) / 1000).toFixed(2)}s]`
    const { modelName, params, sessionId, userId } = opts
    const ids = { sessionId, userId }
    const queryName = (kind) => `model-${modelName}-${kind}`

    // 1. Locate and download the model file (cached by etag)
    const cacheName = `infer-${modelName}-${partitionsKey(opts.partitions)}-sqlite`
    const cached = await fs.get({ name: cacheName })

    let sqliteRef
    if (cached) {
      debug(elapsed(), "1/4 using cached model file, revalidating in background", { version: cached.version })
      this.update(opts)
      sqliteRef = cached
    } else {
      debug(elapsed(), "1/4 no cached model file, waiting for download")
      sqliteRef = await this.update(opts)
    }

    const conn = await sqlitez.open({ path: await fs.resolve(sqliteRef) })
    try {
      // 2. The input query gives the model inputs: {name: {data, shape}}. The model compiles meanwhile.
      const [inputs, compiled] = await Promise.all([
        conn.namedStringQuery({ name: queryName("input"), params: [params] }),
        // loadModel({ ref: tfliteRef, wasmUrl: opts.litertWasm }),
      ])
      if (!inputs) throw new Error(`infer: ${queryName("input")} returned no model inputs`)
      debug(elapsed(), "2/4 input query", queryName("input"), "| bound as ?1:", params, "| returned (model inputs):", preview(inputs))

      // 3. Run the model
      debug(elapsed(), "3/4 calling the model on the server", { model_name: modelName, version: opts.partitions.version, model_partitions: opts.partitions })
      const modelStartedAt = performance.now()
      const outputs = await this.callModel(inputs, opts)
      const data = outputs.parse()
      debug(elapsed(), "3/4 model answered", { sec: +((performance.now() - modelStartedAt) / 1000).toFixed(2), ok: outputs?.ok, status: outputs?.status, error: outputs?.error, outputs: preview(data?.outputs) })

      let modelOutputsAndLimit = [JSON.parse(data.outputs), opts.limit ?? DEFAULT_LIMIT]
      debug(elapsed(), "4/4 output query", queryName("output"), "| bound as ?1 (model outputs) and ?2 (limit):", preview(modelOutputsAndLimit))
      // 4. The sqlite turns the outputs into the result
      const result = await conn.namedJsonQuery({ name: queryName("output"), params: modelOutputsAndLimit })
      debug(elapsed(), "output query returned (after namedJsonQuery)", result)
      // The rec_id ties this run to the recshow event that selectResults sends later
      const recId = crypto.randomUUID()
      const results = { ...result, rec_id: recId }
      const serviceOptions = opts.serviceOptions ?? opts.serviceOpts
      const eventOpts = { dataset: opts.dataset, modelName, modelVersion: opts.partitions.version, storeId: opts.partitions.store_id, serviceOptions }
      this.#sendRecEvent(eventOpts, ids, recId, "recraw", results)
        .then(
          () => debug(elapsed(), "recraw event sent", { rec_id: recId, results: results.results?.length }),
          (err) => debug(elapsed(), "recraw event failed", err)
        )
      const totalSec = +((performance.now() - startedAt) / 1000).toFixed(2)
      debug(elapsed(), "run completed", { rec_id: recId, totalSec })
      return results
    } finally {
      await conn.close()
    }
  }

  async callModel(inputs, opts) {
    const { dataset, modelName, partitions } = opts
    return await this._services.apis().tfzModelPredict({
      params: {
        dataset: dataset,
        model_name: modelName,
        version: partitions.version
      },
      payload: {
        model_partitions: partitions,
        inputs: inputs
      }
    }, opts.serviceOptions ?? opts.serviceOpts)
  }

  /**
   * Asks `selectResults` which of the recommended skus are shown and sends their sku ids as the `recshow` event.
   * @param {SelectResultsOptions} opts - Select options.
   * @returns {Promise<*>} What `selectResults` returned, as is (the default sets the `rating` of each sku; see `defaultSelectResults`).
   */
  async selectResults(opts) {
    const selectResults = opts?.selectResults ?? defaultSelectResults
    if (typeof selectResults !== "function") {
      throw new TypeError("infer: opts.selectResults must be a function")
    }

    const { req, results, recId, sessionId, userId } = opts

    if (!recId) throw new TypeError("infer: opts.recId is required")
    if (!sessionId || !userId) {
      throw new TypeError("infer: opts needs sessionId and userId")
    }
    const ids = { sessionId, userId }

    const selected = await selectResults(req, results)
    if (!Array.isArray(selected?.results)) {
      throw new TypeError("infer: selectResults must return { results: [...] }")
    }

    // The event carries only the sku ids
    const shown = {
      skus: selected.results.map((item, i) => {
        const skuId = item?.sku_id
        if (!skuId) throw new TypeError(`infer: selectResults returned an item with no sku_id at results[${i}]`)
        return { sku_id: skuId }
      })
    }

    debug("selectResults picked", { rec_id: recId, sku_ids: shown.skus.map((s) => s.sku_id) })
    this.#sendRecEvent(opts, ids, recId, "recshow", shown)
      .then(
        () => debug("recshow event sent", { rec_id: recId }),
        (err) => debug("recshow event failed", err)
      )

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
   */
  async #sendRecEvent(opts, ids, recId, kind, payload) {
    await this._services.apis().ingestzEventRelevanceFashionEvent({
      params: { dataset: opts.dataset },
      payload: {
        id: crypto.randomUUID(),
        store_id: opts.storeId,
        session_id: ids.sessionId,
        user_id: ids.userId,
        created_at: new Date().toISOString().split('.')[0] + 'Z',
        [kind]: {
          rec_id: recId,
          model_id: opts.modelName,
          model_version: opts.modelVersion,
          payload: JSON.stringify(payload),
        },
      },
    }, opts.serviceOptions ?? opts.serviceOpts)
  }

  /**
   * Finds where the model file is: the signed URLs of the model API, unless opts has them (tests).
   * @private
   * @param {InferOptions} opts - Inference options.
   * @returns {Promise<{sqlite: string}>} Download URLs of the files.
   */
  async _assetUrls(opts) {
    if (opts.sqliteUrl) return { sqlite: opts.sqliteUrl }

    const { modelName, partitions, dataset } = opts
    const resp = await this._services.apis().dszModelGet({
      params: { dataset, model_name: modelName },
      payload: { model_version: opts.partitions.version, model_partitions: partitions },
    }, opts.serviceOptions ?? opts.serviceOpts)

    const data = resp.parse()
    const byName = Object.fromEntries((data?.urls ?? []).map((file) => [file.name, file.url]))
    const missing = Object.values(FILES).find((name) => !byName[name])
    if (missing) {
      throw new Error(`infer: the model api gave no "${missing}" file: ${resp.error?.message ?? JSON.stringify(Object.keys(byName))}`)
    }
    return { sqlite: byName[FILES.sqlite] }
  }
}

/**
 * @param {string|undefined} url - A (signed) URL.
 * @returns {string|undefined} The URL without its query string, so signatures never reach the logs.
 */
function withoutQuery(url) {
  return typeof url === "string" ? url.split("?")[0] : url
}

/**
 * @param {*} value - Anything to log.
 * @param {number} [max] - Max characters kept.
 * @returns {string} The JSON of the value, cut with its full length when longer than `max` (model tensors are huge).
 */
function preview(value, max = 600) {
  const text = typeof value === "string" ? value : JSON.stringify(value)
  if (typeof text !== "string" || text.length <= max) return text
  return `${text.slice(0, max)}... (${text.length} chars in total)`
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
 * @param {InferUpdateOptions} opts - Options given to `Infer.update`.
 */
function validateUpdate(opts) {
  for (const key of ["modelName", "dataset"]) {
    if (!opts || !opts[key]) throw new TypeError(`infer: opts.${key} is required`)
  }

  if (!isPartitions(opts.partitions)) {
    throw new TypeError("infer: opts.partitions must be an object of string values")
  }
}

/**
 * @param {InferOptions} opts - Options given to `Infer.run`.
 */
function validate(opts) {
  for (const key of ["modelName", "dataset", "sessionId", "userId"]) {
    if (!opts || !opts[key]) throw new TypeError(`infer: opts.${key} is required`)
  }

  if (!isPartitions(opts.partitions)) {
    throw new TypeError("infer: opts.partitions must be an object of string values")
  }
  if (opts.params === undefined) throw new TypeError("infer: opts.params is required")
  if (opts.limit !== undefined && !(Number.isInteger(opts.limit) && opts.limit > 0)) {
    throw new TypeError("infer: opts.limit must be a positive integer")
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
