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
    // const signedPutUrl = "https://storage.googleapis.com/testBucket/bla/bla.txt?..."; 
  
    // // Reset and show progress UI
    // progressBar.style.display = 'block';
    // progressBar.value = 0;
    // statusText.innerText = "Starting upload...";
  
    // const xhr = new XMLHttpRequest();
  
    // // 1. Listen to the upload progress event
    // xhr.upload.addEventListener('progress', (event) => {
    //   if (event.lengthComputable) {
    //     const percentComplete = Math.round((event.loaded / event.total) * 100);
    //     progressBar.value = percentComplete;
    //     statusText.innerText = `Uploading: ${percentComplete}%`;
    //   }
    // });
  
    // // 2. Handle completion or failure
    // xhr.addEventListener('load', () => {
    //   if (xhr.status >= 200 && xhr.status < 300) {
    //     statusText.innerText = "Upload complete!";
    //   } else {
    //     statusText.innerText = `Upload failed with status ${xhr.status}`;
    //   }
    // });
  
    // xhr.addEventListener('error', () => {
    //   statusText.innerText = "Network error occurred during upload.";
    // });
  
    // // 3. Open and send the raw file data
    // xhr.open('PUT', signedPutUrl, true);
    
    // // Note: Only set Content-Type here if you also specified 
    // // ContentType: "..." in your Go SignedURLOptions.
    // // xhr.setRequestHeader('Content-Type', file.type); 
  
    // xhr.send(file);
  }
}

export { Ingestz };
