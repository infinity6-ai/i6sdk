/**
 * @fileoverview Ingest client for uploading files to the ingestion service.
 */

/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * Key-value mapping representing dataset table partition values, except the ingest ID.
 * @typedef {Record<string, string>} Partitions
 */

/**
 * Parameters for requesting a signed upload URL from the ingest service.
 * @typedef {Object} GetUrlParams
 * @property {string} dataset - Dataset name.
 * @property {string} table - Inbox table name.
 * @property {Partitions} partitions - Values for the table partitions, except the ingest id one.
 */

/**
 * Parameters for uploading a single file to the ingest service.
 * @typedef {Object} UploadParams
 * @property {string} dataset - Dataset the file is ingested into.
 * @property {string} table - Inbox table name.
 * @property {Partitions} partitions - Values for the table partitions, except the ingest id one.
 * @property {File} file - The single file to upload.
 */

/**
 * Client for handling file uploads to the ingest service.
 */
class Ingest {
  /**
   * Constructs an Ingest client instance.
   * @param {Services} services - The services manager instance.
   */
  constructor(services) {
    /**
     * @private
     * @type {Services}
     */
    this._services = services;
  }

  /**
   * Asks ingestz for a signed upload URL (ingestz-get-url API) and uploads a single file to it.
   * @param {UploadParams} params - Upload parameters.
   * @param {string} params.dataset - Dataset the file is ingested into.
   * @param {string} params.table - Inbox table name.
   * @param {Partitions} params.partitions - Values for the table partitions, except the ingest id one.
   * @param {File} params.file - The single file to upload.
   * @returns {Promise<void>} Resolves on a successful upload, rejects otherwise.
   * @throws {Error} If file is not a single File or if upload fails.
   */
  async upload({ dataset, table, partitions, file }) {
    if (!(file instanceof File)) {
      throw new Error("file must be a single File.");
    }

    const signedPutUrl = await this.#getUrl({ dataset, table, partitions });
    await this.#put(signedPutUrl, file);
  }

  /**
   * Gets one signed PUT URL from the ingestz-get-url API.
   * @private
   * @param {GetUrlParams} params - Parameters for requesting the upload URL.
   * @param {string} params.dataset - Dataset name.
   * @param {string} params.table - Inbox table name.
   * @param {Partitions} params.partitions - Values for the table partitions, except the ingest id one.
   * @returns {Promise<string>} Resolves with the signed PUT URL.
   * @throws {Error} If the request fails or no upload URL is returned.
   */
  async #getUrl({ dataset, table, partitions }) {
    //Esse método chama o que é gerado
    const resp = await  this._services.invoker().invoke({
      method: 'POST',
      path: `api/ingest/get-url/dataset/${encodeURIComponent(dataset)}`,
      json: { table, partitions, amount: 1 },
    })
    if (!resp.ok) {
      throw new Error(`Get upload URL failed with status ${resp.status}`);
    }

    const { uploads } = await resp.json();
    if (!uploads?.[0]) {
      throw new Error("Get upload URL returned no upload URL.");
    }
    return uploads[0];
  }

  /**
   * Uploads a file to a signed PUT URL using XMLHttpRequest.
   * @private
   * @param {string} signedPutUrl - The signed URL to PUT the file to.
   * @param {File} file - The file to upload.
   * @returns {Promise<void>} Resolves on a successful upload, rejects on network or HTTP error.
   */
  #put(signedPutUrl, file) {
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest();

      xhr.addEventListener('load', () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          resolve();
        } else {
          reject(new Error(`Upload failed with status ${xhr.status}`));
        }
      });

      xhr.addEventListener('error', () => {
        reject(new Error("Network error occurred during upload."));
      });

      xhr.open('PUT', signedPutUrl, true);

      xhr.send(file);
    });
  }
}

export { Ingest };

