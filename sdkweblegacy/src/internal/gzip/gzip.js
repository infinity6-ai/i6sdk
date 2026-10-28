/**
 * @fileoverview Gzip compression utilities.
 */

import { encodeUrlSafe, decodeUrlSafeToBytes } from '../b64/b64.js';

/**
 * Gzips a given JSON object or array.
 * @param {any} obj - The JSON serializable object to compress.
 * @returns {Promise<Uint8Array>} The gzipped bytes.
 */
export async function gzipJson(obj) {
  const str = JSON.stringify(obj);
  const bytes = new TextEncoder().encode(str);
  
  if (typeof CompressionStream === 'undefined') {
    throw new Error("CompressionStream is not supported in this environment.");
  }
  
  const stream = new Response(bytes).body.pipeThrough(new CompressionStream('gzip'));
  return new Uint8Array(await new Response(stream).arrayBuffer());
}

/**
 * Decompresses a gzipped byte array into a JSON object.
 * @param {Uint8Array|ArrayBuffer} compressedBytes - The gzipped data.
 * @returns {Promise<any>} The parsed JSON object.
 */
export async function gunzipJson(compressedBytes) {
  if (typeof DecompressionStream === 'undefined') {
    throw new Error("DecompressionStream is not supported in this environment.");
  }
  
  const stream = new Response(compressedBytes).body.pipeThrough(new DecompressionStream('gzip'));
  const uncompressedBytes = new Uint8Array(await new Response(stream).arrayBuffer());
  const str = new TextDecoder().decode(uncompressedBytes);
  return JSON.parse(str);
}

/**
 * Gzips a JSON object and encodes it as a Base64 URL safe string.
 * @param {any} obj - The JSON serializable object.
 * @returns {Promise<string>} The Base64 URL safe encoded string of the gzipped JSON.
 */
export async function gzipJsonToBase64Url(obj) {
  const compressedBytes = await gzipJson(obj);
  return encodeUrlSafe(compressedBytes);
}

/**
 * Decodes a Base64 URL safe string, gunzips it, and parses it as JSON.
 * @param {string} b64Str - The Base64 URL safe string containing gzipped JSON.
 * @returns {Promise<any>} The parsed JSON object.
 */
export async function gunzipJsonFromBase64Url(b64Str) {
  const compressedBytes = decodeUrlSafeToBytes(b64Str);
  return await gunzipJson(compressedBytes);
}
