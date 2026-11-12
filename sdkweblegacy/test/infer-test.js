((QUnit) => {

    QUnit.module("i6sdk-web infer", function (hooks) {

        // Test assets (gitignored): the test with the real model is skipped when they are not served
        const assets = {
            sqliteUrl: new URL('../src/infer/gen/sqlite.bin', location.href).href,
            tfliteUrl: new URL('../src/infer/gen/model.tflite', location.href).href,
        }
        const base = { modelName: 'relevance-fashion-fbt', partitions: { store_id: 's1', version: 'v1' }, dataset: 'dectest', sessionId: 'sess-1', userId: 'user-1' }
        const cart = { skus: [{ sku_id: '150' }, { sku_id: '77' }] }

        /** @type {import('../src/infer/infer.js').Infer} */
        let infer, events

        const assertRejects = async (assert, promise, message) => {
            try {
                await promise
                assert.ok(false, 'should have been rejected')
            } catch (e) {
                assert.true(e instanceof Error, 'is an Error')
                if (message) {
                    assert.true(message.test(e.message), `message ${e.message} matches ${message}`)
                }
            }
        }

        const stubServices = (events, apis = {}) => ({
            config: () => ({}),
            apis: () => ({ ingestzEventRelevanceFashionEvent: async (req) => { events.push(req); return { ok: true, status: 200 } }, ...apis }),
        })

        hooks.beforeEach(async function () {
            const { Infer } = await import('../src/infer/infer.js');
            events = []
            infer = new Infer(stubServices(events))
        })

        QUnit.test("run validates its opts", async function (assert) {
            const ok = { ...base, params: cart, ...assets }
            await assertRejects(assert, infer.run(), /opts.modelName/)
            for (const key of ['modelName', 'dataset', 'sessionId', 'userId']) {
                await assertRejects(assert, infer.run({ ...ok, [key]: '' }), new RegExp(`opts.${key}`))
            }
            for (const partitions of [undefined, 'x', ['a'], { store_id: 1 }]) {
                await assertRejects(assert, infer.run({ ...ok, partitions }), /opts.partitions/)
            }
        })

        QUnit.test("run reports an asset that cannot be downloaded", async function (assert) {
            const opts = { ...base, params: {}, sqliteUrl: new URL('nope.sqlite', location.href).href, tfliteUrl: assets.tfliteUrl }
            await assertRejects(assert, infer.run(opts), /nope\.sqlite/)
        })

        QUnit.test("run with the real model", async function (assert) {
            if (!(await fetch(assets.tfliteUrl, { method: 'HEAD' })).ok) return assert.ok(true, 'skipped: no assets in src/infer/gen')

            const opts = { ...base, params: cart, ...assets }
            const { results, rec_id } = await infer.run({ ...opts, limit: 5 })

            assert.strictEqual(results.length, 5, 'limit is honored')
            assert.strictEqual(events.length, 1, 'run sends the recraw event')
            assert.strictEqual(events[0].payload.recraw.rec_id, rec_id, 'with the rec_id it returns')
            assert.true(results.every((r) => Array.isArray(r.skus) && r.skus.every((k) => typeof k.sku_id === 'string')), 'skus are {sku_id}')
            assert.true(results.every((r, i) => i === 0 || results[i - 1].rating >= r.rating), 'best rating first')
            await assertRejects(assert, infer.run({ ...opts, limit: 0 }), /opts.limit/)
        })

        QUnit.test("run and selectResults give serviceOpts to the api calls", async function (assert) {
            const serviceOpts = { datasetToken: 'tok-1' }
            const calls = []
            const services = stubServices([], {
                dszModelGet: async (req, opts) => { calls.push(['dszModelGet', opts]); return { ok: true, status: 200, parse: () => ({ urls: [] }) } },
                ingestzEventRelevanceFashionEvent: async (req, opts) => { calls.push(['event', opts]); return { ok: true, status: 200 } },
            })
            const { Infer } = await import('../src/infer/infer.js')

            await assertRejects(assert, new Infer(services).run({ ...base, params: cart, serviceOpts }), /no "sqlite.bin" file/)
            const opts = {
                dataset: 'dectest', modelName: 'relevance-fashion-fbt', modelVersion: 'v1', storeId: 's1', req: { input: cart },
                results: { results: [{ skus: [{ sku_id: '1' }], rating: 0.5 }] }, recId: 'rec-1', sessionId: 'sess-1', userId: 'user-1', serviceOpts,
            }
            await new Infer(services).selectResults(opts)
            assert.deepEqual(calls, [['dszModelGet', serviceOpts], ['event', serviceOpts]])
        })

        QUnit.test("run finds the model files with the model api", async function (assert) {
            if (!(await fetch(assets.tfliteUrl, { method: 'HEAD' })).ok) return assert.ok(true, 'skipped: no assets in src/infer/gen')

            const requests = []
            const services = stubServices([], { dszModelGet: async (req) => {
                requests.push(req)
                return { ok: true, status: 200, parse: () => ({ urls: [{ name: 'tflite.bin', url: assets.tfliteUrl }, { name: 'sqlite.bin', url: assets.sqliteUrl }] }) }
            } })
            const { Infer } = await import('../src/infer/infer.js')

            const result = await new Infer(services).run({ ...base, params: cart })
            assert.strictEqual(result.results.length, 20, 'default limit')
            assert.deepEqual(requests, [{
                params: { dataset: 'dectest', model_name: 'relevance-fashion-fbt' },
                payload: { model_version: 'v1', model_partitions: { store_id: 's1', version: 'v1' } },
            }])

            const noSqlite = stubServices([], { dszModelGet: async () => ({ ok: true, status: 200, parse: () => ({ urls: [] }) }) })
            await assertRejects(assert, new Infer(noSqlite).run({ ...base, params: cart }), /no "sqlite.bin" file/)
        })

        QUnit.test("selectResults returns the selection as is and sends its sku ids as recshow", async function (assert) {
            const requests = []
            const services = { config: () => ({}), apis: () => ({ ingestzEventRelevanceFashionEvent: async (req) => {
                requests.push(req)
                return { ok: true, status: 200 }
            } }) }
            const { Infer } = await import('../src/infer/infer.js')
            const req = { input: { dataset: 'dectest' } }
            const results = { results: [{ skus: [{ sku_id: '1' }, { sku_id: '2' }], rating: 0.5 }, { skus: [{ sku_id: '1' }], rating: 0.8 }] }
            const selected = { results: [{ sku_id: '1', enriched: { name: 'shirt' } }] }
            const calls = []
            const opts = {
                dataset: 'dectest', modelName: 'relevance-fashion-fbt', modelVersion: 'v1', storeId: 's1', req, results, recId: 'rec-1',
                sessionId: 'sess-1', userId: 'user-1',
                selectResults: async (...args) => { calls.push(args); return selected },
            }

            assert.deepEqual(await new Infer(services).selectResults(opts),
                selected, 'returned as is')
            assert.deepEqual(calls, [[req, results]], 'selectResults gets req and results')
            assert.strictEqual(requests.length, 1)
            const sent = requests.map((r) => {
                assert.deepEqual(r.params, { dataset: 'dectest' })
                const event = r.payload
                assert.strictEqual(event.store_id, 's1')
                assert.deepEqual([event.session_id, event.user_id], ['sess-1', 'user-1'])
                assert.true(/^[0-9a-f-]{36}$/.test(event.id), 'id is a uuid')
                const kind = event.recraw ? 'recraw' : 'recshow'
                assert.strictEqual(event[kind === 'recraw' ? 'recshow' : 'recraw'], undefined, `${kind} only`)
                return [kind, { ...event[kind], payload: JSON.parse(event[kind].payload) }]
            })
            const rec = (payload) => ({ rec_id: 'rec-1', model_id: 'relevance-fashion-fbt', model_version: 'v1', payload })
            assert.deepEqual(sent, [['recshow', rec({ skus: [{ sku_id: '1' }] })]], 'recshow only, with only the sku ids')

            // Default: per result, the first sku not in the cart nor already picked, with the rating of that result
            const defaults = {
                ...opts, selectResults: undefined, req: { input: { dataset: 'dectest', skus: [{ sku_id: '1' }] } },
                results: { results: [
                    { skus: [{ sku_id: '1' }, { sku_id: '2' }], rating: 0.9 },
                    { skus: [{ sku_id: '2' }, { sku_id: '3' }], rating: 0.7 },
                    { skus: [{ sku_id: '1' }, { sku_id: '3' }], rating: 0.5 },
                ] },
            }
            assert.deepEqual(await new Infer(services).selectResults(defaults),
                { results: [{ sku_id: '2', rating: 0.9 }, { sku_id: '3', rating: 0.7 }] }, 'default skips the cart and duplicates')

            const returning = (ret) => ({ ...opts, selectResults: async () => ret })
            await assertRejects(assert, new Infer(services).selectResults({ ...opts, selectResults: 'x' }), /opts.selectResults/)
            await assertRejects(assert, new Infer(services).selectResults(returning([])), /must return/)
            await assertRejects(assert, new Infer(services).selectResults(returning({ results: [{ sku_id: '2' }, {}] })), /no sku_id at results\[1\]$/)
            await assertRejects(assert, new Infer(services).selectResults({ ...opts, recId: undefined }), /opts.recId/)
            await assertRejects(assert, new Infer(services).selectResults({ ...opts, sessionId: undefined }), /sessionId and userId/)
            await assertRejects(assert, new Infer(services).selectResults({ ...opts, userId: undefined }), /sessionId and userId/)
        })

        QUnit.test("runModel reports inputs that do not fit the model", async function (assert) {
            if (!(await fetch(assets.tfliteUrl, { method: 'HEAD' })).ok) return assert.ok(true, 'skipped: no assets in src/infer/gen')

            const { loadModel, runModel } = await import('../src/infer/litert.js')
            const { downloader } = await import('../src/internal/downloader/downloader.js')
            const ref = await downloader.update({ name: 'infer-relevance-fashion-fbt-tflite', url: assets.tfliteUrl })
            const model = await loadModel({ ref })
            const fields = { data: Array(70).fill(0), shape: [1, 10, 7] }
            const padMask = { data: [1, 0, 0, 0, 0, 0, 0, 0, 0, 0], shape: [1, 10] }

            await assertRejects(assert, runModel(model, {}), /is missing/)
            await assertRejects(assert, runModel(model, { fields: { data: [1, 2, 3] }, pad_mask: padMask }), /no shape/)
            await assertRejects(assert, runModel(model, { fields: { data: [1, 2, 3], shape: [3] }, pad_mask: padMask }), /the model needs/)
            await assertRejects(assert, runModel(model, { fields: { data: Array(69).fill(0), shape: [1, 10, 7] }, pad_mask: padMask }), /69 values/)

            const outputs = await runModel(model, { fields, pad_mask: padMask })
            assert.strictEqual(outputs.signature_probs.length, 7185)
        })
    })

})(QUnit)
