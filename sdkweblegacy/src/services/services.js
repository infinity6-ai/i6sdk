/**
 * @fileoverview write it
 */

/**
 * @typedef {import("../index.js").I6SdkConfig} I6SdkConfig
 */

/**
 * write it
 */
class Services {
    /**
     * @param {I6SdkConfig} config - Configuration object.
     * @param {string} [config.baseUrl] - Base URL of the server.
     */
    constructor(config) {
        /**
         * @private
         * @type {I6SdkConfig}
         */
        this._config = config;
    }

    invoker() {
        import { Invoker } from "./invoker/invoker.js";
        return new Invoker(this._config)
    }

    auth() {
        import { Auth } from "./auth/ingest.js";
        return new Auth(this._config)
    }

    ingest() {
        import { Ingest } from "./ingest/ingest.js";
        return new Ingest(this._config)
    }

}

export { Services };

