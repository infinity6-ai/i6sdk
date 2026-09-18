// função que o cara passa subscription, store e um arquivo (form)

class Ingestz {
  async upload(subscription, store, inputfiles) {
    console.log("Subscription", subscription)
    console.log("Store", store)
    console.log("InputFiles", inputfiles)

    const fileInput = inputfiles[0]

    console.log("FileInput", fileInput)

    // const fileInput = document.getElementById('fileInput');
    // const progressBar = document.getElementById('progressBar');
    // const statusText = document.getElementById('statusText');

    // if (fileInput.files.length === 0) return alert("Select a file first");
    const file = fileInput.files[0];
    console.log("File", file)

    // // The 'x' variable from your Go backend
    const signedPutUrl = "https://storage.googleapis.com/i6-poc/test-sdk/bla.bin?x-goog-signature=4e3b48f7f467b9ca77884daef746d68e0286a9906542bc7341ce9a0a0a6a0097f0f9be8c491a9648245d6b090f1dc3a7c0c0eb42d670f8376ccebde10e6b4d6fdc56e9df8a5553a7856055fd949c52fb63584dda733c82841051c5af1fc0b0630fb4ae838a5bdedd75e3bb86b7f0bb151310b47bf2589d0f743d573efcc42b4da385cb4091ace3ed0cb8ca52aa5293e8c6bee72a5ae3a07dbc5ea54e5ca6c3cc7a80f2b3c03e904e106f51f626ba45773458b77325f236d8f10b6e6b12495805acd26dc7fc65e3dd71ee2276e1127ebba647d864a3f880d77dd3bdf97e4a72f2346fcd166d8c78a9d4cfb080e16e004f94f41a930130ff2ab01d0a6af91f03b3&x-goog-algorithm=GOOG4-RSA-SHA256&x-goog-credential=testing%40i6-core-prod.iam.gserviceaccount.com%2F20260918%2Fus%2Fstorage%2Fgoog4_request&x-goog-date=20260918T143738Z&x-goog-expires=604800&x-goog-signedheaders=host";

    // // Reset and show progress UI
    // progressBar.style.display = 'block';
    // progressBar.value = 0;
    // statusText.innerText = "Starting upload...";

    const xhr = new XMLHttpRequest();

    // // 1. Listen to the upload progress event
    // xhr.upload.addEventListener('progress', (event) => {
    //   if (event.lengthComputable) {
    //     const percentComplete = Math.round((event.loaded / event.total) * 100);
    //     progressBar.value = percentComplete;
    //     statusText.innerText = `Uploading: ${percentComplete}%`;
    //   }
    // });

    // // 2. Handle completion or failure
    xhr.addEventListener('load', () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        console.log("Upload complete!");
      } else {
        console.log(`Upload failed with status ${xhr.status}`);
      }
    });

    xhr.addEventListener('error', () => {
      console.log("Network error occurred during upload.");
    });

    // // 3. Open and send the raw file data
    xhr.open('PUT', signedPutUrl, true);

    // // Note: Only set Content-Type here if you also specified 
    // // ContentType: "..." in your Go SignedURLOptions.
    // xhr.setRequestHeader('Content-Type', file.type); 

    xhr.send(file);
  }
}

export { Ingestz };
