/**
 * @fileoverview Gzip compression utilities.
 */

import { encodeUrlSafe, decodeUrlSafeToBytes } from '../b64/b64.js';

/**
 * Compresses data (string, Uint8Array, or ArrayBuffer) using gzip.
 * @param {string|Uint8Array|ArrayBuffer} data - The data to compress.
 * @returns {Promise<Uint8Array>} The gzipped bytes.
 */
export async function gzip(data) {
  let bytes;
  if (typeof data === 'string') {
    bytes = new TextEncoder().encode(data);
  } else if (data instanceof ArrayBuffer) {
    bytes = new Uint8Array(data);
  } else {
    bytes = new Uint8Array(data);
  }

  if (typeof CompressionStream === 'undefined') {
    throw new Error("CompressionStream is not supported in this environment.");
  }

  const stream = new Response(bytes).body.pipeThrough(new CompressionStream('gzip'));
  return new Uint8Array(await new Response(stream).arrayBuffer());
}

/**
 * Decompresses gzipped bytes to a UTF-8 string.
 * @param {Uint8Array|ArrayBuffer} compressedBytes - The gzipped data.
 * @returns {Promise<string>} The decompressed string.
 */
export async function gunzipString(compressedBytes) {
  if (typeof DecompressionStream === 'undefined') {
    throw new Error("DecompressionStream is not supported in this environment.");
  }

  const stream = new Response(compressedBytes).body.pipeThrough(new DecompressionStream('gzip'));
  const uncompressedBytes = new Uint8Array(await new Response(stream).arrayBuffer());
  return new TextDecoder().decode(uncompressedBytes);
}

/**
 * Decompresses gzipped bytes and returns raw Uint8Array.
 * @param {Uint8Array|ArrayBuffer} compressedBytes - The gzipped data.
 * @returns {Promise<Uint8Array>} The decompressed bytes.
 */
export async function gunzip(compressedBytes) {
  if (typeof DecompressionStream === 'undefined') {
    throw new Error("DecompressionStream is not supported in this environment.");
  }

  const stream = new Response(compressedBytes).body.pipeThrough(new DecompressionStream('gzip'));
  return new Uint8Array(await new Response(stream).arrayBuffer());
}

/**
 * Gzips a string and encodes it as a Base64 URL safe string.
 * @param {string} str - The string to compress and encode.
 * @returns {Promise<string>} The Base64 URL safe string.
 */
export async function gzipStringToBase64Url(str) {
  const compressedBytes = await gzip(str);
  return encodeUrlSafe(compressedBytes);
}

/**
 * Decodes a Base64 URL safe string and decompresses it as a UTF-8 string.
 * @param {string} b64Str - The Base64 URL safe string.
 * @returns {Promise<string>} The decompressed string.
 */
export async function gunzipStringFromBase64Url(b64Str) {
  const compressedBytes = decodeUrlSafeToBytes(b64Str);
  return await gunzipString(compressedBytes);
}

/**
 * Gzips a given JSON object or array.
 * @param {any} obj - The JSON serializable object to compress.
 * @returns {Promise<Uint8Array>} The gzipped bytes.
 */
export async function gzipJson(obj) {
  return await gzip(JSON.stringify(obj));
}

/**
 * Decompresses a gzipped byte array into a JSON object.
 * @param {Uint8Array|ArrayBuffer} compressedBytes - The gzipped data.
 * @returns {Promise<any>} The parsed JSON object.
 */
export async function gunzipJson(compressedBytes) {
  const str = await gunzipString(compressedBytes);
  return JSON.parse(str);
}

/**
 * Gzips a JSON object and encodes it as a Base64 URL safe string.
 * @param {any} obj - The JSON serializable object.
 * @returns {Promise<string>} The Base64 URL safe encoded string of the gzipped JSON.
 */
export async function gzipJsonToBase64Url(obj) {
  return await gzipStringToBase64Url(JSON.stringify(obj));
}

/**
 * Decodes a Base64 URL safe string, gunzips it, and parses it as JSON.
 * @param {string} b64Str - The Base64 URL safe string containing gzipped JSON.
 * @returns {Promise<any>} The parsed JSON object.
 */
export async function gunzipJsonFromBase64Url(b64Str) {
  const str = await gunzipStringFromBase64Url(b64Str);
  return JSON.parse(str);
}
