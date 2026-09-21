/**
 * @fileoverview HTTP request invoker.
 */

/**
 * @typedef {Object} InvokerConfig
 * @property {string} baseUrl - Default base URL of the ingestz server.
 */

/**
 * @typedef {RequestInit & {
 *   url?: string,
 *   path?: string,
 *   query?: Record<string, any>
 * }} InvokeOptions
 */

/**
 * Class representing an HTTP invoker for executing API requests.
 */
class Invoker {
    /**
     * @param {InvokerConfig} config
     * @param {string} config.baseUrl - Default base URL of the ingestz server.
     */
    constructor(config) {
        /**
         * @private
         * @type {InvokerConfig}
         */
        this._config = config;
    }

    /**
     * Sends an HTTP request using the Fetch API.
     *
     * @param {InvokeOptions} opts - Request options including path/URL, query params, and fetch options.
     * @param {string} [opts.url] - The full URL to request. If omitted, defaults to `${baseUrl}/${path}`.
     * @param {string} [opts.path] - The path to append to the baseUrl if `url` is not provided.
     * @param {Record<string, any>} [opts.query] - Key-value pairs to append as query parameters.
     * @returns {Promise<Response>} A promise that resolves to the fetch Response.
     */
    async invoke(opts) {
        let url = opts.url || `${this._config.baseUrl}/${opts.path}`

        if (opts.query) {
            const queryString = new URLSearchParams(params).toString()
            if (!url.contains('?')) {
                url = url + '?' + queryString
            } else {
                url = url + '&' + queryString
            }
        }

        return fetch(url, opts)
    }

}

export { Auth };
