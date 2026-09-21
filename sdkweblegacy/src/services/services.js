/**
 * @fileoverview Service container managing SDK service instances and configuration.
 */

import { Invoker } from "../invoker/invoker.js";
import { Auth } from "../auth/auth.js";
import { Ingest } from "../ingest/ingest.js";

/**
 * @typedef {import("../types.js").I6SdkConfig} I6SdkConfig
 */

/**
 * Service container that instantiates and provides access to SDK service clients.
 */
class Services {
    /**
     * Constructs a Services instance.
     * @param {I6SdkConfig} config - Configuration object for the SDK services.
     */
    constructor(config) {
        /**
         * @private
         * @type {I6SdkConfig}
         */
        this._config = config;
    }

    /**
     * Returns the SDK configuration.
     * @returns {I6SdkConfig} The configuration object.
     */
    config() {
        return this._config
    }

    bind(evt, fn) {
        this._config.dispatcher.addEventListener(evt, fn);
    }

    /**
     * Creates and returns an Invoker service instance.
     * @returns {Invoker} An instance of the Invoker client.
     */
    invoker() {
        const ret = new Invoker(this)
        if (ret.prepare) {
            ret.prepare()
        }
        return ret
    }

    /**
     * Creates and returns an Auth service instance.
     * @returns {Auth} An instance of the Auth client.
     */
    auth() {
        const ret = new Auth(this)
        if (ret.prepare) {
            ret.prepare()
        }
        return ret
    }

    /**
     * Creates and returns an Ingest service instance.
     * @returns {Ingest} An instance of the Ingest client.
     */
    ingest() {
        const ret = new Ingest(this)
        if (ret.prepare) {
            ret.prepare()
        }
        return ret
    }


}

export { Services };

