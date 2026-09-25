// Downloader depends on fs
// By now we dont support resuming the Download

class Downloader {
    /**
     * Download the file by name
     * @param {Object} opts
     * @param {string} opts.name - required
     * @param {string} opts.url - required
     * @returns {Promise<FileRef>} the FileRef of the downloaded file (basically to get the version)
    */
    async download(opts) {
        
    }
}

const downloader = new Downloader()

export { downloader };

