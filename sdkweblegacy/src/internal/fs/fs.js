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
 * Gets a directory handle if it exists, or returns null on NotFoundError.
 * @param {FileSystemDirectoryHandle} parent
 * @param {string} name
 * @returns {Promise<FileSystemDirectoryHandle|null>}
 */
const getDirIfExists = async (parent, name) => {
  try {
    return await parent.getDirectoryHandle(name)
  } catch (e) {
    if (e.name === 'NotFoundError') return null
    throw e
  }
}

/**
 * Gets a file handle if it exists, or returns null on NotFoundError.
 * @param {FileSystemDirectoryHandle} parent
 * @param {string} name
 * @returns {Promise<FileSystemFileHandle|null>}
 */
const getFileIfExists = async (parent, name) => {
  try {
    return await parent.getFileHandle(name)
  } catch (e) {
    if (e.name === 'NotFoundError') return null
    throw e
  }
}

/**
 * Reads text of a file if it exists, or returns null.
 * @param {FileSystemDirectoryHandle} dir
 * @param {string} fileName
 * @returns {Promise<string|null>}
 */
const readTextIfExists = async (dir, fileName) => {
  const fileHandle = await getFileIfExists(dir, fileName)
  if (!fileHandle) return null
  return (await fileHandle.getFile()).text()
}

/**
 * Removes an entry if it exists; does nothing if it does not exist.
 * @param {FileSystemDirectoryHandle} dir
 * @param {string} name
 * @param {Object} [opts]
 * @returns {Promise<boolean>} True if removed, false if not found.
 */
const removeEntryIfExists = async (dir, name, opts) => {
  try {
    await dir.removeEntry(name, opts)
    return true
  } catch (e) {
    if (e.name === 'NotFoundError') return false
    throw e
  }
}

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
 * Returns the base directory if it exists, or null.
 * @returns {Promise<FileSystemDirectoryHandle|null>}
 */
const baseDirIfExists = async () => {
  let dir = await navigator.storage.getDirectory()
  for (const segment of base.split('/')) {
    dir = await getDirIfExists(dir, segment)
    if (!dir) return null
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

    const parent = await baseDirIfExists()
    if (!parent) return

    const dir = await getDirIfExists(parent, opts.name)
    if (!dir) return

    const versionText = await readTextIfExists(dir, 'version.txt')
    const released = versionText ? 'v' + versionText : null

    for await (const [entryName, entry] of dir.entries()) {
      if (released && entry.kind === 'directory' && entryName < released) {
        await removeEntryIfExists(dir, entryName, { recursive: true })
      }
    }

    for await (const _ of dir.entries()) return

    await removeEntryIfExists(parent, opts.name)
  }

  /**
   * Cleans all files, see cleanFile.
   * @returns {Promise<void>}
   */
  async clean() {
    const parent = await baseDirIfExists()
    if (!parent) return

    const names = []
    for await (const [entryName, entry] of parent.entries()) {
      if (entry.kind === 'directory') names.push(entryName)
    }
    for (const name of names) {
      await this.cleanFile({ name })
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

    const baseHandle = await baseDirIfExists()
    if (!baseHandle) return null

    const dir = await getDirIfExists(baseHandle, opts.name)
    if (!dir) return null

    const version = opts.version || await readTextIfExists(dir, 'version.txt')
    if (!version) return null

    const versionDir = await getDirIfExists(dir, 'v' + version)
    if (!versionDir) return null

    const etag = await readTextIfExists(versionDir, 'etag.txt')

    return {
      name: opts.name,
      version,
      etag
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

    if (latest) {
      const latestDir = await getDirIfExists(dir, latest)
      if (latestDir && await readTextIfExists(latestDir, 'etag.txt') === opts.etag) {
        return { name: opts.name, version: latest.substring(1), etag: opts.etag }
      }
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

    const baseHandle = await baseDirIfExists()
    if (!baseHandle) return

    const dir = await getDirIfExists(baseHandle, opts.name)
    if (!dir) return

    if (await readTextIfExists(dir, 'version.txt') === opts.version) {
      throw new Error(`fs: version ${opts.version} of ${opts.name} is the released one`)
    }

    await removeEntryIfExists(dir, 'v' + opts.version, { recursive: true })
  }
}

/** Shared FS instance. */
const fs = new FS()

export { fs };
