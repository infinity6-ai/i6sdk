// Downloader depends on fs.js
// By now we dont support resuming the Download
// We cant download the entire file in the memory (not even in blocks)
// Use browser API to "Download this file in this opfs"

class Downloader {
    /**
     * Download the file by name if necessary (checking etag)
     * @param {Object} opts
     * @param {string} opts.name - required
     * @param {string} opts.url - required
     * @returns {Promise<FileRef>} the FileRef of the downloaded file (basically to get the version)
    */
    async update(opts) {
        // Resolve o name

        // Primeiro criamos um head (metodo do gcs) num arquivo para pegar o etag

        // Se o etag que existir no head for igual ao daqui, nem faz o download
        
        // 

        // Depois passa o etag pra criar o arquivo no fs
        
        // Depois faz o download
    }
}

const downloader = new Downloader()

export { downloader };

