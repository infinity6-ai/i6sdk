/**
 * @fileoverview Runs an i6 model (e.g. "fbt") in the browser and returns the result.
 *
 * Flow: download the model's sqlite + tflite, run the named query `model-${model}-input` (its result is
 * the `ctx`; its `inputs` are the model inputs), run the tflite model, and run the named query
 * `model-${model}-output` with the model outputs. Its result is returned as is. Nothing here is model specific.
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
 * URLs of the model assets.
 * @typedef {Object} InferAssets
 * @property {string} sqlite - URL of the model's .sqlite file (must answer HEAD with an etag).
 * @property {string} tflite - URL of the model's .tflite file (must answer HEAD with an etag).
 * @property {string} [litertWasm] - URL of the folder with the LiteRT wasm files (ending in "/"). Defaults to the jsdelivr copy.
 */

/**
 * Options of an inference run.
 * @typedef {Object} InferOptions
 * @property {string} model - Model name (e.g. "fbt"), also used to name the cached assets and the named queries.
 * @property {*} params - Model input, any JSON. It is the param of `model-${model}-input` and goes to `model-${model}-output`.
 * @property {InferAssets} assets - Where to download the model assets from.
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
   * @returns {Promise<*>} The `ret` of the named query `model-${model}-output`.
   * @throws {I6Error} With code "invalid_argument" if opts is malformed or the model inputs do not fit the model,
   *   "asset" if a download fails, "model" if loading or running the model fails, "query" if the sqlite fails.
   */
  async run(opts) {
    validate(opts)
    const { model, params, assets } = opts

    //se não vier nada no assets
    let bla = await this._services.apis().dszSignedUrl({
      params: {
        dataset: "dectest"
      },
      payload: {
        paths: [
          { 
            method: "GET", 
            path: "i6tb://results.relevance-fashion-fbt-sqlite/store_id=s1/version=v1/db.sqlite"
          },
          { 
            method: "GET", 
            path: "i6tb://results.relevance-fashion-fbt-model/store_id=s1/version=v1/model.tflite"
          }
        ]
      },
    })

    //Se vier assets, usa o proprio assets
    //Enviar sem assets é o padrão
    //Asset é usado pra teste 

    // 1. Get the model files (cached by etag) and load the model
    const sqliteRef = await downloadFile(`infer-${model}-sqlite`, assets.sqlite)
    const tfliteRef = await downloadFile(`infer-${model}-tflite`, assets.tflite)
    const compiled = await loadModel({ ref: tfliteRef, wasmUrl: assets.litertWasm })
    debug("model ready", describeModel(compiled))

    // The sqlite is only needed while the stages below run
    const conn = await sqlitez.open({ path: await fs.resolve(sqliteRef) })
    try {
      // 2. Read what the model needs from the sqlite: the model inputs (`inputs`) and anything the output query needs
      const ctx = await conn.namedJsonQuery({ name: `model-${model}-input`, param: params })
      if (!ctx?.inputs) throw new I6Error(`infer: model-${model}-input returned no "inputs"`, { code: "query" })
      debug("input query returned", Object.keys(ctx))

      // 3. Run the model
      const outputs = await runModel(compiled, ctx.inputs)

      // 4. The sqlite turns the outputs into the result
      const result = await conn.namedJsonQuery({ name: `model-${model}-output`, param: { params, ctx, outputs: asArrays(outputs) } })
      debug("output query returned", result)
      return result
    } finally {
      await conn.close()
    }
  }
}

/**
 * @param {InferOptions} opts - Options given to `Infer.run`.
 * @throws {I6Error} With code "invalid_argument" if a required option is missing.
 */
function validate(opts) {
  if (!opts || !opts.model) throw new I6Error("infer: opts.model is required", { code: "invalid_argument" })
  if (!opts.assets || !opts.assets.sqlite || !opts.assets.tflite) throw new I6Error("infer: opts.assets.sqlite and opts.assets.tflite are required", { code: "invalid_argument" })
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
 * @param {string} name - Name of the cached file.
 * @param {string} url - URL of the file.
 * @returns {Promise<FileRef>} The cached file.
 * @throws {I6Error} With code "asset" if the download fails.
 */
async function downloadFile(name, url) {
  try {
    const ref = await downloader.update({ name, url })
    debug(`${name} is ready`, ref)
    return ref
  } catch (cause) {
    throw new I6Error(`infer: could not get ${url}: ${cause.message}`, { code: "asset", cause })
  }
}

export { Infer };
