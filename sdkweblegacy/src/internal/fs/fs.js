/**
 * @fileoverview Versioned file store on the browser's origin private file system (OPFS).
 *
 * Layout under `i6/fs`:
 *   {name}/version.txt          the released (current) version, if any
 *   {name}/v{version}/etag.txt  the etag the version was created for
 *   {name}/v{version}/blob.bin  the content (written by the caller, see `resolve`)
 *
 * Versions are ISO timestamps, so they sort in creation order.
 */

/** OPFS path of the store, relative to the origin's root directory. */
const base = 'i6/fs'

/**
 * Reference to one version of a stored file.
 * @typedef {Object} FileRef
 * @property {string} name - File name, a single OPFS directory name (e.g. "model").
 * @property {string} version - Version id, an ISO timestamp (e.g. "2026-09-29T12:00:00.000Z").
 * @property {string} etag - Etag of the remote content this version was created for.
 */

/**
 * @param {FileSystemDirectoryHandle} dir - Directory holding the file.
 * @param {string} fileName - Name of the text file.
 * @returns {Promise<string>} The file content.
 */
const readText = async (dir, fileName) => (await (await dir.getFileHandle(fileName)).getFile()).text()

/**
 * @param {FileSystemDirectoryHandle} dir - Directory holding the file.
 * @param {string} fileName - Name of the file, created or replaced.
 * @param {Blob|BufferSource|string} data - The content.
 * @returns {Promise<void>}
 */
const writeFile = async (dir, fileName, data) => {
  const writable = await (await dir.getFileHandle(fileName, { create: true })).createWritable()
  await writable.write(data)
  await writable.close()
}

/**
 * @returns {Promise<FileSystemDirectoryHandle>} The base directory.
 * @throws {DOMException} NotFoundError if it does not exist.
 */
const baseDir = async () => {
  let dir = await navigator.storage.getDirectory()
  for (const segment of base.split('/')) {
    dir = await dir.getDirectoryHandle(segment)
  }
  return dir
}

/**
 * Versioned file store. Use the `fs` singleton.
 *
 * A file is created as a new unreleased version, filled, then released to become the
 * current one; releasing removes the older versions.
 */
class FS {
  /**
   * Removes the versions older than version.txt, and the whole file if it is left empty.
   * Without version.txt no version is removed.
   * @param {Object} opts
   * @param {string} opts.name - required
   * @returns {Promise<void>}
   */
  async cleanFile(opts) {
    if (!opts || !opts.name) throw new TypeError('fs: opts.name is required')

    try {
      const parent = await baseDir()
      const dir = await parent.getDirectoryHandle(opts.name)
      const released = await readText(dir, 'version.txt').then((v) => 'v' + v, () => null)

      for await (const [entryName, entry] of dir.entries()) {
        if (released && entry.kind === 'directory' && entryName < released) {
          await dir.removeEntry(entryName, { recursive: true })
        }
      }

      for await (const _ of dir.entries()) return

      await parent.removeEntry(opts.name)
    } catch (e) {
      if (e.name !== 'NotFoundError') throw e
    }
  }

  /**
   * Cleans all files, see cleanFile.
   * @returns {Promise<void>}
   */
  async clean() {
    try {
      const names = []
      for await (const [entryName, entry] of (await baseDir()).entries()) {
        if (entry.kind === 'directory') names.push(entryName)
      }
      for (const name of names) {
        await this.cleanFile({ name })
      }
    } catch (e) {
      if (e.name !== 'NotFoundError') throw e
    }
  }

  /**
   * Resolves a version of a file.
   * @param {Object} opts
   * @param {string} opts.name - required
   * @param {string} [opts.version] - defaults to the content of version.txt
   * @returns {Promise<FileRef|null>} a new full FileRef, or null if it does not exist.
   */
  async get(opts) {
    if (!opts || !opts.name) throw new TypeError('fs: opts.name is required')

    try {
      const dir = await (await baseDir()).getDirectoryHandle(opts.name)

      const version = opts.version || await readText(dir, 'version.txt')
      const versionDir = await dir.getDirectoryHandle('v' + version)

      return {
        name: opts.name,
        version,
        etag: await readText(versionDir, 'etag.txt')
      }
    } catch (e) {
      if (e.name === 'NotFoundError') return null
      throw e
    }
  }

  /**
   * Creates a new version, unless the latest version already has this etag.
   * The new version is not released.
   * @param {Object} opts
   * @param {string} opts.name - required
   * @param {string} opts.etag - required
   * @param {Blob|BufferSource|string} [opts.data] - content, required when a new version is created
   * @returns {Promise<FileRef>} a new full FileRef.
   */
  async create(opts) {
    if (!opts || !opts.name) throw new TypeError('fs: opts.name is required')
    if (!opts.etag) throw new TypeError('fs: opts.etag is required')

    let dir = await navigator.storage.getDirectory()

    for (const segment of [...base.split('/'), opts.name]) {
      dir = await dir.getDirectoryHandle(segment, { create: true })
    }

    let latest = null

    for await (const [entryName, entry] of dir.entries()) {
      if (entry.kind === 'directory' && entryName.startsWith('v') && (!latest || entryName > latest)) latest = entryName
    }

    if (latest && await readText(await dir.getDirectoryHandle(latest), 'etag.txt').catch(() => null) === opts.etag) {
      return { name: opts.name, version: latest.substring(1), etag: opts.etag }
    }

    if (opts.data === undefined) throw new TypeError('fs: opts.data is required to create a new version')

    const version = new Date().toISOString()
    const versionDir = await dir.getDirectoryHandle('v' + version, { create: true })

    await writeFile(versionDir, 'etag.txt', opts.etag)

    return { name: opts.name, version, etag: opts.etag }
  }

  /**
   * Returns the bin file path, full (with base). Does not check that it exists, see get.
   * @param {Object} opts
   * @param {string} opts.name - required
   * @param {string} opts.version - required
   * @returns {Promise<string>} e.g. "i6/fs/{name}/v{version}/blob.bin"
   */
  async resolve(opts) {
    if (!opts || !opts.name) throw new TypeError('fs: opts.name is required')
    if (!opts.version) throw new TypeError('fs: opts.version is required')

    return `${base}/${opts.name}/v${opts.version}/blob.bin`
  }

  /**
   * Reads the content (blob.bin) of a version.
   * @param {Object} opts
   * @param {string} opts.name - required
   * @param {string} opts.version - required
   * @returns {Promise<Uint8Array>} the bytes of the file.
   * @throws {DOMException} NotFoundError if the version or its blob does not exist.
   */
  async read(opts) {
    if (!opts || !opts.name) throw new TypeError('fs: opts.name is required')
    if (!opts.version) throw new TypeError('fs: opts.version is required')

    let dir = await navigator.storage.getDirectory()
    for (const segment of [...base.split('/'), opts.name, 'v' + opts.version]) {
      dir = await dir.getDirectoryHandle(segment)
    }
    const file = await (await dir.getFileHandle('blob.bin')).getFile()
    return new Uint8Array(await file.arrayBuffer())
  }

  /**
   * Makes a version the current one (version.txt), then cleans the older versions.
   * @param {Object} opts
   * @param {string} opts.name - required
   * @param {string} opts.version - required, must exist
   * @returns {Promise<FileRef>} a new full FileRef of the released version.
   */
  async release(opts) {
    if (!opts || !opts.name) throw new TypeError('fs: opts.name is required')
    if (!opts.version) throw new TypeError('fs: opts.version is required')

    const ref = await this.get(opts)
    if (!ref) throw new Error(`fs: version ${opts.version} of ${opts.name} does not exist`)

    await writeFile(await (await baseDir()).getDirectoryHandle(opts.name), 'version.txt', opts.version)
    await this.cleanFile(opts)

    return ref
  }

  /**
   * Removes a version that was never released (e.g. its download failed). Does nothing if it does not exist.
   * @param {Object} opts
   * @param {string} opts.name - required
   * @param {string} opts.version - required, must not be the released one
   * @returns {Promise<void>}
   */
  async discard(opts) {
    if (!opts || !opts.name) throw new TypeError('fs: opts.name is required')
    if (!opts.version) throw new TypeError('fs: opts.version is required')

    try {
      const dir = await (await baseDir()).getDirectoryHandle(opts.name)
      if (await readText(dir, 'version.txt').catch(() => null) === opts.version) {
        throw new Error(`fs: version ${opts.version} of ${opts.name} is the released one`)
      }
      await dir.removeEntry('v' + opts.version, { recursive: true })
    } catch (e) {
      if (e.name !== 'NotFoundError') throw e
    }
  }
}

/** Shared FS instance. */
const fs = new FS()

export { fs };

