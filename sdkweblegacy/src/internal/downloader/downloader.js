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
    constructor() {
        /** @type {Map<string, Promise<FileRef>>} */
        this.downloads = new Map()
    }

    /**
     * Downloads or revalidates the file by name.
     * If minPreloadTime is specified:
     *   - If no file exists, downloads and awaits the new file.
     *   - If a file exists and is older than minPreloadTime, starts updating in the background and returns the existing file immediately.
     *   - If a file exists and is newer than minPreloadTime, returns the existing file immediately without network requests.
     * If minPreloadTime is omitted, always performs an update check and waits for it (legacy behavior).
     * @param {Object} opts
     * @param {string} opts.name - required, name to store the file under (see fs.js)
     * @param {string} opts.url - required, URL of the file; must answer HEAD with an etag
     * @param {number} [opts.minPreloadTime] - optional, seconds before triggering a background update check
     * @returns {Promise<FileRef>} the FileRef of the up-to-date file (basically to get the version)
     * @throws {TypeError} If name or url is missing.
     * @throws {Error} If the HEAD or GET request fails, or the HEAD response has no etag.
     *   When the server answered, the error has `status` (the HTTP status code).
    */
    async update(opts) {
        if (!opts || !opts.name) throw new TypeError('downloader: opts.name is required')
        if (!opts.url) throw new TypeError('downloader: opts.url is required')

        const current = await fs.get({ name: opts.name })

        if (opts.minPreloadTime === undefined) {
            return this._startDownload(opts)
        }

        const lastChecked = current
            ? new Date(current.checkedAt || current.version).getTime()
            : 0
        const age = current ? (Date.now() - lastChecked) / 1000 : Infinity

        if (!current || age >= opts.minPreloadTime) {
            const promise = this._startDownload(opts)
            if (!current) return promise
        }

        return current
    }

    /**
     * Starts an in-flight download or re-uses an existing one.
     * @param {Object} opts
     * @returns {Promise<FileRef>}
     * @private
     */
    _startDownload(opts) {
        if (this.downloads.has(opts.name)) {
            return this.downloads.get(opts.name)
        }

        const promise = this._fetchAndUpdate(opts).finally(() => {
            this.downloads.delete(opts.name)
        })

        // Catch background errors so detached background promises do not trigger unhandledrejection
        promise.catch(() => {})

        this.downloads.set(opts.name, promise)
        return promise
    }

    /**
     * @param {Object} opts
     * @param {string} opts.name
     * @param {string} opts.url
     * @returns {Promise<FileRef>}
     * @private
     */
    async _fetchAndUpdate(opts) {
        // HEAD (gcs) to get the current etag
        const head = await fetch(opts.url, { method: 'HEAD' })
        if (!head.ok) throw httpError(`downloader: HEAD ${opts.url} failed (${head.status})`, head.status)
            
        const etag = head.headers.get('etag')
        if (!etag) throw new Error('downloader: no etag in HEAD response')

        // Same etag as the released version, nothing to download
        const current = await fs.get({ name: opts.name })
        if (current && current.etag === etag) {
            const checkedAt = await fs.touch({ name: opts.name })
            return { ...current, checkedAt }
        }

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
