// Downloader depends on fs.js
// By now we dont support resuming the Download
// We cant download the entire file in the memory (not even in blocks)
// Use browser API to "Download this file in this opfs"

import { fs } from '../fs/fs.js'

class Downloader {
    /**
     * Download the file by name if necessary (checking etag)
     * @param {Object} opts
     * @param {string} opts.name - required
     * @param {string} opts.url - required
     * @returns {Promise<FileRef>} the FileRef of the downloaded file (basically to get the version)
    */
    async update(opts) {
        if (!opts || !opts.name) throw new TypeError('downloader: opts.name is required')
        if (!opts.url) throw new TypeError('downloader: opts.url is required')

        // HEAD (gcs) to get the current etag
        const head = await fetch(opts.url, { method: 'HEAD' })
        if (!head.ok) throw new Error(`downloader: HEAD ${opts.url} failed (${head.status})`)
            
        const etag = head.headers.get('etag')
        if (!etag) throw new Error('downloader: no etag in HEAD response')

        // Same etag as the released version, nothing to download
        const current = await fs.get({ name: opts.name })
        if (current && current.etag === etag) return current

        // Create the version in the fs, the blob is not written by it
        const ref = await fs.create({ name: opts.name, etag, data: '' })

        const response = await fetch(opts.url)
        if (!response.ok) throw new Error(`downloader: GET ${opts.url} failed (${response.status})`)

        // Stream the body straight into the opfs file, never held in memory
        let dir = await navigator.storage.getDirectory()
        for (const segment of (await fs.resolve(ref)).split('/')) {
            dir = segment === 'blob.bin'
                ? await dir.getFileHandle(segment, { create: true })
                : await dir.getDirectoryHandle(segment)
        }
        await response.body.pipeTo(await dir.createWritable())

        return fs.release(ref)
    }
}

const downloader = new Downloader()

export { downloader };
