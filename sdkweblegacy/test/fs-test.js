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

    })

})(QUnit)
