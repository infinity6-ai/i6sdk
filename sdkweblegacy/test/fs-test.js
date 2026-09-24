((QUnit) => {

    QUnit.module("i6sdk-web fs", function (hooks) {

        const name = `test-${Date.now()}`

        /** @type {import('../src/internal/fs/fs.js').fs} */
        let fs

        /**
         * Writes a text file inside i6/fs/{name}/{path}.
         * @param {string[]} path - Directories, then the file name.
         * @param {string} content - File content.
         * @returns {Promise<void>}
         */
        const write = async (path, content) => {
            let dir = await navigator.storage.getDirectory()
            for (const segment of ['i6', 'fs', name, ...path.slice(0, -1)]) {
                dir = await dir.getDirectoryHandle(segment, { create: true })
            }
            const writable = await (await dir.getFileHandle(path[path.length - 1], { create: true })).createWritable()
            await writable.write(content)
            await writable.close()
        }

        hooks.beforeEach(async function () {
            ({ fs } = await import('../src/internal/fs/fs.js'))
        })

        hooks.afterEach(async function () {
            try {
                const fsDir = await (await (await navigator.storage.getDirectory()).getDirectoryHandle('i6')).getDirectoryHandle('fs')
                await fsDir.removeEntry(name, { recursive: true })
            } catch (e) { }
        })

        QUnit.test("get requires name", async function (assert) {
            await assert.rejects(fs.get({}), /opts.name/)
        })

        QUnit.test("get missing file returns null", async function (assert) {
            assert.strictEqual(await fs.get({ name }), null)
        })

        QUnit.test("get resolves version.txt", async function (assert) {
            await write(['version.txt'], 'V1')
            await write(['vV1', 'etag.txt'], 'e1')
            assert.deepEqual(await fs.get({ name }), { name, version: 'V1', etag: 'e1' })
        })

        QUnit.test("get with version", async function (assert) {
            await write(['version.txt'], 'V1')
            await write(['vV1', 'etag.txt'], 'e1')
            await write(['vV2', 'etag.txt'], 'e2')
            assert.deepEqual(await fs.get({ name, version: 'V2' }), { name, version: 'V2', etag: 'e2' })
            assert.strictEqual(await fs.get({ name, version: 'V3' }), null)
        })

        QUnit.test("get without version.txt returns null", async function (assert) {
            await write(['vV1', 'etag.txt'], 'e1')
            assert.strictEqual(await fs.get({ name }), null)
        })

        QUnit.test("cleanFile requires name", async function (assert) {
            await assert.rejects(fs.cleanFile({}), /opts.name/)
        })

        QUnit.test("cleanFile on missing file or without version.txt does nothing", async function (assert) {
            await fs.cleanFile({ name })
            await write(['vV1', 'etag.txt'], 'e1')
            await fs.cleanFile({ name })
            assert.deepEqual(await fs.get({ name, version: 'V1' }), { name, version: 'V1', etag: 'e1' })
        })

        QUnit.test("cleanFile removes versions older than version.txt", async function (assert) {
            await write(['version.txt'], 'V2')
            await write(['vV1', 'etag.txt'], 'e1')
            await write(['vV2', 'etag.txt'], 'e2')
            await write(['vV3', 'etag.txt'], 'e3')
            await fs.cleanFile({ name })
            assert.strictEqual(await fs.get({ name, version: 'V1' }), null, "older removed")
            assert.deepEqual(await fs.get({ name }), { name, version: 'V2', etag: 'e2' }, "released kept")
            assert.deepEqual(await fs.get({ name, version: 'V3' }), { name, version: 'V3', etag: 'e3' }, "newer kept")
        })

        QUnit.test("cleanFile keeps the directory when it is not empty", async function (assert) {
            await write(['version.txt'], 'V1')
            await write(['vV1', 'etag.txt'], 'e1')
            await fs.cleanFile({ name })
            assert.notStrictEqual(await fs.get({ name }), null)
        })

        QUnit.test("cleanFile removes the directory when it is empty", async function (assert) {
            const fsDir = await (await (await navigator.storage.getDirectory()).getDirectoryHandle('i6')).getDirectoryHandle('fs', { create: true })
            await fsDir.getDirectoryHandle(name, { create: true })
            await fs.cleanFile({ name })
            await assert.rejects(fsDir.getDirectoryHandle(name), (e) => e.name === 'NotFoundError')
        })

    })

})(QUnit)
