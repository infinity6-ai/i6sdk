class Invoker {
    /**
     * @param {Object} config
     * @param {string} config.baseUrl - Default base URL of the ingestz server.
     */
    constructor(config) {
        this._config = config;
    }

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
