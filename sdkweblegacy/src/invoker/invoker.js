/**
 * @fileoverview HTTP client for sending API requests to i6 services.
 */

import { apiErrorFromResponse, apiErrorFromNetwork } from "../errors.js";

/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * Options for sending HTTP requests via the Invoker.
 * @typedef {Object} InvokeOptions
 * @property {string} [url] - The full URL to request. If omitted, defaults to `${baseUrl}/${path}`.
 * @property {string} [path] - The path to append to the baseUrl if `url` is not provided.
 * @property {Record<string, any>} [query] - Key-value pairs to append as query parameters.
 * @property {*} [json] - JSON-serializable body data. If provided, sets Content-Type header to application/json and serializes to body.
 * @property {string} [method] - HTTP method (e.g., 'GET', 'POST', 'PUT', 'DELETE').
 * @property {HeadersInit} [headers] - HTTP headers for the request.
 * @property {BodyInit|null} [body] - HTTP request body.
 * @property {RequestCredentials} [credentials] - Request credentials mode (e.g., 'include', 'same-origin', 'omit').
 * @property {RequestMode} [mode] - Request mode (e.g., 'cors', 'no-cors', 'same-origin').
 * @property {AbortSignal|null} [signal] - An AbortSignal to cancel the request.
 */

/**
 * Invoker client for executing HTTP requests against the i6 API.
 */
class Invoker {

    /**
     * Constructs an Invoker instance.
     * @param {Services} services - The services manager instance.
     */
    constructor(services) {
        /**
         * @private
         * @type {Services}
         */
        this._services = services;
    }

    /**
     * Sends an HTTP request using the Fetch API.
     *
     * @param {InvokeOptions} opts - Request options including path/URL, query params, and fetch options.
     * @param {string} [opts.url] - The full URL to request. If omitted, defaults to `${baseUrl}/${path}`.
     * @param {string} [opts.path] - The path to append to the baseUrl if `url` is not provided.
     * @param {Record<string, any>} [opts.query] - Key-value pairs to append as query parameters.
     * @param {*} [opts.json] - JSON-serializable body data. If provided, sets Content-Type header to application/json and serializes to body.
     * Credentials are always set to "include" so the session cookie is sent. Note that `opts`
     * is modified (credentials, and headers/body when `json` is given). An HTTP error status
     * does not throw; check `ok` on the returned Response.
     * @returns {Promise<Response>} A promise that resolves to the fetch Response.
     * @throws {TypeError} If the network request fails (from fetch).
     */
    async invoke(opts) {
        let url = opts.url || `${this._services.config().baseUrl}/${opts.path}`

        if (opts.query) {
            const search = new URLSearchParams()
            for (const [key, value] of Object.entries(opts.query)) {
                if (value === undefined || value === null) {
                    continue
                }
                // A list becomes one repeated key per item (?a=1&a=2), not "1,2".
                for (const item of Array.isArray(value) ? value : [value]) {
                    // Check if the item is a valid object/array, otherwise use String coercion
                    const serializedValue = (typeof item === 'object' && item !== null) 
                        ? JSON.stringify(item) 
                        : String(item);

                    search.append(key, serializedValue);
                }
            }
            const queryString = search.toString()
            if (queryString) {
                url = url + (url.includes('?') ? '&' : '?') + queryString
            }
        }

        if (opts.json) {
            // Keep the caller's headers; only Content-Type is forced.
            const headers = new Headers(opts.headers)
            headers.set('Content-Type', 'application/json')
            opts.headers = headers
            opts.body = JSON.stringify(opts.json)
        }

        opts.credentials = "include";

        const resp = await fetch(url, opts)
        return resp
    }

    /**
     * Like `invoke`, but never throws: resolves with `{status, ok, <key>}` on success and with
     * `{status, ok: false, error}` (see `I6ApiError`) on an HTTP error or a network failure.
     * A success with status 204 or an empty body resolves with `{status, ok: true}` and no `<key>`.
     * @param {InvokeOptions} opts - Same as `invoke`.
     * @param {"payload"|"body"} [key] - Where the parsed JSON of a success is put.
     * @returns {Promise<{status: number, ok: boolean, error?: import("../errors.js").I6ApiError} & Record<string, any>>}
     */
    async invokeResult(opts, key = "payload") {
        let resp
        try {
            resp = await this.invoke(opts)
        } catch (cause) {
            return { status: 0, ok: false, error: apiErrorFromNetwork(cause) }
        }

        const result = { status: resp.status, ok: resp.ok }
        try {
            if (resp.ok) {
                const text = resp.status === 204 ? "" : await resp.text()
                if (text) {
                    result[key] = JSON.parse(text)
                }
            } else {
                result.error = apiErrorFromResponse(resp.status, await resp.text())
            }
        } catch (cause) {
            result.ok = false
            result.error = apiErrorFromNetwork(cause)
        }
        return result
    }

}

export { Invoker };
