/**
 * @fileoverview Runs an i6 model (e.g. "fbt") in the browser and returns the result.
 *
 * Flow: download the model's sqlite + tflite, read the features of the input products from
 * the sqlite, feed them to the tflite model, rank the model scores and map the best ones
 * back to product ids through the sqlite. Nothing is rendered, the caller uses the result.
 */

/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * URLs of the model assets.
 * @typedef {Object} InferAssets
 * @property {string} sqlite - URL of the model's .sqlite file (must answer HEAD with an etag).
 * @property {string} tflite - URL of the model's .tflite file (must answer HEAD with an etag).
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
  constructor(services) {}

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
    const litert = await loadLiteRt()
    const compiled = await loadModel(litert, modelBytes)

    // 2. Read the input products from the sqlite
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

/**
 * Downloads (or reads from the browser cache, checking the etag) the model assets and validates them.
 * @param {string} model - Model name, used to name the cached files.
 * @param {InferAssets} assets - Asset URLs.
 * @returns {Promise<ModelAssets>} The bytes of both files.
 * @throws {I6Error} With code "asset" if a download fails or a file is not a valid sqlite/tflite.
 */
async function downloadAssets(model, assets) {}

/**
 * Loads the bundled sql.js (wasm) once.
 * @returns {Promise<Object>} The initialized sql.js module (`SQL`).
 */
async function loadSqlJs() {}

/**
 * Loads the bundled LiteRT runtime (wasm) once.
 * @returns {Promise<Object>} The LiteRT module.
 */
async function loadLiteRt() {}

/**
 * Compiles the tflite model, reusing the previous compilation when the bytes are the same.
 * @param {Object} litert - The LiteRT module.
 * @param {Uint8Array} modelBytes - Content of the .tflite file.
 * @returns {Promise<Object>} The compiled model.
 */
async function loadModel(litert, modelBytes) {}

/**
 * Reads from the sqlite the rows of the input products and the store mask.
 * @param {Uint8Array} dbBytes - Content of the .sqlite file.
 * @param {string[]} productIds - Unique ids of the input products.
 * @returns {Promise<ShopData>} The product rows by id and the mask of every signature index.
 */
async function fetchShopData(dbBytes, productIds) {}

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
function buildInputTensors(litert, model, productIds, productFeaturesMap) {}

/**
 * Runs the model and frees the input and output tensors.
 * @param {Object} model - The compiled model.
 * @param {Record<string, Object>} tensors - Input name -> tensor, from buildInputTensors.
 * @returns {Promise<Float32Array>} One score per signature index (first output of the model).
 */
async function runModel(model, tensors) {}

/**
 * Ranks the scores, dropping the signatures that no product has, and keeps the best ones.
 * @param {ArrayLike<number>} scores - One score per signature index.
 * @param {number[]} storeMask - Per signature index: 1 if available, 0 otherwise.
 * @param {number} topN - How many signature indices to keep.
 * @returns {number[]} Signature indices, best score first.
 */
function rankSignatures(scores, storeMask, topN) {}

/**
 * Maps signature indices to product ids through the sqlite (index -> rec_out -> product).
 * Indices with no product are skipped; when several products share a signature the first one is used.
 * @param {Uint8Array} dbBytes - Content of the .sqlite file.
 * @param {number[]} signatureIndices - Signature indices, best first.
 * @returns {Promise<string[]>} Product ids in the same order.
 */
async function fetchRecommendedProductIds(dbBytes, signatureIndices) {}

export { Infer };
