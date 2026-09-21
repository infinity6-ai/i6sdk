/**
 * @fileoverview Main entry point for the i6 Web Legacy SDK.
 */

import { Auth } from "./auth/auth.js";
import { Ingest } from "./ingest/ingest.js";

/**
 * @typedef {Object} I6SdkConfig
 * @property {string} [baseUrl] - Base URL of the i6 server. Defaults to the sandbox environment.
 */

/**
 * Default configuration for the SDK.
 * @type {I6SdkConfig}
 */
const DEFAULT_CONFIG = {
  baseUrl: "https://decsuite.sandbox.rs.infinity6.ai",
  // baseUrl: "http://localhost:9021",
}

/**
 * Main SDK class providing access to i6 services.
 */
class I6Sdk {

  /**
   * Constructs an instance of the i6 SDK.
   * @param {I6SdkConfig} [config] - Configuration options for the SDK.
   * @param {string} [config.baseUrl] - Base URL of the i6 server. Defaults to the sandbox one.
   */
  constructor(config) {
    /**
     * @private
     * @type {I6SdkConfig}
     */
    this._config = {
      ...DEFAULT_CONFIG,
      ...config,
    };
  }

  /**
   * Returns a greeting message.
   * @returns {string} The greeting message.
   */
  hello() {
    return 'Hello from i6sdk-web!';
  }

  /**
   * Returns an Ingest service instance.
   * @returns {Promise<Ingest>} An instance of the Ingest client.
   */
  async ingest() {
    return new Ingest(this._config)
  }

  /**
   * Returns an Auth service instance.
   * @returns {Promise<Auth>} An instance of the Auth client.
   */
  async auth() {
    return new Auth(this._config)
  }
}

window.I6Sdk = I6Sdk;
export { I6Sdk };

