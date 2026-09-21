/**
 * @fileoverview Main entry point for the i6 Web Legacy SDK.
 */

import { Services } from "./services/services.js";

/**
 * @typedef {import("./types.js").I6SdkConfig} I6SdkConfig
 * @typedef {import("./ingest/ingest.js").Ingest} Ingest
 * @typedef {import("./auth/auth.js").Auth} Auth
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
   */
  constructor(config) {
    /**
     * @private
     * @type {Services}
     */
    this._service = new Services(this._config = {
      ...DEFAULT_CONFIG,
      ...config,
    })
  }

  /**
   * Returns an Ingest service instance.
   * @returns {Ingest} An instance of the Ingest client.
   */
  ingest() {
    return this._service.ingest()
  }

  /**
   * Returns an Auth service instance.
   * @returns {Auth} An instance of the Auth client.
   */
  auth() {
    return this._service.auth()
  }

  /**
   * Registers an event listener on the internal event dispatcher.
   * @param {string} evt - The event type or name to listen for.
   * @param {EventListenerOrEventListenerObject|Function} fn - The callback function or event listener object invoked when the event occurs.
   * @returns {void}
   */
  bind(evt, fn) {
    this._service.bind(evt, fn)
  }
}

window.I6Sdk = I6Sdk;
export { I6Sdk };

