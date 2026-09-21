import { I6Sdk } from '../src/index.js'
import { Ingest } from '../src/ingest/ingest.js'

const i6sdk = new I6Sdk()
document.getElementById('output').textContent = i6sdk.hello()

/** @type {Ingest} */
const ingest = await i6sdk.ingest()

$(window).ready(async () => {
    const form = $("form")

    $(form).submit(() => {
        ingest.upload({
            dataset: "dectest",
            table: "dummy-catalog",
            partitions: { store: "s1" },
            file: form.find("[type=file]")[0].files[0],
        })
            .then(() => console.log("Upload succeeded"))
            .catch((err) => console.error("Upload failed", err))
        return false
    })
})
