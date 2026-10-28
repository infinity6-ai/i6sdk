/**
 * @fileoverview Base64 URL safe encoding utilities.
 */

/**
 * Encodes a string or Uint8Array to a Base64 URL safe string.
 * @param {string|Uint8Array|ArrayBuffer} data
 * @returns {string}
 */
export function encodeUrlSafe(data) {
  let bytes;
  if (typeof data === 'string') {
    bytes = new TextEncoder().encode(data);
  } else if (data instanceof ArrayBuffer) {
    bytes = new Uint8Array(data);
  } else {
    bytes = new Uint8Array(data);
  }

  let binaryString = '';
  const chunkSize = 8192;
  for (let i = 0; i < bytes.length; i += chunkSize) {
    binaryString += String.fromCharCode.apply(null, bytes.subarray(i, i + chunkSize));
  }
  
  const base64 = btoa(binaryString);
  return base64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

/**
 * Decodes a Base64 URL safe string to a Uint8Array.
 * @param {string} b64UrlSafe
 * @returns {Uint8Array}
 */
export function decodeUrlSafeToBytes(b64UrlSafe) {
  let base64 = b64UrlSafe.replace(/-/g, '+').replace(/_/g, '/');
  while (base64.length % 4) {
    base64 += '=';
  }
  const binaryString = atob(base64);
  const bytes = new Uint8Array(binaryString.length);
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }
  return bytes;
}

/**
 * Decodes a Base64 URL safe string to a UTF-8 string.
 * @param {string} b64UrlSafe
 * @returns {string}
 */
export function decodeUrlSafe(b64UrlSafe) {
  const bytes = decodeUrlSafeToBytes(b64UrlSafe);
  return new TextDecoder().decode(bytes);
}
