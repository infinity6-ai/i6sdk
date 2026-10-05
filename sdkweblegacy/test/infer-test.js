((QUnit) => {

    QUnit.module("i6sdk-web infer", function (hooks) {

        // Test assets (gitignored): the test with the real model is skipped when they are not served
        const assets = {
            sqlite: new URL('../src/infer/gen/modeldb.sqlite', location.href).href,
            tflite: new URL('../src/infer/gen/model.tflite', location.href).href,
        }
        const productIds = ['9132438814883', '9132177490083']

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
            const ok = { model: 'fbt', params: { productIds }, assets }
            await rejectsWithCode(assert, infer.run(), 'invalid_argument', /opts.model/)
            await rejectsWithCode(assert, infer.run({ ...ok, model: '' }), 'invalid_argument', /opts.model/)
            await rejectsWithCode(assert, infer.run({ ...ok, assets: { sqlite: assets.sqlite } }), 'invalid_argument', /opts.assets/)
        })

        QUnit.test("run reports an asset that cannot be downloaded", async function (assert) {
            const opts = { model: 'nomodel', params: {}, assets: { sqlite: new URL('nope.sqlite', location.href).href, tflite: assets.tflite } }
            await rejectsWithCode(assert, infer.run(opts), 'asset', /could not get/)
        })

        QUnit.test("run with the real model", async function (assert) {
            if (!(await fetch(assets.tflite, { method: 'HEAD' })).ok) return assert.ok(true, 'skipped: no assets in src/infer/gen')

            const opts = { model: 'fbt', params: { productIds }, assets }
            const result = await infer.run(opts)

            assert.strictEqual(result.productIds.length, 20)
            assert.strictEqual(result.top.length, 20)
            assert.deepEqual(result.productIds, result.top.map((t) => t.product_id))
            assert.true(result.top.every((t, i) => i === 0 || result.top[i - 1].score >= t.score), 'best score first')
            assert.strictEqual(result.top[0].sig_index, 4611, 'best signature of this cart (from the first port run)')
        })

        QUnit.test("runModel reports inputs that do not fit the model", async function (assert) {
            const response = await fetch(assets.tflite)
            if (!response.ok) return assert.ok(true, 'skipped: no assets in src/infer/gen')

            const { loadModel, runModel } = await import('../src/infer/infer.js')
            const model = await loadModel(new Uint8Array(await response.arrayBuffer()))
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
