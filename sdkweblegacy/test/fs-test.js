((QUnit) => {

    QUnit.module("i6sdk-web fs", function (hooks) {

        const name = `test-${Date.now()}`

        /** @type {import('../src/internal/fs/fs.js').fs} */
        let fs

        /**
         * Writes a text file inside i6/fs/{name}/{path}.
         * @param {string[]} path - Directories, then the file name.
         * @param {string} content - File content.
         * @param {string} [fileName] - File name, defaults to the test one.
         * @returns {Promise<void>}
         */
        const write = async (path, content, fileName = name) => {
            let dir = await navigator.storage.getDirectory()
            for (const segment of ['i6', 'fs', fileName, ...path.slice(0, -1)]) {
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

        QUnit.test("create requires name and etag", async function (assert) {
            await assert.rejects(fs.create({ etag: 'e1' }), /opts.name/)
            await assert.rejects(fs.create({ name }), /opts.etag/)
        })

        QUnit.test("create writes a new version", async function (assert) {
            const ref = await fs.create({ name, etag: 'e1' })
            assert.strictEqual(ref.name, name)
            assert.strictEqual(ref.etag, 'e1')
            assert.strictEqual(new Date(ref.version).toISOString(), ref.version)
            assert.deepEqual(await fs.get({ name, version: ref.version }), ref)
            assert.strictEqual(await fs.get({ name }), null, "not released")

            const versionDir = await (await (await (await (await navigator.storage.getDirectory())
                .getDirectoryHandle('i6')).getDirectoryHandle('fs')).getDirectoryHandle(name)).getDirectoryHandle('v' + ref.version)
            try {
                await versionDir.getFileHandle('blob.bin')
                assert.true(false, "blob.bin is not created")
            } catch (e) {
                assert.strictEqual(e.name, 'NotFoundError', "blob.bin is not created")
            }
        })

        QUnit.test("create with the latest etag returns the latest version", async function (assert) {
            const ref = await fs.create({ name, etag: 'e1' })
            assert.deepEqual(await fs.create({ name, etag: 'e1' }), ref)
        })

        QUnit.test("create with another etag creates a newer version", async function (assert) {
            const ref1 = await fs.create({ name, etag: 'e1' })
            await new Promise((resolve) => setTimeout(resolve, 5))
            const ref2 = await fs.create({ name, etag: 'e2' })
            assert.true(ref2.version > ref1.version)
            assert.deepEqual(await fs.get({ name, version: ref1.version }), ref1, "old kept")
            assert.deepEqual(await fs.create({ name, etag: 'e2' }), ref2, "latest is the new one")
        })

        QUnit.test("release requires name and version", async function (assert) {
            await assert.rejects(fs.release({ version: 'V1' }), /opts.name/)
            await assert.rejects(fs.release({ name }), /opts.version/)
        })

        QUnit.test("release fails if the version does not exist", async function (assert) {
            await assert.rejects(fs.release({ name, version: 'V1' }), /does not exist/)
            await write(['vV1', 'etag.txt'], 'e1')
            await assert.rejects(fs.release({ name, version: 'V2' }), /does not exist/)
            assert.strictEqual(await fs.get({ name }), null, "nothing released")
        })

        QUnit.test("release sets version.txt and cleans older versions", async function (assert) {
            const ref1 = await fs.create({ name, etag: 'e1' })
            await new Promise((resolve) => setTimeout(resolve, 5))
            const ref2 = await fs.create({ name, etag: 'e2' })
            await new Promise((resolve) => setTimeout(resolve, 5))
            const ref3 = await fs.create({ name, etag: 'e3' })

            assert.deepEqual(await fs.release(ref2), ref2)
            assert.deepEqual(await fs.get({ name }), ref2)
            assert.strictEqual(await fs.get({ name, version: ref1.version }), null, "older removed")
            assert.deepEqual(await fs.get({ name, version: ref3.version }), ref3, "newer kept")
        })

        QUnit.test("resolve requires name and version", async function (assert) {
            await assert.rejects(fs.resolve({ version: 'V1' }), /opts.name/)
            await assert.rejects(fs.resolve({ name }), /opts.version/)
        })

        QUnit.test("resolve returns the full blob path", async function (assert) {
            assert.strictEqual(await fs.resolve({ name: 'a', version: 'V1' }), 'i6/fs/a/vV1/blob.bin')

            const ref = await fs.create({ name, etag: 'e1' })
            const path = await fs.resolve(ref)
            assert.strictEqual(path, `i6/fs/${name}/v${ref.version}/blob.bin`)
        })

        QUnit.test("clean cleans all files", async function (assert) {
            const other = name + '-other'
            await write(['version.txt'], 'V2')
            await write(['vV1', 'etag.txt'], 'e1')
            await write(['vV2', 'etag.txt'], 'e2')
            await write(['version.txt'], 'V2', other)
            await write(['vV1', 'etag.txt'], 'e1', other)
            await write(['vV2', 'etag.txt'], 'e2', other)

            await fs.clean()

            for (const n of [name, other]) {
                assert.strictEqual(await fs.get({ name: n, version: 'V1' }), null, `${n}: older removed`)
                assert.deepEqual(await fs.get({ name: n }), { name: n, version: 'V2', etag: 'e2' }, `${n}: released kept`)
            }

            const fsDir = await (await (await navigator.storage.getDirectory()).getDirectoryHandle('i6')).getDirectoryHandle('fs')
            await fsDir.removeEntry(other, { recursive: true })
        })

    })

})(QUnit)
