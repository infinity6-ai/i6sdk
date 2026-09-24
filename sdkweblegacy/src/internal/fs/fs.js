
const base = 'i6/fs'

// file is actually a directory
// {base}/{file}/version.txt
// {base}/{file}/v{versioon}/blob.bin
// {base}/{file}/v{versioon}/etag.txt

// version is current date toISOString

/**
 * [write it].
 * @typedef {Object} FileRef
 * @property {string} name - [write it: bla, ble].
 * @property {string} version - [write it].
 * @property {string} etag - [write it].
 */

// never changes incomming opts/fileref, returns awlays a copy

/**
 * @param {FileSystemDirectoryHandle} dir - Directory holding the file.
 * @param {string} fileName - Name of the text file.
 * @returns {Promise<string>} The file content.
 */
const readText = async (dir, fileName) => (await (await dir.getFileHandle(fileName)).getFile()).text()

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

  async create(opts) {
    // opts.file = "bla" - required
    // opts.etag = "yyy" - required

    // if latest version etag does not match opts.tag, create a new version with this opts.tag

    // return new full FileRef
  }

  async release(opts) {
    // opts.file = "bla" - required
    // opts.version = "version" - required

    // set version.txt to this version
    // async cleanFile(opts)  
  }

  async resolve(opts) {
    // opts.file = "bla" - required
    // opts.version = "version" - required

    // return "bin file path full (with base)"
  }

}

const fs = new FS()

export { fs };

