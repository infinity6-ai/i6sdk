class Ingestz {
  /**
   * Uploads a single file to the signed storage URL.
   * @param {Object} params
   * @param {string} params.subscription
   * @param {string} params.store
   * @param {HTMLInputElement} params.inputfile - A single file input element.
   * @returns {Promise<void>} Resolves on a successful upload, rejects otherwise.
   */
  upload({ subscription, store, inputfile }) {
    console.log("Subscription", subscription)
    console.log("Store", store)

    const file = inputfile.files[0];

    const signedPutUrl = "https://storage.googleapis.com/i6-poc/test-sdk/bla.bin?x-goog-signature=4e3b48f7f467b9ca77884daef746d68e0286a9906542bc7341ce9a0a0a6a0097f0f9be8c491a9648245d6b090f1dc3a7c0c0eb42d670f8376ccebde10e6b4d6fdc56e9df8a5553a7856055fd949c52fb63584dda733c82841051c5af1fc0b0630fb4ae838a5bdedd75e3bb86b7f0bb151310b47bf2589d0f743d573efcc42b4da385cb4091ace3ed0cb8ca52aa5293e8c6bee72a5ae3a07dbc5ea54e5ca6c3cc7a80f2b3c03e904e106f51f626ba45773458b77325f236d8f10b6e6b12495805acd26dc7fc65e3dd71ee2276e1127ebba647d864a3f880d77dd3bdf97e4a72f2346fcd166d8c78a9d4cfb080e16e004f94f41a930130ff2ab01d0a6af91f03b3&x-goog-algorithm=GOOG4-RSA-SHA256&x-goog-credential=testing%40i6-core-prod.iam.gserviceaccount.com%2F20260918%2Fus%2Fstorage%2Fgoog4_request&x-goog-date=20260918T143738Z&x-goog-expires=604800&x-goog-signedheaders=host";

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

export { Ingestz };
