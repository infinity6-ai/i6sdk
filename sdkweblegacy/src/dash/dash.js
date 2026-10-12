/**
 * @fileoverview Queries an i6 dash in the browser: downloads the dash's sqlite, runs the named query
 * `dash-${dashName}` on it and returns its result as is. Nothing here is dash specific.
 *
 * Set `window.I6_DEBUG = true` to see what happens in the console.
 */

import { I6Error } from "../errors.js";
import { fs } from "../internal/fs/fs.js";
import { sqlitez } from "../internal/sqlitez/sqlitez.js";
import { downloadFile, isPartitions, partitionsKey } from "../internal/assets/assets.js";

/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * Options of a dash query.
 * @typedef {Object} DashOptions
 * @property {string} dashName - Dash name (e.g. "product-relevance-fashion-dash-funnel"); names the cached sqlite and the named query.
 * @property {Object<string, string>} partitions - Partitions of the dash, name -> value (e.g. `{ store_id: "s1", version: "v1" }`).
 * @property {string} dataset - Dataset that holds the dash file.
 * @property {*} params - Query input, any JSON. Bound as `?1` of the named query.
 * @property {string} [sqliteUrl] - Tests only: URL of the sqlite (must answer HEAD with an etag). Default: the model API.
 */

/** Name of the sqlite in the model API response. */
const SQLITE_FILE = "sqlite.bin"

/** Prints only when `window.I6_DEBUG` is true. */
const debug = (...args) => { if (globalThis.I6_DEBUG) console.log("dash:", ...args) }

/**
 * Dash client.
 */
class Dash {
  /**
   * Constructs a Dash client instance.
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
   * Downloads the dash sqlite (cached by etag), runs the named query `dash-${dashName}` with `params` and returns its result.
   * @param {DashOptions} opts - Dash options.
   * @returns {Promise<*>} The `ret` of the named query.
   * @throws {I6Error} With code "invalid_argument" if opts is malformed, "asset" if the download fails,
   *   "query" if the sqlite fails.
   */
  async query(opts) {
    validate(opts)

    const url = await this._sqliteUrl(opts)
    const ref = await downloadFile(`dash-${opts.dashName}-${partitionsKey(opts.partitions)}-sqlite`, url)

    const conn = await sqlitez.open({ path: await fs.resolve(ref) })
    try {
      const name = `dash-${opts.dashName}`
      const result = await conn.namedJsonQuery({ name, params: [opts.params] })
      debug(`${name} returned`, result)
      return result
    } finally {
      await conn.close()
    }
  }

  /**
   * Finds where the dash sqlite is: the signed URL of the model API (like Infer), unless opts has `sqliteUrl` (tests).
   * @private
   * @param {DashOptions} opts - Dash options.
   * @returns {Promise<string>} Download URL of the sqlite.
   * @throws {I6Error} With code "asset" if the API gave no URL.
   */
  async _sqliteUrl(opts) {
    if (opts.sqliteUrl) return opts.sqliteUrl

    const { dashName, partitions, dataset } = opts
    const apis = await this._services.apis()
    const resp = await apis.dszSignedUrl({
      params: { dataset: dataset },
      payload: {
        paths: [
          { method: "GET", path: "i6tb://results.relevance-fashion-store-dash/store_id=1/sqlite.bin" },
        ]
      }
    })
    
    const url = resp.payload?.urls?.[0]?.url
    if (!url) throw new I6Error(`dash: the api gave no "${SQLITE_FILE}" file: ${JSON.stringify(resp.error ?? null)}`, { code: "asset" })

    return url
  }
}

/**
 * @param {DashOptions} opts - Options given to `Dash.query`.
 * @throws {I6Error} With code "invalid_argument" if a required option is missing.
 */
function validate(opts) {
  for (const key of ["dashName", "dataset"]) {
    if (!opts || !opts[key]) throw new I6Error(`dash: opts.${key} is required`, { code: "invalid_argument" })
  }
  if (!isPartitions(opts.partitions)) {
    throw new I6Error("dash: opts.partitions must be an object of string values", { code: "invalid_argument" })
  }
  if (opts.params === undefined) throw new I6Error("dash: opts.params is required", { code: "invalid_argument" })
}

export { Dash };
