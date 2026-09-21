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
        const ret = new Invoker()
        if (ret.prepare) {
            ret.prepare(this._config)
        }
        return ret
    }

    auth() {
        import { Auth } from "./auth/auth.js";
        const ret = new Auth()
        if (ret.prepare) {
            ret.prepare(this._config)
        }
        return ret
    }

    ingest() {
        import { Ingest } from "./ingest/ingest.js";
        const ret = new Ingest()
        if (ret.prepare) {
            ret.prepare(this._config)
        }
        return ret
    }

}

export { Services };

