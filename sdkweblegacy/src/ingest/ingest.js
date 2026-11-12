/**
 * @fileoverview Ingest client for uploading files to the ingestion service.
 */

/**
 * @typedef {import("../services/services.js").Services} Services
 * @typedef {import("../invoker/invoker.js").ServiceOptions} ServiceOptions
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
 * @property {ServiceOptions} [serviceOpts] - Given to the ingestz-get-url API call.
 */

/**
 * Progress of an upload in flight.
 * @typedef {Object} UploadProgress
 * @property {number} loaded - Bytes sent so far.
 * @property {number} total - Total bytes to send (the file size when the browser can't tell).
 * @property {number} percent - Progress from 0 to 100.
 */

/**
 * Parameters for uploading a single file to the ingest service.
 * @typedef {Object} UploadParams
 * @property {string} dataset - Dataset the file is ingested into.
 * @property {string} table - Inbox table name.
 * @property {Partitions} partitions - Values for the table partitions, except the ingest id one.
 * @property {File} file - The single file to upload.
 * @property {(progress: UploadProgress) => void} [onProgress] - Called as the file is sent.
 * @property {AbortSignal} [signal] - Aborts the upload; the promise then rejects with the signal's reason (an AbortError by default).
 * @property {string} [contentType] - Content-Type of the PUT. Defaults to the file's own type.
 * @property {ServiceOptions} [serviceOpts] - Given to the ingestz-get-url API call (the PUT to the signed URL needs none).
 */

/**
 * Result of a successful upload.
 * @typedef {Object} UploadResult
 * @property {string} ingestId - Id the server assigned to this ingest.
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
   * @returns {Promise<UploadResult>} Resolves with the ingest id on a successful upload.
   */
  async upload({ dataset, table, partitions, file, onProgress, signal, contentType, serviceOpts }) {
    if (!(file instanceof File)) {
      throw new TypeError("file must be a single File.");
    }
    signal?.throwIfAborted();

    const { ingestId, url } = await this.#getUrl({ dataset, table, partitions, serviceOpts });
    signal?.throwIfAborted();

    await this.#put(url, file, { onProgress, signal, contentType });
    return { ingestId };
  }

  /**
   * Gets one signed PUT URL from the ingestz-get-url API.
   * @private
   * @param {GetUrlParams} params - Parameters for requesting the upload URL.
   * @returns {Promise<{ingestId: string, url: string}>} The ingest id and the signed PUT URL.
   */
  async #getUrl({ dataset, table, partitions, serviceOpts }) {
    const resp = await this._services.apis().ingestzGetUrl({
      params: { dataset },
      payload: { table, partitions, amount: 1 },
    }, serviceOpts)
    if (!resp.ok) {
      throw new Error(`Get upload URL failed with status ${resp.status}: ${resp.error?.message ?? ""}`);
    }

    const uploads = resp.payload?.uploads;
    if (!uploads?.[0]) {
      throw new Error("Get upload URL returned no upload URL.");
    }
    return { ingestId: resp.payload.ingest_id, url: uploads[0] };
  }

  /**
   * Uploads a file to a signed PUT URL using XMLHttpRequest (fetch has no upload progress).
   * @private
   * @param {string} signedPutUrl - The signed URL to PUT the file to.
   * @param {File} file - The file to upload.
   * @param {Pick<UploadParams, "onProgress"|"signal"|"contentType">} opts
   * @returns {Promise<void>} Resolves on a successful upload, rejects on network or HTTP error or abort.
   */
  #put(signedPutUrl, file, { onProgress, signal, contentType }) {
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest();
      const abort = () => xhr.abort();
      const abortError = () =>
        signal?.reason ?? new DOMException("Upload aborted.", "AbortError");
      const done = (fn, arg) => {
        signal?.removeEventListener('abort', abort);
        fn(arg);
      };

      xhr.addEventListener('load', () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          done(resolve);
        } else {
          done(reject, new Error(`Upload failed with status ${xhr.status}: ${xhr.responseText}`));
        }
      });

      xhr.addEventListener('error', () => {
        done(reject, new TypeError("Network error occurred during upload."));
      });

      xhr.addEventListener('abort', () => {
        done(reject, abortError());
      });

      if (onProgress) {
        xhr.upload.addEventListener('progress', (e) => {
          const total = e.lengthComputable ? e.total : file.size;
          onProgress({
            loaded: e.loaded,
            total,
            percent: total > 0 ? (e.loaded / total) * 100 : 100,
          });
        });
      }

      xhr.open('PUT', signedPutUrl, true);

      if (contentType) {
        xhr.setRequestHeader('Content-Type', contentType);
      }

      signal?.addEventListener('abort', abort, { once: true });
      xhr.send(file);
    });
  }
}

export { Ingest };
