import { Invoker } from "../invoker/invoker.js";
import { Auth } from "../auth/auth.js";
import { Ingest } from "../ingest/ingest.js";

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

    config() {
        return this._config
    }

    invoker() {
        const ret = new Invoker(this)
        if (ret.prepare) {
            ret.prepare()
        }
        return ret
    }

    auth() {
        const ret = new Auth(this)
        if (ret.prepare) {
            ret.prepare()
        }
        return ret
    }

    ingest() {
g        const ret = new Ingest(this)
        if (ret.prepare) {
            ret.prepare()
        }
        return ret
    }

}

export { Services };

