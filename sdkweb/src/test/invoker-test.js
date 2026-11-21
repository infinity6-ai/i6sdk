((QUnit) => {

    QUnit.module("i6sdk-web invokeResult", function (hooks) {

        let Invoker, Resp
        let realFetch
        let invoker

        hooks.beforeEach(async function () {
            ({ Invoker, Resp } = await import('../main/invoker/invoker.js'))
            realFetch = globalThis.fetch
            invoker = new Invoker({ config: () => ({ baseUrl: 'https://api.test' }), tokens: () => ({ getValidAccessToken: async () => null, getValidDatasetAccessToken: async () => null }) })
        })

        hooks.afterEach(function () {
            globalThis.fetch = realFetch
        })

        const answer = (body, status) => { globalThis.fetch = async () => new Response(body, { status }) }

        QUnit.test("success returns Resp with status, ok and parse()", async function (assert) {
            answer('{"id":"a"}', 200)
            const r = await invoker.invokeResult({ path: 'x' })
            assert.true(r instanceof Resp)
            assert.strictEqual(r.status, 200)
            assert.true(r.ok)
            assert.strictEqual(r.text, '{"id":"a"}')
            assert.deepEqual(r.parse(), { id: 'a' })
        })

        QUnit.test("a 204 or an empty success body returns null on parse()", async function (assert) {
            answer(null, 204)
            const r1 = await invoker.invokeResult({ path: 'x' })
            assert.strictEqual(r1.status, 204)
            assert.true(r1.ok)
            assert.strictEqual(r1.parse(), null)

            answer('', 200)
            const r2 = await invoker.invokeResult({ path: 'x' })
            assert.strictEqual(r2.status, 200)
            assert.true(r2.ok)
            assert.strictEqual(r2.parse(), null)
        })

        QUnit.test("http errors keep the server text and get a code from the status", async function (assert) {
            const cases = [
                [400, 'email_invalid', 'bad_request'],
                [401, 'session expired', 'unauthenticated'],
                [403, 'nope', 'forbidden'],
                [404, 'Not Found', 'not_found'],
                [500, 'boom', 'unavailable'],
                [503, 'down', 'unavailable'],
            ]
            for (const [status, text, code] of cases) {
                answer(text, status)
                const r = await invoker.invokeResult({ path: 'x' })
                assert.strictEqual(r.status, status)
                assert.false(r.ok)
                assert.deepEqual(r.error, { status, code, message: text, cause: text })
                assert.throws(() => r.parse(), new RegExp(text))
            }
        })

        QUnit.test("an empty error body still has a message", async function (assert) {
            answer('', 400)
            const r = await invoker.invokeResult({ path: 'x' })
            assert.strictEqual(r.error.message, 'HTTP 400')
            assert.strictEqual(r.error.code, 'bad_request')
        })

        QUnit.test("a network failure throws", async function (assert) {
            const cause = new TypeError('Failed to fetch')
            globalThis.fetch = async () => { throw cause }
            await assert.rejects(invoker.invokeResult({ path: 'x' }), cause)
        })

        QUnit.test("uses custom fetcher from ServiceOptions if provided", async function (assert) {
            globalThis.fetch = async () => new Response('global', { status: 200 })
            const fetcher = async (url, opts) => new Response('{"id":"custom"}', { status: 200 })
            
            const r = await invoker.invokeResult({ path: 'x', opts: { fetcher } })
            assert.strictEqual(r.status, 200)
            assert.true(r.ok)
            assert.deepEqual(r.parse(), { id: 'custom' })
        })

        QUnit.test("uses config.fetcher when ServiceOptions fetcher is not provided", async function (assert) {
            globalThis.fetch = async () => new Response('global', { status: 200 })
            const configFetcher = async (url, opts) => new Response('{"id":"config-custom"}', { status: 200 })
            const customInvoker = new Invoker({
                config: () => ({ baseUrl: 'https://api.test', fetcher: configFetcher }),
                tokens: () => ({ getValidAccessToken: async () => null, getValidDatasetAccessToken: async () => null })
            })

            const r = await customInvoker.invokeResult({ path: 'x' })
            assert.strictEqual(r.status, 200)
            assert.true(r.ok)
            assert.deepEqual(r.parse(), { id: 'config-custom' })
        })

        QUnit.test("uses config.fetch alias when config.fetcher is not set", async function (assert) {
            globalThis.fetch = async () => new Response('global', { status: 200 })
            const configFetch = async (url, opts) => new Response('{"id":"config-alias"}', { status: 200 })
            const customInvoker = new Invoker({
                config: () => ({ baseUrl: 'https://api.test', fetch: configFetch }),
                tokens: () => ({ getValidAccessToken: async () => null, getValidDatasetAccessToken: async () => null })
            })

            const r = await customInvoker.invokeResult({ path: 'x' })
            assert.strictEqual(r.status, 200)
            assert.true(r.ok)
            assert.deepEqual(r.parse(), { id: 'config-alias' })
        })

        QUnit.test("per-request ServiceOptions fetcher takes precedence over config.fetcher", async function (assert) {
            const configFetcher = async (url, opts) => new Response('{"id":"config-fetcher"}', { status: 200 })
            const reqFetcher = async (url, opts) => new Response('{"id":"req-fetcher"}', { status: 200 })
            const customInvoker = new Invoker({
                config: () => ({ baseUrl: 'https://api.test', fetcher: configFetcher }),
                tokens: () => ({ getValidAccessToken: async () => null, getValidDatasetAccessToken: async () => null })
            })

            const r = await customInvoker.invokeResult({ path: 'x', opts: { fetcher: reqFetcher } })
            assert.strictEqual(r.status, 200)
            assert.true(r.ok)
            assert.deepEqual(r.parse(), { id: 'req-fetcher' })
        })

        QUnit.test("calls prepareRequest from ServiceOptions if provided", async function (assert) {
            let fetchedUrl = ""
            globalThis.fetch = async (url, opts) => {
                fetchedUrl = url
                return new Response('{"id":"prep"}', { status: 200 })
            }
            
            const prepareRequest = async (opts) => {
                opts.url += "?sig=test"
            }
            
            const r = await invoker.invokeResult({ path: 'x', opts: { prepareRequest } })
            assert.strictEqual(r.status, 200)
            assert.true(r.ok)
            assert.deepEqual(r.parse(), { id: 'prep' })
            assert.strictEqual(fetchedUrl, 'https://api.test/x?sig=test')
        })

        QUnit.test("returns 401 when user token is required but missing, and 403 for dataset token", async function (assert) {
            const rUser = await invoker.invokeResult({ path: 'x', opts: { labels: { token_types: 'user' } } })
            assert.strictEqual(rUser.status, 401, "status is 401 for missing user credentials")
            assert.false(rUser.ok)
            assert.strictEqual(rUser.error.code, "unauthenticated")

            const rDataset = await invoker.invokeResult({ path: 'x', opts: { labels: { token_types: 'dataset' } } })
            assert.strictEqual(rDataset.status, 403, "status is 403 for missing dataset credentials")
            assert.false(rDataset.ok)
            assert.strictEqual(rDataset.error.code, "forbidden")
        })
    })

})(QUnit)
