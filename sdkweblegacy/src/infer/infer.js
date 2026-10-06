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
 * One model input, as nested rows (e.g. `[[1, 2], [3, 4]]`). The shape is read from the nesting and a missing
 * leading batch dimension of 1 is added.
 * @typedef {Object} ModelInput
 * @property {Array} data - Nested rows.
 * @property {number} [fill] - Value used to pad up to the model's shape. Too many rows are cut from the front (the last
 *   ones are kept). Without `fill`, data that does not match the model's shape is an error.
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
  constructor(services) { }

  /**
   * Runs the model for the given params and returns what the output query gives.
   * @param {InferOptions} opts - Inference options.
   * @returns {Promise<*>} The `ret` of the named query `model-${model}-output`.
   * @throws {I6Error} With code "invalid_argument" if opts is malformed or the model inputs do not fit the model,
   *   "asset" if a download fails, "model" if loading or running the model fails, "query" if the sqlite fails.
   */
  async run(opts) {
    if (!opts || !opts.model) throw new I6Error("infer: opts.model is required", { code: "invalid_argument" })
    if (!opts.assets || !opts.assets.sqlite || !opts.assets.tflite) throw new I6Error("infer: opts.assets.sqlite and opts.assets.tflite are required", { code: "invalid_argument" })

    const { model, params, assets } = opts

    // 1. Get the model files (cached by etag) and load the model
    const sqliteRef = await download(`infer-${model}-sqlite`, assets.sqlite)
    const tfliteRef = await download(`infer-${model}-tflite`, assets.tflite)
    const compiled = await loadModel({ ref: tfliteRef, wasmUrl: assets.litertWasm })

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
      const outputsJson = Object.fromEntries(Object.entries(outputs).map(([name, values]) => [name, Array.from(values)]))
      const result = await conn.namedJsonQuery({ name: `model-${model}-output`, param: { params, ctx, outputs: outputsJson } })
      debug("output query returned", result)
      return result
    } finally {
      await conn.close()
    }
  }
}

/** Must match the @litertjs/core version in package.json. */
const LITERT_VERSION = "2.5.3"

/** Where the LiteRT wasm files are loaded from when the caller does not say. */
const LITERT_WASM_URL = `https://cdn.jsdelivr.net/npm/@litertjs/core@${LITERT_VERSION}/wasm/`

/** Typed array of each dtype the runner accepts. */
const DTYPE_CTORS = { float32: Float32Array, int32: Int32Array, uint8: Uint8Array }

/**
 * The last loaded model: `key` is the name and version of its file, so it is reused until the file changes.
 * @type {{key: string, litert: Object, compiled: Object}|null}
 */
let loaded = null

/**
 * Brings one remote file to the browser cache, downloading it only if its etag changed.
 * @param {string} name - Name of the cached file.
 * @param {string} url - URL of the file.
 * @returns {Promise<FileRef>} The cached file.
 * @throws {I6Error} With code "asset" if the download fails.
 */
async function download(name, url) {
  try {
    const ref = await downloader.update({ name, url })
    debug(`${name} is ready`, ref)
    return ref
  } catch (cause) {
    throw new I6Error(`infer: could not get ${url}: ${cause.message}`, { code: "asset", cause })
  }
}

/**
 * Loads the LiteRT runtime (once) and the tflite model. Everything LiteRT stays inside this file. The file is only
 * read and compiled when its version changed since the last call.
 * @param {Object} opts
 * @param {FileRef} opts.ref - The downloaded .tflite file.
 * @param {string} [opts.wasmUrl] - URL of the folder with the LiteRT wasm files. Defaults to the jsdelivr copy.
 * @returns {Promise<Object>} The compiled model, to give to `runModel` and `describeModel`.
 * @throws {I6Error} With code "asset" if the file is not a tflite file, "model" if the runtime or the model cannot be loaded.
 */
async function loadModel(opts) {
  const key = `${opts.ref.name}@${opts.ref.version}`
  if (loaded?.key === key) return loaded.compiled

  let bytes
  try {
    bytes = await fs.read(opts.ref)
  } catch (cause) {
    throw new I6Error(`infer: could not read the model: ${cause.message}`, { code: "asset", cause })
  }
  if (new TextDecoder().decode(bytes.subarray(4, 8)) !== "TFL3") {
    throw new I6Error(`infer: ${opts.ref.name} is not a tflite file`, { code: "asset" })
  }

  try {
    const litert = await import("@litertjs/core")
    const loading = litert.getGlobalLiteRtPromise()
    if (loading) await loading
    else await litert.loadLiteRt(opts.wasmUrl ?? LITERT_WASM_URL)

    const compiled = await litert.loadAndCompile(bytes, { accelerator: "wasm" })
    if (loaded) loaded.compiled.delete()
    loaded = { key, litert, compiled }

    debug("model compiled", describeModel(compiled))
    return compiled
  } catch (cause) {
    throw new I6Error(`infer: could not load the model: ${cause.message}`, { code: "model", cause })
  }
}

/**
 * Names, dtypes and shapes of what the model takes and returns.
 * @param {Object} model - The compiled model, from `loadModel`.
 * @returns {{inputs: Object[], outputs: Object[]}}
 */
function describeModel(model) {
  const describe = (d) => ({ name: d.name, dtype: d.dtype, shape: Array.from(d.shape) })
  return { inputs: model.getInputDetails().map(describe), outputs: model.getOutputDetails().map(describe) }
}

/**
 * Flat values and shape of one input, fitted to the model's shape: missing batch dimension added, padded with `fill`.
 * @param {string} name - Input name, for error messages.
 * @param {ModelInput} input - What the caller gave.
 * @param {number[]} modelShape - Shape of the model input (-1 is free).
 * @returns {{flat: number[], shape: number[]}}
 * @throws {I6Error} With code "invalid_argument" if the data does not fit the model's shape.
 */
function toTensorData(name, input, modelShape) {
  const invalid = (msg) => new I6Error(`infer: input "${name}" ${msg}`, { code: "invalid_argument" })

  let data = input.data
  let rank = 0
  for (let a = data; Array.isArray(a); a = a[0]) rank++
  if (rank === modelShape.length - 1) data = [data]
  else if (rank !== modelShape.length) throw invalid(`has ${rank} dimensions, the model needs ${modelShape.length}`)

  const make = (dims) => (dims.length === 0 ? input.fill : Array.from({ length: dims[0] }, () => make(dims.slice(1))))
  const fit = (rows, dims) => {
    if (dims.length === 0) return rows
    const n = dims[0] === -1 ? rows.length : dims[0]
    if (rows.length !== n && input.fill === undefined) throw invalid(`has ${rows.length} values where the model needs ${n} (give "fill" to pad)`)
    const fitted = rows.length > n ? rows.slice(rows.length - n) : [...rows, ...Array.from({ length: n - rows.length }, () => make(dims.slice(1)))]
    return fitted.map((row) => fit(row, dims.slice(1)))
  }
  data = fit(data, modelShape)

  const shape = []
  for (let a = data; Array.isArray(a); a = a[0]) shape.push(a.length)
  return { flat: data.flat(Infinity), shape }
}

/**
 * Runs the model. The caller only gives data: every LiteRT tensor is created and freed here.
 * Every model input must be in `inputs` under its exact name.
 * @param {Object} model - The compiled model, from loadModel.
 * @param {Record<string, ModelInput>} inputs - Input name -> data.
 * @returns {Promise<Record<string, ArrayLike<number>>>} Output name -> values, for every output of the model.
 * @throws {I6Error} With code "invalid_argument" if an input is missing or its data does not fit the model, "model" if the run fails.
 */
async function runModel(model, inputs) {
  const tensors = {}
  let raw = {}
  try {
    for (const detail of model.getInputDetails()) {
      const input = inputs?.[detail.name]
      if (!input) throw new I6Error(`infer: input "${detail.name}" is missing`, { code: "invalid_argument" })

      const { flat, shape } = toTensorData(detail.name, input, Array.from(detail.shape))
      tensors[detail.name] = new loaded.litert.Tensor(new DTYPE_CTORS[detail.dtype](flat), shape)
    }

    try {
      raw = await model.run(tensors)
      const outputs = {}
      for (const [name, tensor] of Object.entries(raw)) outputs[name] = await tensor.data()
      return outputs
    } catch (cause) {
      throw new I6Error(`infer: could not run the model: ${cause.message}`, { code: "model", cause })
    }
  } finally {
    Object.values(raw).forEach((t) => t.delete())
    Object.values(tensors).forEach((t) => t.delete())
  }
}

export { Infer, loadModel, describeModel, runModel, toTensorData };
