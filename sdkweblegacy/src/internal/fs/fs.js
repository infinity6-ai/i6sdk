
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
 * @property {string} path - [write it: opfs ref with base].
 * @property {string} etag - [write it].
 */

// never changes incomming opts/fileref, returns awlays a copy

class FS {

  async cleanFile(opts) {
    // opts.name - required

    // clean older versions than version.txt
    // remove entire directory if it is empty
  }

  async clean() {
    // clean all files
  }

  async get(opts) {
    // opts.name - required
    // opts.version - optional

    // resolve version.txt if version is not given
    // return new full fileRef
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

}

const fs = new FS()

export { fs };

