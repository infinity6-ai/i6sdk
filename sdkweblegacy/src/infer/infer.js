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
 * @property {Record<string, ModelInput>} inputs - Tensors fed to the model, by model input name (exact match).
 */

/**
 * One model input tensor, built by the caller.
 * @typedef {Object} ModelInput
 * @property {ArrayLike<number>} data - Flat values, row-major.
 * @property {number[]} [shape] - Tensor shape. Defaults to the model's shape with -1 as 1.
 */

/**
 * Result of an inference run.
 * @typedef {Object} InferResult
 * @property {string[]} productIds - Recommended product ids, best first. Empty when none of the input products is known to the model.
 */

/**
 * Downloaded model assets.
 * @typedef {Object} ModelAssets
 * @property {string} dbPath - OPFS path of the .sqlite file (open it with sqlitez).
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
    const { model, params, assets, inputs } = opts

    // 1. Get the model files (cached by etag) and the runtimes
    const { dbPath, modelBytes } = await downloadAssets(model, assets)
    const litert = await loadLiteRt(assets.litertWasm)
    const compiled = await loadModel(litert, modelBytes)

    // The sqlite is only needed while the stages below run
    const conn = await sqlitez.open({ path: dbPath })
    try {
      // 2. Read the input products from the sqlite
      const productIds = [...new Set(params.productIds)]
      const { storeMask } = await fetchShopData(conn, model, productIds)

      // 3. Run the model
      const tensors = buildInputTensors(litert, compiled, inputs)
      const outputs = await runModel(compiled, tensors)
      const scores = Object.values(outputs)[0]

      // 4. Rank the scores and map the best ones back to product ids
      const signatureIndices = rankSignatures(scores, storeMask, TOP_N)
      return { productIds: await fetchRecommendedProductIds(conn, model, signatureIndices) }
    } finally {
      await conn.close()
      console.log("infer: sqlite closed")
    }
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

/**
 * Reads from the sqlite the rows of the input products and the store mask, with the named
 * query `model-${model}-input` (param `{productIds}`, ret `{products, storeMask}`).
 * @param {SqliteConn} conn - Open connection to the model's sqlite.
 * @param {string} model - Model name (e.g. "fbt").
 * @param {string[]} productIds - Unique ids of the input products.
 * @returns {Promise<ShopData>} The product rows by id and the mask of every signature index.
 * @throws {I6Error} With code "query" if the named query is missing or fails.
 */
async function fetchShopData(conn, model, productIds) {
  const ret = await conn.namedJsonQuery({ name: `model-${model}-input`, param: { productIds } })
  const { products, storeMask } = ret

  const productFeaturesMap = {}
  for (const row of products) productFeaturesMap[row.product_id] = row

  const found = Object.keys(productFeaturesMap)
  console.log(`infer: sqlite: requested ${productIds.length} products ${JSON.stringify(productIds)}, found ${found.length} ${JSON.stringify(found)}`)
  for (const id of found) {
    const row = productFeaturesMap[id]
    console.log(`infer: sqlite: product ${id}: rec_in has ${row.rec_in ? row.rec_in.split("|").length : 0} values, rec_out=${row.rec_out}, sig_index=${row.sig_index}`)
  }
  console.log(`infer: sqlite: storeMask has ${storeMask.length} signatures, ${storeMask.filter((m) => m === 1).length} available`)

  return { productFeaturesMap, storeMask }
}

/** Typed array of each dtype the runner accepts. */
const DTYPE_CTORS = { float32: Float32Array, int32: Int32Array, uint8: Uint8Array }

/**
 * Turns the caller's inputs into model tensors. Every model input must be in `inputs` under its exact name.
 * @param {Object} litert - The LiteRT module.
 * @param {Object} model - The compiled model.
 * @param {Record<string, ModelInput>} inputs - Input name -> data and shape.
 * @returns {Record<string, Object>} Input name -> tensor.
 * @throws {I6Error} With code "invalid_argument" if an input is missing or its data does not fit its shape.
 */
function buildInputTensors(litert, model, inputs) {
  const tensors = {}
  for (const detail of model.getInputDetails()) {
    const input = inputs?.[detail.name]
    if (!input) throw new I6Error(`infer: opts.inputs is missing "${detail.name}"`, { code: "invalid_argument" })

    const shape = input.shape ?? Array.from(detail.shape, (d) => (d === -1 ? 1 : d))
    if (input.data.length !== shape.reduce((a, b) => a * b, 1)) {
      throw new I6Error(`infer: input "${detail.name}" has ${input.data.length} values, does not fit shape [${shape}]`, { code: "invalid_argument" })
    }
    tensors[detail.name] = new litert.Tensor(new DTYPE_CTORS[detail.dtype](input.data), shape)
  }
  return tensors
}

/**
 * Runs the model and frees the input and output tensors.
 * @param {Object} model - The compiled model.
 * @param {Record<string, Object>} tensors - Input name -> tensor, from buildInputTensors.
 * @returns {Promise<Record<string, ArrayLike<number>>>} Output name -> values, for every output of the model.
 */
async function runModel(model, tensors) {
  let raw = {}
  try {
    raw = await model.run(tensors)
    const outputs = {}
    for (const [name, tensor] of Object.entries(raw)) outputs[name] = await tensor.data()
    return outputs
  } finally {
    Object.values(raw).forEach((t) => t.delete())
    Object.values(tensors).forEach((t) => t.delete())
  }
}

/**
 * Ranks the scores, dropping the signatures that no product has, and keeps the best ones.
 * @param {ArrayLike<number>} scores - One score per signature index.
 * @param {number[]} storeMask - Per signature index: 1 if available, 0 otherwise.
 * @param {number} topN - How many signature indices to keep.
 * @returns {number[]} Signature indices, best score first.
 */
function rankSignatures(scores, storeMask, topN) { }

/**
 * Maps signature indices to product ids with the named query `model-${model}-output`
 * (param `{signatureIndices}`, ret `{matches: [{sig_index, sig_value, productIds}]}`).
 * Indices with no product are skipped; when several products share a signature the first one is used.
 * @param {SqliteConn} conn - Open connection to the model's sqlite.
 * @param {string} model - Model name (e.g. "fbt").
 * @param {number[]} signatureIndices - Signature indices, best first.
 * @returns {Promise<string[]>} Product ids in the same order.
 * @throws {I6Error} With code "query" if the named query is missing or fails.
 */
async function fetchRecommendedProductIds(conn, model, signatureIndices) {
  if (signatureIndices.length === 0) return []

  const { matches } = await conn.namedJsonQuery({ name: `model-${model}-output`, param: { signatureIndices } })
  const productIds = matches.map((m) => m.productIds[0])

  console.log(`infer: sqlite: ${signatureIndices.length} signatures -> ${matches.length} with products`, matches.map((m) => ({ sig_index: m.sig_index, candidates: m.productIds.length, chosen: m.productIds[0] })))
  return productIds
}

export { Infer };
