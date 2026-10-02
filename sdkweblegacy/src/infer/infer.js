/**
 * @fileoverview Runs an i6 model (e.g. "fbt") in the browser and returns the result.
 *
 * Flow: download the model's sqlite + tflite, read the features of the input products from
 * the sqlite, feed them to the tflite model, rank the model scores and map the best ones
 * back to product ids through the sqlite. Nothing is rendered, the caller uses the result.
 */

import { I6Error } from "../errors.js";
import { downloader } from "../internal/downloader/downloader.js";
import { fs } from "../internal/fs/fs.js";

/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * URLs of the model assets.
 * @typedef {Object} InferAssets
 * @property {string} sqlite - URL of the model's .sqlite file (must answer HEAD with an etag).
 * @property {string} tflite - URL of the model's .tflite file (must answer HEAD with an etag).
 * @property {string} [litertWasm] - URL of the folder with the LiteRT wasm files (ending in "/"). Defaults to the jsdelivr copy.
 */

/**
 * Parameters of the inference.
 * @typedef {Object} InferParams
 * @property {string[]} productIds - Ids of the products the recommendation is based on (e.g. the cart).
 */

/**
 * Options of an inference run.
 * @typedef {Object} InferOptions
 * @property {string} model - Model name (e.g. "fbt"), also used to name the cached assets.
 * @property {InferParams} params - Model input.
 * @property {InferAssets} assets - Where to download the model assets from.
 */

/**
 * Result of an inference run.
 * @typedef {Object} InferResult
 * @property {string[]} productIds - Recommended product ids, best first. Empty when none of the input products is known to the model.
 */

/**
 * Downloaded model assets.
 * @typedef {Object} ModelAssets
 * @property {Uint8Array} dbBytes - Content of the .sqlite file.
 * @property {Uint8Array} modelBytes - Content of the .tflite file.
 */

/**
 * Product data read from the sqlite for the input products.
 * @typedef {Object} ShopData
 * @property {Record<string, Object>} productFeaturesMap - Product id -> sqlite row (rec_in, rec_out, base_text, ...).
 * @property {number[]} storeMask - Per signature index: 1 if some product has that signature, 0 otherwise.
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
   * Runs the model for the given params and returns the recommended product ids.
   * @param {InferOptions} opts - Inference options.
   * @returns {Promise<InferResult>} The recommended product ids.
   * @throws {I6Error} With code "invalid_argument" if opts is malformed, or "asset"/"model"/"query"
   *   if downloading, running the model or reading the sqlite fails.
   */
  async run(opts) {
    const { model, params, assets } = opts

    // 1. Get the model files (cached by etag) and the runtimes
    const { dbBytes, modelBytes } = await downloadAssets(model, assets)
    const litert = await loadLiteRt(assets.litertWasm)
    const compiled = await loadModel(litert, modelBytes)

    // 2. Read the input products from the sqlite
    // The sqlinput and output now must be names of queries that will be run using sqlitez namedquery (they will be hardcoded, model-fbt-input and model-fbt-output, "fbt" is model name)
    const productIds = [...new Set(params.productIds)]
    const { productFeaturesMap, storeMask } = await fetchShopData(dbBytes, productIds)

    // 3. Run the model
    const tensors = buildInputTensors(litert, compiled, params.productIds, productFeaturesMap)
    if (!tensors) return { productIds: [] }
    const scores = await runModel(compiled, tensors)

    // 4. Rank the scores and map the best ones back to product ids
    const signatureIndices = rankSignatures(scores, storeMask, TOP_N)
    return { productIds: await fetchRecommendedProductIds(dbBytes, signatureIndices) }
  }
}

/** How many recommendations are kept. */
const TOP_N = 20

/** Must match the @litertjs/core version in package.json. */
const LITERT_VERSION = "2.5.3"

/** Where the LiteRT wasm files are loaded from when the caller does not say. */
const LITERT_WASM_URL = `https://cdn.jsdelivr.net/npm/@litertjs/core@${LITERT_VERSION}/wasm/`

/** Last compiled model, with the bytes it was compiled from. */
let cachedModel = null

/**
 * Downloads (or reads from the browser cache, checking the etag) the model assets and validates them.
 * @param {string} model - Model name, used to name the cached files.
 * @param {InferAssets} assets - Asset URLs.
 * @returns {Promise<ModelAssets>} The bytes of both files.
 * @throws {I6Error} With code "asset" if a download fails or a file is not a valid sqlite/tflite.
 */
async function downloadAssets(model, assets) {
  const dbBytes = await downloadFile(`infer-${model}-sqlite`, assets.sqlite)
  const modelBytes = await downloadFile(`infer-${model}-tflite`, assets.tflite)

  if (new TextDecoder().decode(dbBytes.subarray(0, 16)) !== "SQLite format 3\0") {
    throw new I6Error(`infer: ${assets.sqlite} is not a sqlite file`, { code: "asset" })
  }
  if (new TextDecoder().decode(modelBytes.subarray(4, 8)) !== "TFL3") {
    throw new I6Error(`infer: ${assets.tflite} is not a tflite file`, { code: "asset" })
  }
  console.log("infer: assets are valid (sqlite header, TFL3)")

  return { dbBytes, modelBytes }
}

/**
 * Brings one remote file to the browser cache (if its etag changed) and reads it.
 * @param {string} name - Name of the cached file.
 * @param {string} url - URL of the file.
 * @returns {Promise<Uint8Array>} The bytes of the file.
 * @throws {I6Error} With code "asset" if the download or the read fails.
 */
async function downloadFile(name, url) {
  try {
    const before = await fs.get({ name })
    const ref = await downloader.update({ name, url })
    const bytes = await fs.read(ref)
    console.log(`infer: ${name}: ${before && before.version === ref.version ? "cache hit" : "downloaded"}, ${bytes.length} bytes`, ref)
    return bytes
  } catch (cause) {
    throw new I6Error(`infer: could not get ${url}: ${cause.message}`, { code: "asset", cause })
  }
}

/**
 * Loads the bundled sql.js (wasm) once.
 * @returns {Promise<Object>} The initialized sql.js module (`SQL`).
 */
async function loadSqlJs() { }

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
      inputs: compiled.getInputDetails().map(describe),
      outputs: compiled.getOutputDetails().map(describe),
    })
    return compiled
  } catch (cause) {
    throw new I6Error(`infer: could not compile the model: ${cause.message}`, { code: "model", cause })
  }
}

/**
 * Reads from the sqlite the rows of the input products and the store mask.
 * @param {Uint8Array} dbBytes - Content of the .sqlite file.
 * @param {string[]} productIds - Unique ids of the input products.
 * @returns {Promise<ShopData>} The product rows by id and the mask of every signature index.
 */
async function fetchShopData(dbBytes, productIds) { }

/**
 * Builds the model input tensors: the feature vectors (`rec_in`) of the last items, one
 * masked slot to predict, then padding, plus the padding mask. The inputs are found by name
 * ("field" and "pad") and the shapes come from the model.
 * @param {Object} litert - The LiteRT module.
 * @param {Object} model - The compiled model.
 * @param {string[]} productIds - Input product ids, in order (may repeat).
 * @param {Record<string, Object>} productFeaturesMap - Product id -> sqlite row.
 * @returns {Record<string, Object>|null} Input name -> tensor, or null if no product has a usable feature vector.
 */
function buildInputTensors(litert, model, productIds, productFeaturesMap) { }

/**
 * Runs the model and frees the input and output tensors.
 * @param {Object} model - The compiled model.
 * @param {Record<string, Object>} tensors - Input name -> tensor, from buildInputTensors.
 * @returns {Promise<Float32Array>} One score per signature index (first output of the model).
 */
async function runModel(model, tensors) { }

/**
 * Ranks the scores, dropping the signatures that no product has, and keeps the best ones.
 * @param {ArrayLike<number>} scores - One score per signature index.
 * @param {number[]} storeMask - Per signature index: 1 if available, 0 otherwise.
 * @param {number} topN - How many signature indices to keep.
 * @returns {number[]} Signature indices, best score first.
 */
function rankSignatures(scores, storeMask, topN) { }

/**
 * Maps signature indices to product ids through the sqlite (index -> rec_out -> product).
 * Indices with no product are skipped; when several products share a signature the first one is used.
 * @param {Uint8Array} dbBytes - Content of the .sqlite file.
 * @param {number[]} signatureIndices - Signature indices, best first.
 * @returns {Promise<string[]>} Product ids in the same order.
 */
async function fetchRecommendedProductIds(dbBytes, signatureIndices) { }

export { Infer, downloadAssets, loadLiteRt, loadModel };
