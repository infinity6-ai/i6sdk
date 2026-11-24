/**
 * @fileoverview Debug logging of the SDK. Set `window.I6_SDK_DEBUG = true` to see it in the console.
 */

/**
 * @param {string} prefix - Printed first, followed by ":" (e.g. "infer").
 * @returns {(...args: *) => void} A logger that prints only when `window.I6_SDK_DEBUG` is true.
 */
const debugLogger = (prefix) => (...args) => { if (globalThis.I6_SDK_DEBUG) console.log(`${prefix}:`, ...args) }

export { debugLogger };
