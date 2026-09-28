import { I6Sdk } from '../src/index.js'

const i6sdk = new I6Sdk()
window.i6sdk = i6sdk
/** @type {Ingest} */
const ingest = await i6sdk.ingest()

$(window).ready(async () => {
    const form = $("form")
    const button = form.find("button")

    $("#login").click(async () => {
        try {
            await (await i6sdk.auth()).requireUser()
        } catch (err) {
            console.error("Login failed", err)
        }
    })

    $(form).submit(() => {
        button.prop("disabled", true)
        ingest.upload({
            dataset: "dectest",
            table: "dummy-catalog",
            partitions: { store: "s1" },
            file: form.find("[type=file]")[0].files[0],
        })
            .then(() => console.log("Upload succeeded"))
            .catch((err) => console.error("Upload failed", err))
            .finally(() => button.prop("disabled", false))
        return false
    })
})
