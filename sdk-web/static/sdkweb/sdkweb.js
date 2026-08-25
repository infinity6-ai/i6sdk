import { RaasBoughtTogether } from './raas/raas-bought-together.js'


class I6Sdk {
  async init(config) {
    this.config = config
    if (!this.config) {
      throw new Error('SDK configuration is missing.');
    }
  }

  async RaasBoughtTogether(config) {
    const computedConfig = {
      ...this.config,
      ...config
    }
    const ret = new RaasBoughtTogether()
    await ret.init(computedConfig)
    return ret
  }
}

// Expose the SDK to the global scope
window.I6Sdk = I6Sdk;
export { I6Sdk };
