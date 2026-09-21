import { Ingestz } from "./ingestz/ingestz.js";

const DEFAULT_CONFIG = {
  baseUrl: "https://decsuite.sandbox.rs.infinity6.ai",
  // baseUrl: "http://localhost:9021",
}

class I6Sdk {

  /**
   * @param {Object} [config]
   * @param {string} [config.baseUrl] - Base URL of the i6 server. Defaults to the sandbox one.
   */
  constructor(config) {
    this._config = {
      ...DEFAULT_CONFIG,
      ...config,
    };
  }

  hello() {
    return 'Hello from i6sdk-web!';
  }

  async ingestz() {
    return new Ingestz(this._config)
  }
}

window.I6Sdk = I6Sdk;
export { I6Sdk };
