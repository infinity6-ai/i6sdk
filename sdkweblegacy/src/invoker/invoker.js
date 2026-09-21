class Invoker {

    constructor(services) {
        this._services = services;
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
        let url = opts.url || `${this._services.config().baseUrl}/${opts.path}`

        if (opts.query) {
            const queryString = new URLSearchParams(params).toString()
            if (!url.contains('?')) {
                url = url + '?' + queryString
            } else {
                url = url + '&' + queryString
            }
        }

        if (opts.json) {
            opts.headers = {
                'Content-Type': 'application/json'
            }
            opts.body = JSON.stringify(opts.json)
        }

        console.log("🚀 ~ Invoker ~ invoke ~ url:", url)
        return fetch(url, opts)
    }

}

export { Invoker };
