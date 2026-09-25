// Downloader depends on fs


class Downloader {
    /**
     * Download the file by name
     * @param {Object} opts
     * @param {string} opts.name - required
     * @param {string} opts.url - required
     * @returns {Promise<void>}
    */
    async download(opts) {
    }

}

const downloader = new Downloader()

export { downloader };

