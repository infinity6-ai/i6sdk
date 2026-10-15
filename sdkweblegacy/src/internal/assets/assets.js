/**
 * @fileoverview Helpers shared by the services that run on files downloaded from i6 (infer, dash).
 */

import { I6Error } from "../../errors.js";
import { downloader } from "../downloader/downloader.js";

/**
 * @typedef {import("../fs/fs.js").FileRef} FileRef
 */

/**
 * @param {Object<string, string>} partitions - Partitions, name -> value.
 * @returns {string} A stable key of the partitions, safe as part of a cache file name (e.g. "store_id%3Ds1-version%3Dv1").
 */
function partitionsKey(partitions) {
  return Object.entries(partitions).sort(([a], [b]) => (a < b ? -1 : a > b ? 1 : 0))
    .map(([key, value]) => encodeURIComponent(`${key}=${value}`)).join("-")
}

/**
 * @param {*} partitions - Value to check.
 * @returns {boolean} True if it is an object (not an array) whose values are all strings.
 */
function isPartitions(partitions) {
  return !!partitions && typeof partitions === "object" && !Array.isArray(partitions)
    && Object.values(partitions).every((v) => typeof v === "string")
}

/**
 * Brings one remote file to the browser cache, downloading it only if its etag changed.
 * @param {string} name - Name of the cached file.
 * @param {string} url - URL of the file.
 * @returns {Promise<FileRef>} The cached file.
 * @throws {I6Error} With code "asset" if the download fails; `status` is the HTTP status when the server answered
 *   (404 when the file does not exist, e.g. a store or month without data) and `cause` is the underlying error.
 */
async function downloadFile(name, url) {
  try {
    const ref = await downloader.update({ name, url })
    if (globalThis.I6_DEBUG) console.log(`assets: ${name} is ready`, ref)
    return ref
  } catch (cause) {
    throw new I6Error(`could not get ${url}: ${cause.message}`, { code: "asset", status: cause.status, cause })
  }
}

export { partitionsKey, isPartitions, downloadFile };
