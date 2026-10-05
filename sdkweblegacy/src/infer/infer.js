/**
 * @fileoverview Runs an i6 model (e.g. "fbt") in the browser and returns the result.
 *
 * Flow: download the model's sqlite + tflite, run the named query `model-${model}-input` (its result is
 * the `ctx`), let the caller build the model inputs from it, run the tflite model, and run the named query
 * `model-${model}-output` with the model outputs. Its result is returned as is. Nothing here is model specific.
 */

import { I6Error } from "../errors.js";
import { downloader } from "../internal/downloader/downloader.js";
import { fs } from "../internal/fs/fs.js";
import { sqlitez } from "../internal/sqlitez/sqlitez.js";

/**
 * @typedef {import("../services/services.js").Services} Services
 * @typedef {import("../internal/sqlitez/sqlitez.js").SqliteConn} SqliteConn
 */

/**
 * URLs of the model assets.
 * @typedef {Object} InferAssets
 * @property {string} sqlite - URL of the model's .sqlite file (must answer HEAD with an etag).
 * @property {string} tflite - URL of the model's .tflite file (must answer HEAD with an etag).
 * @property {string} [litertWasm] - URL of the folder with the LiteRT wasm files (ending in "/"). Defaults to the jsdelivr copy.
 */

/**
 * One model input tensor, built by the caller.
 * @typedef {Object} ModelInput
 * @property {ArrayLike<number>} data - Flat values, row-major.
 * @property {number[]} [shape] - Tensor shape. Defaults to the model's shape with -1 as 1.
 */

/**
 * Options of an inference run.
 * @typedef {Object} InferOptions
 * @property {string} model - Model name (e.g. "fbt"), also used to name the cached assets and the named queries.
 * @property {*} params - Model input, any JSON. It is the param of `model-${model}-input` and goes to `model-${model}-output`.
 * @property {InferAssets} assets - Where to download the model assets from.
 * @property {(ctx: *) => Record<string, ModelInput>|Promise<Record<string, ModelInput>>} buildInputs - Gets the result
 *   of `model-${model}-input` and returns the model inputs by exact model input name.
 */

/**
 * Downloaded model assets.
 * @typedef {Object} ModelAssets
 * @property {string} dbPath - OPFS path of the .sqlite file (open it with sqlitez).
 * @property {Uint8Array} modelBytes - Content of the .tflite file.
 */

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
   * @throws {I6Error} With code "invalid_argument" if opts is malformed, or "asset"/"model"/"query"
   *   if downloading, running the model or reading the sqlite fails.
   */
  async run(opts) {
    if (!opts || !opts.model) throw new I6Error("infer: opts.model is required", { code: "invalid_argument" })
    if (!opts.assets || !opts.assets.sqlite || !opts.assets.tflite) throw new I6Error("infer: opts.assets.sqlite and opts.assets.tflite are required", { code: "invalid_argument" })
    if (typeof opts.buildInputs !== "function") throw new I6Error("infer: opts.buildInputs must be a function", { code: "invalid_argument" })

    const { model, params, assets, buildInputs } = opts

    // 1. Get the model files (cached by etag) and the runtimes
    const { dbPath, modelBytes } = await downloadAssets(model, assets)
    const litert = await loadLiteRt(assets.litertWasm)
    const compiled = await loadModel(litert, modelBytes)

    // The sqlite is only needed while the stages below run
    const conn = await sqlitez.open({ path: dbPath })
    try {
      // 2. Read what the model needs from the sqlite, the caller turns it into inputs
      const ctx = await conn.namedJsonQuery({ name: `model-${model}-input`, param: params })
      console.log("infer: input query ->", Object.keys(ctx))
      const inputs = await buildInputs(ctx)

      // 3. Run the model
      const outputs = await runModel(litert, compiled, inputs)

      // 4. The sqlite turns the outputs into the result
      const outputsJson = Object.fromEntries(Object.entries(outputs).map(([name, values]) => [name, Array.from(values)]))
      const result = await conn.namedJsonQuery({ name: `model-${model}-output`, param: { params, ctx, outputs: outputsJson } })
      console.log("infer: output query ->", result)
      return result
    } finally {
      await conn.close()
      console.log("infer: sqlite closed")
    }
  }
}

/** Must match the @litertjs/core version in package.json. */
const LITERT_VERSION = "2.5.3"

/** Where the LiteRT wasm files are loaded from when the caller does not say. */
const LITERT_WASM_URL = `https://cdn.jsdelivr.net/npm/@litertjs/core@${LITERT_VERSION}/wasm/`

/** Last compiled model, with the bytes it was compiled from. */
let cachedModel = null

/**
 * Downloads (or reuses from the browser cache, checking the etag) the model assets and validates them.
 * @param {string} model - Model name, used to name the cached files.
 * @param {InferAssets} assets - Asset URLs.
 * @returns {Promise<ModelAssets>} The path of the sqlite and the bytes of the tflite.
 * @throws {I6Error} With code "asset" if a download fails or a file is not a valid sqlite/tflite.
 */
async function downloadAssets(model, assets) {
  const dbRef = await downloadFile(`infer-${model}-sqlite`, assets.sqlite)
  const modelRef = await downloadFile(`infer-${model}-tflite`, assets.tflite)

  // sqlitez.open checks the sqlite header
  const dbPath = await fs.resolve(dbRef)
  const modelBytes = await fs.read(modelRef)

  if (new TextDecoder().decode(modelBytes.subarray(4, 8)) !== "TFL3") {
    throw new I6Error(`infer: ${assets.tflite} is not a tflite file`, { code: "asset" })
  }
  console.log(`infer: assets are ready (sqlite at ${dbPath}, tflite ${modelBytes.length} bytes, TFL3)`)

  return { dbPath, modelBytes }
}

/**
 * Brings one remote file to the browser cache, downloading it only if its etag changed.
 * @param {string} name - Name of the cached file.
 * @param {string} url - URL of the file.
 * @returns {Promise<import("../internal/fs/fs.js").FileRef>} The cached file.
 * @throws {I6Error} With code "asset" if the download fails.
 */
async function downloadFile(name, url) {
  try {
    const before = await fs.get({ name })
    const ref = await downloader.update({ name, url })
    console.log(`infer: ${name}: ${before && before.version === ref.version ? "cache hit" : "downloaded"}`, ref)
    return ref
  } catch (cause) {
    throw new I6Error(`infer: could not get ${url}: ${cause.message}`, { code: "asset", cause })
  }
}

/**
 * Loads the bundled LiteRT runtime (wasm) once.
 * @param {string} [wasmUrl] - URL of the folder with the LiteRT wasm files. Defaults to the jsdelivr copy.
 * @returns {Promise<Object>} The LiteRT module.
 * @throws {I6Error} With code "model" if the runtime cannot be loaded.
 */
async function loadLiteRt(wasmUrl = LITERT_WASM_URL) {
  try {
    const litert = await import("@litertjs/core")
    const loading = litert.getGlobalLiteRtPromise()
    if (loading) {
      await loading
      console.log("infer: litert: already loaded")
    } else {
      await litert.loadLiteRt(wasmUrl)
      console.log(`infer: litert: loaded from ${wasmUrl}`)
    }
    return litert
  } catch (cause) {
    throw new I6Error(`infer: could not load LiteRT: ${cause.message}`, { code: "model", cause })
  }
}

/**
 * Compiles the tflite model, reusing the previous compilation when the bytes are the same.
 * @param {Object} litert - The LiteRT module.
 * @param {Uint8Array} modelBytes - Content of the .tflite file.
 * @returns {Promise<Object>} The compiled model.
 */
async function loadModel(litert, modelBytes) {
  const sameBytes = (a, b) => a.length === b.length && a.every((v, i) => v === b[i])

  if (cachedModel && sameBytes(cachedModel.bytes, modelBytes)) {
    console.log("infer: model: reusing the compiled model")
    return cachedModel.compiled
  }

  try {
    const compiled = await litert.loadAndCompile(modelBytes, { accelerator: "wasm" })
    if (cachedModel) cachedModel.compiled.delete()
    cachedModel = { bytes: modelBytes, compiled }

    const describe = (d) => ({ name: d.name, dtype: d.dtype, shape: Array.from(d.shape) })
    console.log("infer: model: compiled", {
      details: compiled.getInputDetails(),
      inputs: compiled.getInputDetails().map(describe),
      outputs: compiled.getOutputDetails().map(describe),
    })
    return compiled
  } catch (cause) {
    throw new I6Error(`infer: could not compile the model: ${cause.message}`, { code: "model", cause })
  }
}

/** Typed array of each dtype the runner accepts. */
const DTYPE_CTORS = { float32: Float32Array, int32: Int32Array, uint8: Uint8Array }

/**
 * Runs the model. The caller only gives data and shapes: every LiteRT tensor is created and freed here.
 * Every model input must be in `inputs` under its exact name.
 * @param {Object} litert - The LiteRT module, from loadLiteRt.
 * @param {Object} model - The compiled model, from loadModel.
 * @param {Record<string, ModelInput>} inputs - Input name -> data and shape.
 * @returns {Promise<Record<string, ArrayLike<number>>>} Output name -> values, for every output of the model.
 * @throws {I6Error} With code "invalid_argument" if an input is missing or its data does not fit its shape.
 */
async function runModel(litert, model, inputs) {
  const tensors = {}
  let raw = {}
  try {
    for (const detail of model.getInputDetails()) {
      const input = inputs?.[detail.name]
      if (!input) throw new I6Error(`infer: buildInputs did not return "${detail.name}"`, { code: "invalid_argument" })

      const shape = input.shape ?? Array.from(detail.shape, (d) => (d === -1 ? 1 : d))
      if (input.data.length !== shape.reduce((a, b) => a * b, 1)) {
        throw new I6Error(`infer: input "${detail.name}" has ${input.data.length} values, does not fit shape [${shape}]`, { code: "invalid_argument" })
      }
      tensors[detail.name] = new litert.Tensor(new DTYPE_CTORS[detail.dtype](input.data), shape)
    }

    raw = await model.run(tensors)
    const outputs = {}
    for (const [name, tensor] of Object.entries(raw)) outputs[name] = await tensor.data()
    return outputs
  } finally {
    Object.values(raw).forEach((t) => t.delete())
    Object.values(tensors).forEach((t) => t.delete())
  }
}

export { Infer, loadLiteRt, loadModel, runModel };
