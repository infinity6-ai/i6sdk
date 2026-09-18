import { I6Sdk } from '../src/index.js'
import { Ingestz } from '../src/ingestz/ingestz.js'

const i6sdk = new I6Sdk()
document.getElementById('output').textContent = i6sdk.hello()

/** @type {Ingestz} */
const ingestz = await i6sdk.ingestz()

$(window).ready(async () => {
    const form = $("form")

    $(form).submit(() => {
        ingestz.upload({
            dataset: "dectest",
            table: "dummy-catalog",
            partitions: { store: "s1" },
            inputfile: form.find("[type=file]")[0],
        })
            .then(() => console.log("Upload succeeded"))
            .catch((err) => console.error("Upload failed", err))
        return false
    })
})
