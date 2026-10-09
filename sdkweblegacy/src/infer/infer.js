/**
 * @fileoverview Runs an i6 model (e.g. "relevance-fashion-fbt") in the browser and returns the result.
 *
 * Flow: get the model's sqlite + tflite (URLs from the model signed-get API, or from the caller in tests),
 * run the named query `model-${modelName}-input` (its result is the model inputs, `{name: {data, shape}}`),
 * run the tflite model, and run the named query `model-${modelName}-output` with the model outputs.
 * Its result is returned as is. Nothing here is model specific.
 *
 * Set `window.I6_DEBUG = true` to see what happens in the console.
 */

import { I6Error } from "../errors.js";
import { downloadFile, isPartitions, partitionsKey } from "../internal/assets/assets.js";
import { fs } from "../internal/fs/fs.js";
import { sqlitez } from "../internal/sqlitez/sqlitez.js";
import { describeModel, loadModel, runModel } from "./litert.js";

/**
 * @typedef {import("../services/services.js").Services} Services
 * @typedef {import("../internal/fs/fs.js").FileRef} FileRef
 */

/**
 * Options of an inference run.
 * @typedef {Object} InferOptions
 * @property {string} modelName - Model name (e.g. "relevance-fashion-fbt"); names the cached assets and the named queries.
 * @property {string} modelVersion - Model version (e.g. "v1").
 * @property {Object<string, string>} partitions - Partitions of the model, name -> value (e.g. `{ store_id: "s1" }`). Any partitions.
 * @property {string} dataset - Dataset that holds the model files.
 * @property {*} params - Model input, any JSON (e.g. `{ cart_skus: [...] }`). Bound as `?1` of the input query.
 * @property {number} [limit] - Max number of results asked from the output query (its `?2`). Default 20.
 * @property {string} [sqliteUrl] - Tests only, together with `tfliteUrl`: URL of the sqlite (must answer HEAD with an etag). Default: the model API.
 * @property {string} [tfliteUrl] - Tests only, together with `sqliteUrl`: URL of the tflite.
 * @property {string} [litertWasm] - URL of the folder with the LiteRT wasm files (ending in "/"). Defaults to the jsdelivr copy.
 */

/** Default for `InferOptions.limit`. */
const DEFAULT_LIMIT = 20

/** Names of the model files in the model API response. */
const FILES = { sqlite: "sqlite.bin", tflite: "tflite.bin" }

/** Prints only when `window.I6_DEBUG` is true. */
const debug = (...args) => { if (globalThis.I6_DEBUG) console.log("infer:", ...args) }

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
   * @returns {Promise<*>} The `ret` of the named query `model-${modelName}-output`.
   * @throws {I6Error} With code "invalid_argument" if opts is malformed or the model inputs do not fit the model,
   *   "asset" if a download fails, "model" if loading or running the model fails, "query" if the sqlite fails.
   */
  async run(opts) {
    validate(opts)
    const { modelName, params } = opts
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

      // 4. The sqlite turns the outputs into the result
      const result = await conn.namedJsonQuery({ name: queryName("output"), params: outputQueryParams(outputs, opts) })
      debug("output query returned", result)
      return result
    } finally {
      await conn.close()
    }
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
    })

    const byName = Object.fromEntries((resp.payload?.urls ?? []).map((file) => [file.name, file.url]))
    const missing = Object.values(FILES).find((name) => !byName[name])
    if (missing) {
      throw new I6Error(`infer: the model api gave no "${missing}" file: ${JSON.stringify(resp.error ?? Object.keys(byName))}`, { code: "asset" })
    }
    return { sqlite: byName[FILES.sqlite], tflite: byName[FILES.tflite] }
  }
}

/**
 * @param {InferOptions} opts - Options given to `Infer.run`.
 * @throws {I6Error} With code "invalid_argument" if a required option is missing.
 */
function validate(opts) {
  for (const key of ["modelName", "dataset"]) {
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

export { Infer };
