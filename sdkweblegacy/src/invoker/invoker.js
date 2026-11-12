/**
 * @fileoverview HTTP client for sending API requests to i6 services.
 */

import * as b64 from "../internal/b64/b64.js";
import * as gzip from "../internal/gzip/gzip.js";

/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * A custom fetch implementation.
 * @typedef {(url: string, opts: InvokeOptions) => Promise<Response>} Fetcher
 */

/**
 * Hook to prepare or mutate request options before the request is sent (e.g., adding `sig` query parameter).
 * @typedef {(opts: InvokeOptions) => Promise<void> | void} PrepareRequest
 */

/**
 * Service options configuring SDK network requests.
 *
 * This is essential when the `DatasetToken` must never be exposed to the browser.
 * There are two main approaches to keep the `DatasetToken` safe on your backend:
 *
 * ---
 * ### Approach 1: Backend Reverse Proxy (Recommended for simplicity)
 * To avoid all signature complexity ("signature hell"), you can simply proxy requests through
 * your own backend server. The browser sends the request to your backend, and your server forwards
 * it to i6 while attaching the `Authorization: Dataset <datasetToken>` header. This also works for
 * all HTTP methods and request bodies without any signing restrictions.
 *
 * Browser Fetcher example:
 * ```js
 * const fetcher = async (url, opts) => {
 *   // Route the request to your own server proxy endpoint
 *   const proxyUrl = url.replace('https://api.i6.ai', '/api/i6-proxy');
 *   return fetch(proxyUrl, opts);
 * };
 * ```
 *
 * Backend server handler example:
 * ```js
 * // Express / Node.js backend
 * app.all('/api/i6-proxy/*', async (req, res) => {
 *   const targetUrl = req.originalUrl.replace('/api/i6-proxy', 'https://api.i6.ai');
 *   const resp = await fetch(targetUrl, {
 *     method: req.method,
 *     headers: {
 *       ...req.headers,
 *       'Authorization': `Dataset ${process.env.I6_DATASET_TOKEN}`,
 *     },
 *     body: ['GET', 'HEAD'].includes(req.method) ? undefined : req.body,
 *   });
 *   res.status(resp.status).send(await resp.text());
 * });
 * ```
 *
 * ---
 * ### Approach 2: URL Signing (Direct browser-to-i6 calls)
 * If you prefer requests to go directly from the browser to i6 (e.g. to offload traffic from your servers),
 * provide a `prepareRequest` hook that calls your backend server to pre-sign the URL before forwarding it to i6.
 *
 * Example:
 * ```js
 * const prepareRequest = async (opts) => {
 *   // Ask your backend to sign the URL
 *   const res = await fetch('/api/sign-i6-url', {
 *     method: 'POST',
 *     body: JSON.stringify({ url: opts.url, method: opts.method })
 *   });
 *   const { sig } = await res.json();
 *   // Append the sig query parameter to opts.url
 *   opts.url += (opts.url.includes('?') ? '&' : '?') + `sig=${encodeURIComponent(sig)}`;
 * };
 * ```
 *
 * #### How the backend signs the URL
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
 *    Note: The values of fixed params are Base64URL-encoded without padding: `key=b64url(value)`.
 *    Example:
 *    ```text
 *    GET
 *    /api/ingest/dataset/my-dataset/event/relevance-fashion-event/text.txt
 *    20261027T150000000Z
 *    payload
 *    dataset=bXktZGF0YXNldA,sessionid=YWJj,userid=MTIz
 *    ```
 * 3. Sign the message bytes with the Ed25519 private key. To make it using openssl (like in `@cryptz/cryptz25519/cryptz25519openssl`):
 *    ```bash
 *    openssl pkeyutl -sign -inkey private_key.pem -rawin -in message.txt -out signature.bin
 *    ```
 * 4. Encode the signature and public key in Base64URL.
 * 5. Append the `sig` query parameter: `?sig=v1.<exp>.<open_params_pipe_separated>.<pub_key>.<signature>`
 *    (e.g., `sig=v1.20261027T150000000Z.payload.PubKeyBase64...SigBase64...`
 */

/**
 * Options for sending HTTP requests via the Invoker.
 * @typedef {Object} InvokeOptions
 * @property {string} [url] - Full target URL. Overrides baseUrl + path if set.
 * @property {string} [path] - Path appended to baseUrl (e.g. "api/item/create").
 * @property {string} [method="GET"] - HTTP method ("GET", "POST", "PUT", etc.).
 * @property {Record<string, *>} [query] - Query parameters serialized into the URL.
 * @property {*} [json] - Body data serialized as JSON. Sets "Content-Type: application/json".
 * @property {HeadersInit} [headers] - Additional HTTP headers.
 * @property {ServiceOptions} [opts] - Advanced service request hooks (fetcher, prepareRequest, etc.).
 * @property {RequestCredentials} [credentials="include"] - Fetch credentials mode.
 */

function codeFromStatus(status) {
    if (status === 401) return "unauthenticated"
    if (status === 403) return "forbidden"
    if (status === 404) return "not_found"
    if (status >= 500) return "unavailable"
    return "bad_request"
}

function apiErrorFromResponse(status, text) {
    return { status, code: codeFromStatus(status), message: text || `HTTP ${status}`, cause: text }
}

/**
 * Response wrapper holding HTTP transport status and lazy parser.
 * @template T
 */
class Resp {
    #parser
    #parsed = false
    #value = null
    #text = null

    /**
     * @param {number} status - HTTP status code.
     * @param {Headers} headers - HTTP response headers.
     * @param {() => *} parser - Lazy parser function.
     * @param {string} [text] - Raw response body text.
     */
    constructor(status, headers, parser, text = "") {
        this.status = status
        this.headers = headers
        this.ok = this.status >= 200 && this.status < 300
        this.#parser = parser
        this.#text = text
    }

    /**
     * Raw body text of the response.
     * @returns {string}
     */
    get text() {
        return this.#text
    }

    /**
     * Structured error information when response is not ok.
     * @returns {{status: number, code: string, message: string, cause: *}|null}
     */
    get error() {
        if (!this.ok) {
            return apiErrorFromResponse(this.status, this.#text)
        }
        return null
    }

    /**
     * Parses and returns the deserialized response value.
     * Throws an error if deserialization fails or if the response status is not ok.
     * @returns {T}
     */
    parse() {
        if (!this.#parsed) {
            this.#value = this.#parser()
            this.#parsed = true
        }
        return this.#value
    }
}

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
     * @param {{tokenType: string, url: string}} opts
     * @return {{token:string, required: boolean}}
     */
    async resolveToken(opts) {
        let ret = { required: false, token: "" }
        if (!opts.tokenType || opts.tokenType === 'none') {
            return ret
        }
        ret.required = true
        const url = new URL(opts.url);
        if (opts.tokenType === 'user') {
            ret.token = await this._services.tokens().getValidAccessToken()
        } else if (opts.tokenType === 'dataset') {
            const parts = url.pathname.split('/')
            if (parts[0] === '') parts.shift()
            if (parts[2] !== 'dataset') throw new Error(`cannot resolve dataset name ${opts.tokenType}: ${url.pathname}`)
            const dataset = parts[3]
            if (!dataset) throw new Error(`cannot resolve dataset name ${opts.tokenType}: ${url.pathname}`)
            ret.token = await this._services.tokens().getValidDatasetAccessToken(dataset)
        }
        return ret
    }

    /**
     * Sends an HTTP request using the Fetch API.
     *
     * @param {InvokeOptions} opts - Request options including path/URL, query params, and fetch options.
     * @returns {Promise<Response>} A promise that resolves to the fetch Response.
     */
    async invoke(opts) {
        opts.url = opts.url || `${this._services.config().baseUrl}/${opts.path}`

        if (opts.query) {
            const search = new URLSearchParams()
            for (const [key, value] of Object.entries(opts.query)) {
                if (value === undefined || value === null) {
                    continue
                }
                for (const item of Array.isArray(value) ? value : [value]) {
                    const serializedValue = (typeof item === 'object' && item !== null)
                        ? JSON.stringify(item)
                        : String(item);

                    search.append(key, serializedValue);
                }
            }
            const queryString = search.toString()
            if (queryString) {
                opts.url = opts.url + (opts.url.includes('?') ? '&' : '?') + queryString
            }
        }

        opts.headers = new Headers(opts.headers)
        const datasetToken = opts.opts?.datasetToken;
        const bearerToken = opts.opts?.bearerToken;

        if (!opts.headers.has('Authorization') && datasetToken) {
            opts.headers.set('Authorization', 'Dataset ' + datasetToken)
        }
        if (!opts.headers.has('Authorization') && bearerToken) {
            opts.headers.set('Authorization', 'Bearer ' + bearerToken)
        }
        if (!opts.headers.has('Authorization') && opts.opts?.labels?.token_type && opts.opts.labels.token_type !== 'none') {
            const tk = await this.resolveToken({
                tokenType: opts.opts.labels.token_type,
                url: opts.url
            })
            if (tk.token) {
                opts.headers.set('Authorization', 'Bearer ' + tk.token)
            } else if (tk.required) {
                if (opts.opts.labels.token_type == 'user') {
                    return new Response(JSON.stringify({ error: "Unauthorized", message: "Unauthorized" }), { status: 401, statusText: "Unauthorized" });
                }
                return new Response(JSON.stringify({ error: "Forbidden", message: "Forbidden" }), { status: 403, statusText: "Forbidden" });
            }
        }

        if (opts.json) {
            opts.headers.set('Content-Type', 'application/json')
            opts.body = JSON.stringify(opts.json)
        }

        opts.credentials = "include";

        if (opts.opts?.prepareRequest) {
            await opts.opts.prepareRequest(opts)
        }

        if (opts.opts?.fetcher) {
            return opts.opts.fetcher(opts.url, opts)
        }

        return fetch(opts.url, opts)
    }

    /**
     * Resolves with a `Resp` instance holding the HTTP status, headers, and lazy `parse()`.
     * Network errors bubble up directly as native exceptions.
     * @template T
     * @param {InvokeOptions} opts - Same as `invoke`.
     * @returns {Promise<Resp<T>>}
     */
    async invokeResult(opts) {
        const resp = await this.invoke(opts)
        const text = resp.status === 204 ? "" : await resp.text()

        const parser = () => {
            if (!resp.ok) {
                const apiErr = apiErrorFromResponse(resp.status, text)
                const err = new Error(apiErr.message)
                err.status = apiErr.status
                err.code = apiErr.code
                err.cause = apiErr.cause
                throw err
            }
            if (!text) {
                return null
            }
            return JSON.parse(text)
        }

        return new Resp(resp.status, resp.headers, parser, text)
    }

}

export { Invoker, Resp };
