/**
 * @fileoverview The LiteRT side of infer: loads the tflite model and runs it. Everything that touches the
 * LiteRT library (runtime, model, tensors) is in this file. The caller only gives data and gets values back.
 */

import { I6Error } from "../errors.js";
import { fs } from "../internal/fs/fs.js";

/**
 * @typedef {import("../internal/fs/fs.js").FileRef} FileRef
 */

/**
 * One model input: either flat `data` with its `shape` (what the sqlite gives, e.g. `{data: [1, 2, 3, 4], shape: [2, 2]}`),
 * or nested rows without `shape` (e.g. `[[1, 2], [3, 4]]`): the shape is then read from the nesting and a missing
 * leading batch dimension of 1 is added.
 * @typedef {Object} ModelInput
 * @property {Array} data - Flat values (with `shape`) or nested rows.
 * @property {number[]} [shape] - Shape of the flat `data`. It must be the model's shape (-1 is free), no padding is done.
 * @property {number} [fill] - Value used to pad up to the model's shape. Too many rows are cut from the front (the last
 *   ones are kept). Without `fill`, data that does not match the model's shape is an error.
 */

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
let loadedModel = null

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
  if (loadedModel?.key === key) return loadedModel.compiled

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
    if (loadedModel) loadedModel.compiled.delete()
    loadedModel = { key, litert, compiled }

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

  if (input.shape) {
    const shape = input.shape
    const flat = Array.isArray(input.data) ? input.data.flat(Infinity) : null
    if (!flat) throw invalid("has no data array")
    if (shape.length !== modelShape.length || shape.some((n, i) => modelShape[i] !== -1 && modelShape[i] !== n)) {
      throw invalid(`has shape [${shape}], the model needs [${modelShape}]`)
    }
    const size = shape.reduce((a, n) => a * n, 1)
    if (flat.length !== size) throw invalid(`has ${flat.length} values, its shape [${shape}] needs ${size}`)
    return { flat, shape: [...shape] }
  }

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
      tensors[detail.name] = new loadedModel.litert.Tensor(new DTYPE_CTORS[detail.dtype](flat), shape)
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

export { loadModel, describeModel, runModel, toTensorData };
