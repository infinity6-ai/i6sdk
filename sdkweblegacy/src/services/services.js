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
        const ret = new Invoker(this._config)
        if (ret.prepare) {
            ret.prepare()
        }
        return ret
    }

    auth() {
        import { Auth } from "./auth/auth.js";
        const ret = new Auth(this._config)
        if (ret.prepare) {
            ret.prepare()
        }
        return ret
    }

    ingest() {
        import { Ingest } from "./ingest/ingest.js";
        const ret = new Ingest(this._config)
        if (ret.prepare) {
            ret.prepare()
        }
        return ret
    }

}

export { Services };

