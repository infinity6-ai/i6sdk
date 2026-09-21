import { Auth } from "./auth/auth.js";
import { Ingest } from "./ingest/ingest.js";

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

  async ingest() {
    return new Ingest(this._config)
  }

  async auth() {
    return new Auth(this._config)
  }
}

window.I6Sdk = I6Sdk;
export { I6Sdk };
