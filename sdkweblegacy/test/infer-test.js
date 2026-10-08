((QUnit) => {

    QUnit.module("i6sdk-web infer", function (hooks) {

        // Test assets (gitignored): the test with the real model is skipped when they are not served
        const assets = {
            sqliteUrl: new URL('../src/infer/gen/modeldb.sqlite', location.href).href,
            tfliteUrl: new URL('../src/infer/gen/model.tflite', location.href).href,
        }
        const base = { modelName: 'relevance-fashion-fbt', modelVersion: 'v1', storeId: 's1', dataset: 'dectest' }
        const sku_ids = ['9132438814883', '9132177490083']

        /** @type {import('../src/infer/infer.js').Infer} */
        let infer
        let I6Error

        const rejectsWithCode = async (assert, promise, code, message) => {
            try {
                await promise
                assert.ok(false, 'should have been rejected')
            } catch (e) {
                assert.true(e instanceof I6Error, 'is an I6Error')
                assert.strictEqual(e.code, code)
                assert.true(message.test(e.message), `message ${e.message} matches ${message}`)
            }
        }

        hooks.beforeEach(async function () {
            const { Infer } = await import('../src/infer/infer.js');
            ({ I6Error } = await import('../src/errors.js'))
            infer = new Infer()
        })

        QUnit.test("run validates its opts", async function (assert) {
            const ok = { ...base, params: { sku_ids }, ...assets }
            await rejectsWithCode(assert, infer.run(), 'invalid_argument', /opts.modelName/)
            for (const key of Object.keys(base)) {
                await rejectsWithCode(assert, infer.run({ ...ok, [key]: '' }), 'invalid_argument', new RegExp(`opts.${key}`))
            }
        })

        QUnit.test("run reports an asset that cannot be downloaded", async function (assert) {
            const opts = { ...base, params: {}, sqliteUrl: new URL('nope.sqlite', location.href).href, tfliteUrl: assets.tfliteUrl }
            await rejectsWithCode(assert, infer.run(opts), 'asset', /could not get/)
        })

        QUnit.test("run with the real model", async function (assert) {
            if (!(await fetch(assets.tfliteUrl, { method: 'HEAD' })).ok) return assert.ok(true, 'skipped: no assets in src/infer/gen')

            const opts = { ...base, params: { sku_ids }, ...assets }
            const result = await infer.run(opts)

            assert.strictEqual(result.productIds.length, 20)
            assert.strictEqual(result.top.length, 20)
            assert.deepEqual(result.productIds, result.top.map((t) => t.product_id))
            assert.true(result.top.every((t, i) => i === 0 || result.top[i - 1].score >= t.score), 'best score first')
            assert.strictEqual(result.top[0].sig_index, 4611, 'best signature of this cart (from the first port run)')
        })

        QUnit.test("run finds the model files with the model api", async function (assert) {
            if (!(await fetch(assets.tfliteUrl, { method: 'HEAD' })).ok) return assert.ok(true, 'skipped: no assets in src/infer/gen')

            const requests = []
            const services = { apis: () => ({ modelzSignedGet: async (req) => {
                requests.push(req)
                return { ok: true, payload: { urls: [{ name: 'tflite.bin', url: assets.tfliteUrl }, { name: 'sqlite.bin', url: assets.sqliteUrl }] } }
            } }) }
            const { Infer } = await import('../src/infer/infer.js')

            const result = await new Infer(services).run({ ...base, params: { sku_ids } })
            assert.strictEqual(result.productIds.length, 20)
            assert.deepEqual(requests, [{
                params: { dataset: 'dectest', model_name: 'relevance-fashion-fbt' },
                payload: { model_version: 'v1', model_partitions: { store_id: 's1' } },
            }])

            const noSqlite = { apis: () => ({ modelzSignedGet: async () => ({ ok: true, payload: { urls: [] } }) }) }
            await rejectsWithCode(assert, new Infer(noSqlite).run({ ...base, params: { sku_ids } }), 'asset', /no "sqlite.bin" file/)
        })

        QUnit.test("runModel reports inputs that do not fit the model", async function (assert) {
            if (!(await fetch(assets.tfliteUrl, { method: 'HEAD' })).ok) return assert.ok(true, 'skipped: no assets in src/infer/gen')

            const { loadModel, runModel } = await import('../src/infer/litert.js')
            const { downloader } = await import('../src/internal/downloader/downloader.js')
            const ref = await downloader.update({ name: 'infer-relevance-fashion-fbt-tflite', url: assets.tfliteUrl })
            const model = await loadModel({ ref })
            const fields = { data: [[1, 2, 3, 4, 5, 6, 7]], fill: 0 }
            const padMask = { data: [1], fill: 0 }

            await rejectsWithCode(assert, runModel(model, {}), 'invalid_argument', /is missing/)
            await rejectsWithCode(assert, runModel(model, { fields: { data: [1, 2, 3] }, pad_mask: padMask }), 'invalid_argument', /dimensions/)
            await rejectsWithCode(assert, runModel(model, { fields: { data: [[1, 2, 3, 4, 5, 6, 7]] }, pad_mask: padMask }), 'invalid_argument', /needs/)

            const outputs = await runModel(model, { fields, pad_mask: padMask })
            assert.strictEqual(outputs.signature_probs.length, 7185)
        })
    })

})(QUnit)
