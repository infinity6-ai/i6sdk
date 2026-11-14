/**
 * @fileoverview Keeps a remote file cached in the browser's file system (see fs.js),
 * downloading it only when its etag changes.
 */

import { fs } from '../fs/fs.js'

/**
 * @typedef {import('../fs/fs.js').FileRef} FileRef
 */

/**
 * @param {string} message
 * @param {number} status - HTTP status code of the response.
 * @returns {Error} An error with `status`.
 */
function httpError(message, status) {
    return Object.assign(new Error(message), { status })
}

/**
 * Downloads remote files into the versioned file store. Use the `downloader` singleton.
 */
class Downloader {
    /**
     * Downloads the file by name if necessary (checking etag): sends a HEAD request, and
     * if the etag differs from the released version, streams the body into a new version
     * and releases it.
     * @param {Object} opts
     * @param {string} opts.name - required, name to store the file under (see fs.js)
     * @param {string} opts.url - required, URL of the file; must answer HEAD with an etag
     * @returns {Promise<FileRef>} the FileRef of the up-to-date file (basically to get the version)
     * @throws {TypeError} If name or url is missing.
     * @throws {Error} If the HEAD or GET request fails, or the HEAD response has no etag.
     *   When the server answered, the error has `status` (the HTTP status code).
    */
    async update(opts) {
        if (!opts || !opts.name) throw new TypeError('downloader: opts.name is required')
        if (!opts.url) throw new TypeError('downloader: opts.url is required')

        // HEAD (gcs) to get the current etag
        const head = await fetch(opts.url, { method: 'HEAD' })
        if (!head.ok) throw httpError(`downloader: HEAD ${opts.url} failed (${head.status})`, head.status)
            
        const etag = head.headers.get('etag')
        if (!etag) throw new Error('downloader: no etag in HEAD response')

        // Same etag as the released version, nothing to download
        const current = await fs.get({ name: opts.name })
        if (current && current.etag === etag) return current

        // Create the version in the fs, the blob is not written by it
        const ref = await fs.create({ name: opts.name, etag })

        let success = false
        try {
            const response = await fetch(opts.url)
            if (!response.ok) throw httpError(`downloader: GET ${opts.url} failed (${response.status})`, response.status)

            // Stream the body straight into the opfs file, never held in memory
            let dir = await navigator.storage.getDirectory()
            for (const segment of (await fs.resolve(ref)).split('/')) {
                dir = segment === 'blob.bin'
                    ? await dir.getFileHandle(segment, { create: true })
                    : await dir.getDirectoryHandle(segment)
            }
            await response.body.pipeTo(await dir.createWritable())
            success = true
        } finally {
            if (!success) {
                // Do not leave a half-written version behind
                await fs.discard(ref)
            }
        }

        return fs.release(ref)
    }
}

/** Shared Downloader instance. */
const downloader = new Downloader()

export { downloader };
