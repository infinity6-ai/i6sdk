import { Ingest } from "./ingest/ingest.js";

const DEFAULT_BASE_URL = "https://decsuite.sandbox.rs.infinity6.ai";
// const DEFAULT_BASE_URL = "http://localhost:9021";

class I6Sdk {
  _config = {
    baseUrl: DEFAULT_BASE_URL,
  };

  /**
   * @param {Object} [config]
   * @param {string} [config.baseUrl] - Base URL of the i6 server. Defaults to the sandbox one.
   */
  constructor(config) {
    if (config) {
      this.config(config);
    }
  }

  /**
   * @param {Object} config
   * @param {string} [config.baseUrl] - Base URL of the i6 server. Defaults to the sandbox one.
   */
  config(config) {
    if (!config) {
      throw new Error('SDK configuration is missing.');
    }
    this._config = {
      baseUrl: config.baseUrl ?? DEFAULT_BASE_URL,
    };
  }

  hello() {
    return 'Hello from i6sdk-web!';
  }

  async ingest() {
    return new Ingest(this._config)
  }
}

window.I6Sdk = I6Sdk;
export { I6Sdk };
