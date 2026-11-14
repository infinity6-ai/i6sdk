((QUnit) => {

    QUnit.module("i6sdk-web downloader", function (hooks) {

        const name = `test-dl-${Date.now()}`
        const url = 'https://example.test/file.bin'

        /** @type {import('../src/internal/downloader/downloader.js').downloader} */
        let downloader
        /** @type {import('../src/internal/fs/fs.js').fs} */
        let fs

        const realFetch = window.fetch
        let calls

        /**
         * Replaces fetch by a fake server.
         * @param {Object} server
         * @param {string} server.etag - etag returned by HEAD.
         * @param {string} [server.body] - body returned by GET.
         * @param {number} [server.status] - status of both responses.
         */
        const serve = ({ etag, body = '', status = 200 }) => {
            calls = []
            window.fetch = async (input, init = {}) => {
                const method = init.method || 'GET'
                calls.push(method)
                const headers = etag === undefined ? {} : { etag }
                return new Response(method === 'HEAD' ? null : body, { status, headers })
            }
        }

        /**
         * @param {{name: string, version: string}} ref
         * @returns {Promise<string>} the blob.bin content.
         */
        const readBlob = async (ref) => {
            let dir = await navigator.storage.getDirectory()
            const segments = (await fs.resolve(ref)).split('/')
            for (const segment of segments.slice(0, -1)) {
                dir = await dir.getDirectoryHandle(segment)
            }
            return (await (await dir.getFileHandle('blob.bin')).getFile()).text()
        }

        hooks.beforeEach(async function () {
            ({ downloader } = await import('../src/internal/downloader/downloader.js'));
            ({ fs } = await import('../src/internal/fs/fs.js'))
        })

        hooks.afterEach(async function () {
            window.fetch = realFetch
            try {
                const fsDir = await (await (await navigator.storage.getDirectory()).getDirectoryHandle('i6')).getDirectoryHandle('fs')
                await fsDir.removeEntry(name, { recursive: true })
            } catch (e) { }
        })

        QUnit.test("update requires name and url", async function (assert) {
            await assert.rejects(downloader.update({ url }), /opts.name/)
            await assert.rejects(downloader.update({ name }), /opts.url/)
        })

        QUnit.test("update fails when HEAD fails or has no etag", async function (assert) {
            serve({ etag: 'e1', status: 404 })
            await assert.rejects(downloader.update({ name, url }), /HEAD/)
            serve({ etag: undefined })
            await assert.rejects(downloader.update({ name, url }), /etag/)
            assert.strictEqual(await fs.get({ name }), null, "nothing released")
        })

        QUnit.test("update downloads and releases a new file", async function (assert) {
            serve({ etag: 'e1', body: 'hello' })
            const ref = await downloader.update({ name, url })

            assert.strictEqual(ref.name, name)
            assert.strictEqual(ref.etag, 'e1')
            assert.deepEqual(calls, ['HEAD', 'GET'])
            assert.deepEqual(await fs.get({ name }), ref, "released")
            assert.strictEqual(await readBlob(ref), 'hello')
        })

        QUnit.test("update with the same etag does not download", async function (assert) {
            serve({ etag: 'e1', body: 'hello' })
            const ref1 = await downloader.update({ name, url })

            serve({ etag: 'e1', body: 'changed' })
            const ref2 = await downloader.update({ name, url })

            assert.deepEqual(calls, ['HEAD'])
            assert.strictEqual(ref2.version, ref1.version)
            assert.strictEqual(ref2.etag, ref1.etag)
            assert.true(ref2.checkedAt >= ref1.checkedAt)
            assert.deepEqual(await fs.get({ name }), ref2)
            assert.strictEqual(await readBlob(ref2), 'hello')
        })

        QUnit.test("concurrent update calls deduplicate in flight", async function (assert) {
            serve({ etag: 'e1', body: 'hello' })
            const [ref1, ref2] = await Promise.all([
                downloader.update({ name, url }),
                downloader.update({ name, url })
            ])
            assert.deepEqual(calls, ['HEAD', 'GET'])
            assert.deepEqual(ref1, ref2)
            assert.strictEqual(downloader.downloads.has(name), false)
        })

        QUnit.test("update with another etag downloads a new version", async function (assert) {
            serve({ etag: 'e1', body: 'one' })
            const ref1 = await downloader.update({ name, url })
            await new Promise((resolve) => setTimeout(resolve, 5))

            serve({ etag: 'e2', body: 'two' })
            const ref2 = await downloader.update({ name, url })

            assert.true(ref2.version > ref1.version)
            assert.deepEqual(await fs.get({ name }), ref2, "new one released")
            assert.strictEqual(await fs.get({ name, version: ref1.version }), null, "old one cleaned")
            assert.strictEqual(await readBlob(ref2), 'two')
        })

        QUnit.test("update fails when GET fails and does not release", async function (assert) {
            serve({ etag: 'e1', body: 'one' })
            const okFetch = window.fetch
            window.fetch = async (input, init = {}) => init.method === 'HEAD'
                ? okFetch(input, init)
                : new Response('boom', { status: 500 })

            await assert.rejects(downloader.update({ name, url }), /GET/)
            assert.strictEqual(await fs.get({ name }), null, "nothing released")
        })

        QUnit.test("update retries after a failed download", async function (assert) {
            serve({ etag: 'e1', body: 'one' })
            const okFetch = window.fetch
            window.fetch = async (input, init = {}) => init.method === 'HEAD'
                ? okFetch(input, init)
                : new Response('boom', { status: 500 })
            await assert.rejects(downloader.update({ name, url }))

            serve({ etag: 'e1', body: 'one' })
            const ref = await downloader.update({ name, url })
            assert.deepEqual(await fs.get({ name }), ref)
            assert.strictEqual(await readBlob(ref), 'one')
        })

        QUnit.test("update streams the body without reading it into memory", async function (assert) {
            serve({ etag: 'e1' })
            let pulled = 0
            const okFetch = window.fetch
            window.fetch = async (input, init = {}) => {
                if (init.method === 'HEAD') return okFetch(input, init)
                const body = new ReadableStream({
                    pull(controller) {
                        if (pulled === 3) return controller.close()
                        controller.enqueue(new TextEncoder().encode(`chunk${pulled++};`))
                    }
                })
                return new Response(body, { status: 200 })
            }

            const ref = await downloader.update({ name, url })
            assert.strictEqual(pulled, 3)
            assert.strictEqual(await readBlob(ref), 'chunk0;chunk1;chunk2;')
        })

    })

})(QUnit)
