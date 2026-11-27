const debugLogger = (prefix) => (...args) => {
  if (globalThis.I6_SDK_DEBUG) console.log(`${prefix}:`, ...args);
};
export { debugLogger };
