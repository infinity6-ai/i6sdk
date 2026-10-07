/**
 * @fileoverview Runs an i6 model (e.g. "relevance-fashion-fbt") in the browser and returns the result.
 *
 * Flow: get the model's sqlite + tflite (the URLs come from the dataset signed-url API, unless the caller passes
 * them, which is meant for tests), run the named query `model-${modelName}-input` (its result is
 * the `ctx`; its `inputs` are the model inputs), run the tflite model, and run the named query
 * `model-${modelName}-output` with the model outputs. Its result is returned as is. Nothing here is model specific.
 *
 * Set `window.I6_DEBUG = true` to see what happens in the console.
 */

import { I6Error } from "../errors.js";
import { downloader } from "../internal/downloader/downloader.js";
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
 * @property {string} storeId - Store the model was trained for (e.g. "s1").
 * @property {string} dataset - Dataset that holds the model files.
 * @property {*} params - Model input, any JSON (e.g. `{ sku_ids: [...] }`). It is the param of `model-${modelName}-input` and goes to `model-${modelName}-output`.
 * @property {string} [sqliteUrl] - Only for tests: URL of the model's .sqlite (must answer HEAD with an etag). An "i6tb://" path is signed first.
 *   Default: discovered with the API.
 * @property {string} [tfliteUrl] - Only for tests: URL of the model's .tflite (same rules as `sqliteUrl`). Default: discovered with the API.
 * @property {string} [litertWasm] - URL of the folder with the LiteRT wasm files (ending in "/"). Defaults to the jsdelivr copy.
 */

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
    const { modelName, modelVersion, storeId, params } = opts
    const cacheName = (kind) => `infer-${modelName}-${storeId}-${modelVersion}-${kind}`

    // 1. Get the model files (cached by etag) and load the model
    const [sqliteUrls, tfliteUrls] = await this._assetUrls(opts)
    const [sqliteRef, tfliteRef] = await Promise.all([
      downloadFile(cacheName("sqlite"), sqliteUrls),
      downloadFile(cacheName("tflite"), tfliteUrls),
    ])
    const compiled = await loadModel({ ref: tfliteRef, wasmUrl: opts.litertWasm })
    debug("model ready", describeModel(compiled))

    const conn = await sqlitez.open({ path: await fs.resolve(sqliteRef) })
    try {
      // 2. Read what the model needs from the sqlite: the model inputs (`inputs`) and anything the output query needs
      const ctx = await conn.namedJsonQuery({ name: `model-${modelName}-input`, param: params })
      if (!ctx?.inputs) throw new I6Error(`infer: model-${modelName}-input returned no "inputs"`, { code: "query" })
      debug("input query returned", Object.keys(ctx))

      // 3. Run the model
      const outputs = await runModel(compiled, ctx.inputs)

      // 4. The sqlite turns the outputs into the result
      const result = await conn.namedJsonQuery({ name: `model-${modelName}-output`, param: { params, ctx, outputs: asArrays(outputs) } })
      debug("output query returned", result)
      return result
    } finally {
      await conn.close()
    }
  }

  /**
   * Finds where the model files are. A URL given in opts is used as is (an "i6tb://" path is signed);
   * the others are the model's default paths in the dataset, signed with one API call.
   * @private
   * @param {InferOptions} opts - Inference options.
   * @returns {Promise<[string[], string[]]>} For the sqlite and the tflite: the URL to try first, then its fallbacks.
   * @throws {I6Error} With code "asset" if the API cannot sign them.
   */
  async _assetUrls(opts) {
    const { modelName, modelVersion, storeId, dataset } = opts
    const dir = (kind) => `i6tb://results.${modelName}-${kind}/store_id=${storeId}/version=${modelVersion}`

    const locations = [
      opts.sqliteUrl ?? `${dir("sqlite")}/db.sqlite`,
      opts.tfliteUrl ?? `${dir("model")}/model.tflite`,
    ]

    const paths = locations.filter(isI6tbPath)
    if (paths.length === 0) return locations.map((url) => [url])

    const resp = await this._services.apis().dszSignedUrl({
      params: { dataset },
      payload: { paths: paths.map((path) => ({ method: "GET", path })) },
    })
    
    const signed = resp.payload?.urls
    if (!resp.ok || signed?.length !== paths.length) {
      throw new I6Error(`infer: could not get signed urls for ${paths.join(", ")}: ${JSON.stringify(resp.error ?? resp.status)}`, { code: "asset" })
    }

    let next = 0
    return locations.map((location) => {
      if (!isI6tbPath(location)) return [location]
      const { url, fallbacks } = signed[next++]
      return [url, ...(fallbacks ?? [])]
    })
  }
}

/** @param {string} location - A URL or an "i6tb://" dataset path. */
const isI6tbPath = (location) => location.startsWith("i6tb://")

/**
 * @param {InferOptions} opts - Options given to `Infer.run`.
 * @throws {I6Error} With code "invalid_argument" if a required option is missing.
 */
function validate(opts) {
  for (const key of ["modelName", "modelVersion", "storeId", "dataset"]) {
    if (!opts || !opts[key]) throw new I6Error(`infer: opts.${key} is required`, { code: "invalid_argument" })
  }
}

/**
 * @param {Record<string, ArrayLike<number>>} outputs - Model outputs (typed arrays).
 * @returns {Record<string, number[]>} The same outputs as plain arrays, which JSON can hold.
 */
function asArrays(outputs) {
  return Object.fromEntries(Object.entries(outputs).map(([name, values]) => [name, Array.from(values)]))
}

/**
 * Brings one remote file to the browser cache, downloading it only if its etag changed.
 * Tries the urls in order; the next one is used only when the previous answered 404.
 * @param {string} name - Name of the cached file.
 * @param {string[]} urls - URL of the file, then its fallbacks.
 * @returns {Promise<FileRef>} The cached file.
 * @throws {I6Error} With code "asset" if the download fails.
 */
async function downloadFile(name, urls) {
  let failure
  for (const url of urls) {
    try {
      const ref = await downloader.update({ name, url })
      debug(`${name} is ready`, ref)
      return ref
    } catch (cause) {
      failure = new I6Error(`infer: could not get ${url}: ${cause.message}`, { code: "asset", cause })
      if (cause.status !== 404) break
    }
  }
  throw failure
}

export { Infer };
