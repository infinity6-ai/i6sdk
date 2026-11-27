((QUnit) => {

    QUnit.module("i6sdkweb generated APIs", function (hooks) {

        let I6Sdk
        let realFetch
        let lastRequest

        hooks.beforeEach(async function () {
            ({ I6Sdk } = await import('../main/index.js'))
            realFetch = globalThis.fetch
            lastRequest = null
        })

        hooks.afterEach(function () {
            globalThis.fetch = realFetch
        })

        const mockFetch = (responseBody, status = 200) => {
            globalThis.fetch = async (url, init) => {
                lastRequest = { url: url.toString(), init }
                return new Response(JSON.stringify(responseBody), {
                    status,
                    headers: { 'content-type': 'application/json' },
                })
            }
        }

        QUnit.test("dszListDatasets sends params and parses response", async function (assert) {
            mockFetch({
                payload: {
                    datasets: [
                        { id: "ds_fashion", name: "Fashion Dataset", grants: 1 }
                    ]
                }
            })

            const sdk = new I6Sdk({ servUrls: { decsuite: 'https://api.i6.test' } })
            sdk.tokens().setUserToken({
                user_access_token: "mock-user-token-123",
                expires_at: Date.now() + 1000000,
            })

            const resp = await sdk.apis().dszListDatasets({ params: { userId: 'u_100' } })

            const getHeader = (name) => {
                const h = lastRequest.init.headers
                return h.get ? h.get(name) : (h[name] || h[name.toLowerCase()])
            }

            assert.true(resp.ok)
            assert.strictEqual(resp.status, 200)
            const data = resp.parse()
            assert.strictEqual(data.payload.datasets.length, 1)
            assert.strictEqual(data.payload.datasets[0].id, "ds_fashion")
            assert.ok(lastRequest.url.includes("/api/ds/user-id/u_100/datasets"))
            assert.strictEqual(getHeader("Authorization"), "Bearer mock-user-token-123")
        })

        QUnit.test("pipezStartRelevanceFashionRelevanceFashionBatchIngest sends payload and parses id", async function (assert) {
            mockFetch({ payload: { id: "run_batch_999" } })

            const sdk = new I6Sdk({ servUrls: { decsuite: 'https://api.i6.test' } })
            sdk.tokens().setDatasetToken("fashion-ds", {
                access_token: "mock-ds-token-abc",
                dataset: "fashion-ds",
                grants: 8,
                expires_at: Date.now() + 1000000,
            })

            const resp = await sdk.apis().pipezStartRelevanceFashionRelevanceFashionBatchIngest({
                params: { dataset: 'fashion-ds' },
                payload: { ingest_id: 'ing_55', table: 'catalog' }
            })

            const getHeader = (name) => {
                const h = lastRequest.init.headers
                return h.get ? h.get(name) : (h[name] || h[name.toLowerCase()])
            }

            assert.true(resp.ok)
            assert.strictEqual(resp.parse().payload.id, "run_batch_999")
            assert.ok(lastRequest.url.includes("/api/pipe/dataset/fashion-ds/start/relevance_fashion/relevance-fashion-batch-ingest"))
            assert.strictEqual(lastRequest.init.method, "POST")
            assert.strictEqual(getHeader("Authorization"), "Bearer mock-ds-token-abc")
            const body = JSON.parse(lastRequest.init.body)
            assert.deepEqual(body, { ingest_id: 'ing_55', table: 'catalog' })
        })

        QUnit.test("pipezStartRelevanceFashionRelevanceFashionMain interpolates dataset in path", async function (assert) {
            mockFetch({ payload: { id: "run_main_111" } })

            const sdk = new I6Sdk({ servUrls: { decsuite: 'https://api.i6.test' } })
            sdk.tokens().setDatasetToken("fashion-ds", {
                access_token: "mock-ds-token-abc",
                dataset: "fashion-ds",
                grants: 8,
                expires_at: Date.now() + 1000000,
            })

            const resp = await sdk.apis().pipezStartRelevanceFashionRelevanceFashionMain({
                params: { dataset: 'fashion-ds' }
            })

            const getHeader = (name) => {
                const h = lastRequest.init.headers
                return h.get ? h.get(name) : (h[name] || h[name.toLowerCase()])
            }

            assert.true(resp.ok)
            assert.strictEqual(resp.parse().payload.id, "run_main_111")
            assert.ok(lastRequest.url.includes("/api/pipe/dataset/fashion-ds/start/relevance_fashion/relevance-fashion-main"))
            assert.strictEqual(lastRequest.init.method, "POST")
            assert.strictEqual(getHeader("Authorization"), "Bearer mock-ds-token-abc")
        })

    })

})(QUnit)
