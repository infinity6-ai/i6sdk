((QUnit) => {

    QUnit.module("i6sdk-web infer", function (hooks) {

        // Test assets (gitignored): the test with the real model is skipped when they are not served
        const assets = {
            sqliteUrl: new URL('../src/infer/gen/sqlite.bin', location.href).href,
            tfliteUrl: new URL('../src/infer/gen/model.tflite', location.href).href,
        }
        const base = { modelName: 'relevance-fashion-fbt', partitions: { store_id: 's1', version: 'v1' }, dataset: 'dectest' }
        const cart = { cart_skus: [{ sku_id: '150' }, { sku_id: '77' }] }

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
            const ok = { ...base, params: cart, ...assets }
            await rejectsWithCode(assert, infer.run(), 'invalid_argument', /opts.modelName/)
            for (const key of ['modelName', 'dataset']) {
                await rejectsWithCode(assert, infer.run({ ...ok, [key]: '' }), 'invalid_argument', new RegExp(`opts.${key}`))
            }
            for (const partitions of [undefined, 'x', ['a'], { store_id: 1 }]) {
                await rejectsWithCode(assert, infer.run({ ...ok, partitions }), 'invalid_argument', /opts.partitions/)
            }
        })

        QUnit.test("run reports an asset that cannot be downloaded", async function (assert) {
            const opts = { ...base, params: {}, sqliteUrl: new URL('nope.sqlite', location.href).href, tfliteUrl: assets.tfliteUrl }
            await rejectsWithCode(assert, infer.run(opts), 'asset', /could not get/)
        })

        QUnit.test("run with the real model", async function (assert) {
            if (!(await fetch(assets.tfliteUrl, { method: 'HEAD' })).ok) return assert.ok(true, 'skipped: no assets in src/infer/gen')

            const opts = { ...base, params: cart, ...assets }
            const { results } = await infer.run({ ...opts, limit: 5 })

            assert.strictEqual(results.length, 5, 'limit is honored')
            assert.true(results.every((r) => Array.isArray(r.skus) && r.skus.every((k) => typeof k.sku_id === 'string')), 'skus are {sku_id}')
            assert.true(results.every((r, i) => i === 0 || results[i - 1].rating >= r.rating), 'best rating first')
            await rejectsWithCode(assert, infer.run({ ...opts, limit: 0 }), 'invalid_argument', /opts.limit/)
        })

        QUnit.test("run finds the model files with the model api", async function (assert) {
            if (!(await fetch(assets.tfliteUrl, { method: 'HEAD' })).ok) return assert.ok(true, 'skipped: no assets in src/infer/gen')

            const requests = []
            const services = { apis: () => ({ dszModelGet: async (req) => {
                requests.push(req)
                return { ok: true, payload: { urls: [{ name: 'tflite.bin', url: assets.tfliteUrl }, { name: 'sqlite.bin', url: assets.sqliteUrl }] } }
            } }) }
            const { Infer } = await import('../src/infer/infer.js')

            const result = await new Infer(services).run({ ...base, params: cart })
            assert.strictEqual(result.results.length, 20, 'default limit')
            assert.deepEqual(requests, [{
                params: { dataset: 'dectest', model_name: 'relevance-fashion-fbt' },
                payload: { model_version: 'v1', model_partitions: { store_id: 's1', version: 'v1' } },
            }])

            const noSqlite = { apis: () => ({ dszModelGet: async () => ({ ok: true, payload: { urls: [] } }) }) }
            await rejectsWithCode(assert, new Infer(noSqlite).run({ ...base, params: cart }), 'asset', /no "sqlite.bin" file/)
        })

        QUnit.test("runModel reports inputs that do not fit the model", async function (assert) {
            if (!(await fetch(assets.tfliteUrl, { method: 'HEAD' })).ok) return assert.ok(true, 'skipped: no assets in src/infer/gen')

            const { loadModel, runModel } = await import('../src/infer/litert.js')
            const { downloader } = await import('../src/internal/downloader/downloader.js')
            const ref = await downloader.update({ name: 'infer-relevance-fashion-fbt-tflite', url: assets.tfliteUrl })
            const model = await loadModel({ ref })
            const fields = { data: Array(70).fill(0), shape: [1, 10, 7] }
            const padMask = { data: [1, 0, 0, 0, 0, 0, 0, 0, 0, 0], shape: [1, 10] }

            await rejectsWithCode(assert, runModel(model, {}), 'invalid_argument', /is missing/)
            await rejectsWithCode(assert, runModel(model, { fields: { data: [1, 2, 3] }, pad_mask: padMask }), 'invalid_argument', /no shape/)
            await rejectsWithCode(assert, runModel(model, { fields: { data: [1, 2, 3], shape: [3] }, pad_mask: padMask }), 'invalid_argument', /the model needs/)
            await rejectsWithCode(assert, runModel(model, { fields: { data: Array(69).fill(0), shape: [1, 10, 7] }, pad_mask: padMask }), 'invalid_argument', /69 values/)

            const outputs = await runModel(model, { fields, pad_mask: padMask })
            assert.strictEqual(outputs.signature_probs.length, 7185)
        })
    })

})(QUnit)
