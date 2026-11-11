/**
 * @fileoverview HTTP client for sending API requests to i6 services.
 */

import { apiErrorFromResponse, apiErrorFromNetwork } from "../errors.js";

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
 *    (e.g., `sig=v1.20261027T150000000Z.payload.PubKeyBase64...SigBase64...`)
 *
 * @typedef {object} ServiceOptions
 * @property {string} [datasetToken] - [Go: string]
 * @property {string} [bearerToken] - [Go: string]
 * @property {Record<string, string>} [labels] - Labels
 * @property {Fetcher} [fetcher] - Custom fetch implementation.
 * @property {PrepareRequest} [prepareRequest] - Hook to mutate request options before sending (e.g., adding `sig` query parameter).
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
 * @property {boolean} [sinkable]
 */

function getValue(obj, path) {
  return path.split('.').reduce((currentObject, key) => {
    return currentObject ? currentObject[key] : undefined;
  }, obj);
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
        opts.url = opts.url || `${this._services.config().baseUrl}/${opts.path}`

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
                opts.url = opts.url + (opts.url.includes('?') ? '&' : '?') + queryString
            }
        }

        opts.headers = new Headers(opts.headers)
        const datasetToken = opts.opts?.datasetToken || this._services.config().datasetToken;
        const bearerToken = opts.opts?.bearerToken || this._services.config().bearerToken || (typeof sessionStorage !== 'undefined' ? sessionStorage.getItem('i6_user_access_token') : null);

        if (!opts.headers.has('Authorization') && datasetToken) {
            opts.headers.set('Authorization', 'Dataset ' + datasetToken)
        }
        if (!opts.headers.has('Authorization') && bearerToken) {
            opts.headers.set('Authorization', 'Bearer ' + bearerToken)
        }
        if (!opts.headers.has('Authorization')) {
            const url = new URL(opts.url);
            const split = url.pathname.split('/')
            const idx = split.indexOf('dataset')
            let token = ''
            if (idx < 0) {
                token = await this._services.tokens().getValidAccessToken()  
                // console.log('user token', token) 
            } else {
                const dataset = split[idx+1]
                token = await this._services.tokens().getValidDatasetAccessToken(dataset)
                // console.log('dataset token', dataset, token)
            }
            if (token) {
                opts.headers.set('Authorization',  'Bearer ' + token)
            }
        }
        console.log(opts)

        if (opts.json) {
            opts.headers.set('Content-Type', 'application/json')
            opts.body = JSON.stringify(opts.json)
        }

        opts.credentials = "include";

        if (opts.opts?.prepareRequest) {
            await opts.opts.prepareRequest(opts)
        }

        // if (opts.sinkable) {
        //     //     if (opts.method?.toUpperCase() != "GET") throw new I6Error('Unsupported: sinkable must be GET')
        //     //     if (opts.body) throw new I6Error('Unsupported: sinkable must not have body')
        //     //     const u = new URL(opts.url)
        //     //     const pathname = u.pathname.replace(/^\/+/, '')
        //     //     let search = u.search.startsWith('?') ? u.search.substring(1) : u.search
        //     //     if (search) {
        //     //         search = b64.encodeUrlSafe(await gzip.gzip((search)))        
        //     //     }
        //     //     if (search.length <= 1900) {
        //     //         opts.url = `${this._services.config().baseUrl}/api/sink/${pathname}?s=${search}`
        //     //     }
        //     console.log('sinkable apis disabled', opts.url)
        //     return new Response(null, {
        //         status: 204,
        //         statusText: 'No Content',
        //     })
        // }

        if (opts.opts?.fetcher) {
            return opts.opts.fetcher(opts.url, opts)
        }

        return fetch(opts.url, opts)
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
