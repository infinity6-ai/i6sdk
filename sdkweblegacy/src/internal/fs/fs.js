
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

class FS {

  async cleanFile(opts) {
    // opts.name - required

    // clean older versions than version.txt
    // remove entire directory if it is empty
  }

  async clean() {
    // clean all files
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
      let dir = await navigator.storage.getDirectory()
      
      for (const segment of [...base.split('/'), opts.name]) {
        dir = await dir.getDirectoryHandle(segment)
      }
      
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

