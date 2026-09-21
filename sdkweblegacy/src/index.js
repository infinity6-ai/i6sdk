/**
 * @fileoverview Main entry point for the i6 Web Legacy SDK.
 */

import { Services } from "./services/services.js";

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
    this._service = new Services(this._config = {
      ...DEFAULT_CONFIG,
      ...config,
    })
  }

  /**
   * Returns an Ingest service instance.
   * @returns {Promise<Ingest>} An instance of the Ingest client.
   */
  ingest() {
    return this._service.ingest()
  }

  /**
   * Returns an Auth service instance.
   * @returns {Promise<Auth>} An instance of the Auth client.
   */
  auth() {
    return this._service.auth()
  }

}

window.I6Sdk = I6Sdk;
export { I6Sdk };

