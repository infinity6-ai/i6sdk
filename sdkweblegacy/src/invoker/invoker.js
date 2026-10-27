/**
 * @fileoverview HTTP client for sending API requests to i6 services.
 */

import { apiErrorFromResponse, apiErrorFromNetwork } from "../errors.js";

/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * A custom fetch implementation.
 *
 * This is primarily used to implement URL signing for secure requests so that the
 * `DatasetToken` never comes to the browser.
 *
 * ### Using Fetcher to sign URLs
 * Instead of setting `datasetToken` in the browser SDK, you can provide a custom `fetcher`
 * that calls your backend server to sign the request.
 *
 * Example:
 * ```js
 * const fetcher = async (url, opts) => {
 *   // Ask your backend to sign the URL
 *   const res = await fetch('/api/sign-i6-url', {
 *     method: 'POST',
 *     body: JSON.stringify({ url, method: opts.method })
 *   });
 *   const { signedUrl } = await res.json();
 *   // Send the actual request to i6 using the signed URL
 *   return fetch(signedUrl, opts);
 * };
 * ```
 *
 * ### How the backend signs the URL
 * The `DatasetToken` is a JWT securely stored on your backend. The Ed25519 private key is located inside
 * the JWT payload's `ext.priv_key` field (Base64URL encoded). The public key `jti` is registered in our storage.
 * Signed URLs will work while this specific `jti` is not revoked.
 *
 * Nowadays, we can sign only:
 * - `GET` requests
 * - Path
 * - Expiration timestamp (`exp`)
 * - Allowed open query params
 * - Allowed query params with specific fixed values
 *
 * **Example using `ingestzEventRelevanceFashionEvent`:**
 * The browser sends some fixed values (like `userid`, `sessionid`, `dataset`) and some open values (like `click` or `purchase` inside `payload`).
 *
 * 1. Define your expiration time (`exp`) in ISO8601BasicMillis (e.g., `20261027T150000000Z`).
 * 2. Build the canonical message:
 *    `METHOD\nPATH\nEXP\nOPEN_PARAMS_SORTED_COMMA_SEPARATED\nFIXED_PARAMS_KEY_SORTED_COMMA_SEPARATED`
 *    Example:
 *    ```text
 *    GET
 *    /api/ingest/dataset/my-dataset/event/relevance-fashion-event/text.txt
 *    20261027T150000000Z
 *    payload
 *    dataset=my-dataset,sessionid=abc,userid=123
 *    ```
 * 3. Sign the message bytes with the Ed25519 private key. To make it using openssl (like in `@cryptz/cryptz25519/cryptz25519openssl`):
 *    ```bash
 *    openssl pkeyutl -sign -inkey private_key.pem -rawin -in message.txt -out signature.bin
 *    ```
 * 4. Encode the signature and public key in Base64URL.
 * 5. Append the `sig` query parameter: `?sig=v1.<exp>.<open_params_pipe_separated>.<pub_key>.<signature>`
 *    (e.g., `sig=v1.20261027T150000000Z.payload.PubKeyBase64...SigBase64...`)
 *
 * @typedef {(url: string, opts: InvokeOptions) => Promise<Response>} Fetcher
 */

/**
 * @typedef {object} ServiceOptions
 * @property {string} [datasetToken] - [Go: string]
 * @property {string} [sessionId] - Id of the user session; sent in the recommendation events.
 * @property {string} [userId] - Id of the user; sent in the recommendation events.
 * @property {Fetcher} [fetcher] - Custom fetch implementation.
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
 * @property {ServiceOptions} [opts]
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
                    // Objects are sent as JSON, everything else as its string
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

        opts.headers = new Headers(opts.headers)
        if (opts.opts?.datasetToken) {
            opts.headers.set('Authorization', 'Dataset ' + opts.opts?.datasetToken)
        }

        if (opts.json) {
            opts.headers.set('Content-Type', 'application/json')
            opts.body = JSON.stringify(opts.json)
        }

        opts.credentials = "include";

        if (opts.opts?.fetcher) {
            return opts.opts.fetcher(url, opts)
        }

        return fetch(url, opts)
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
