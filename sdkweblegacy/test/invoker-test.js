((QUnit) => {

    QUnit.module("i6sdk-web invokeResult", function (hooks) {

        let Invoker
        let realFetch
        let invoker

        hooks.beforeEach(async function () {
            ({ Invoker } = await import('../src/invoker/invoker.js'))
            realFetch = globalThis.fetch
            invoker = new Invoker({ config: () => ({ baseUrl: 'https://api.test' }) })
        })

        hooks.afterEach(function () {
            globalThis.fetch = realFetch
        })

        const answer = (body, status) => { globalThis.fetch = async () => new Response(body, { status }) }

        QUnit.test("success puts the json in payload, or in body when asked", async function (assert) {
            answer('{"id":"a"}', 200)
            assert.deepEqual(await invoker.invokeResult({ path: 'x' }), { status: 200, ok: true, payload: { id: 'a' } })
            assert.deepEqual(await invoker.invokeResult({ path: 'x' }, 'body'), { status: 200, ok: true, body: { id: 'a' } })
        })

        QUnit.test("a 204 or an empty success body is ok with no payload", async function (assert) {
            answer(null, 204)
            assert.deepEqual(await invoker.invokeResult({ path: 'x' }, 'body'), { status: 204, ok: true })
            answer('', 200)
            assert.deepEqual(await invoker.invokeResult({ path: 'x' }), { status: 200, ok: true })
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
                assert.deepEqual(r, { status, ok: false, error: { status, code, message: text, cause: text } })
            }
        })

        QUnit.test("an empty error body still has a message", async function (assert) {
            answer('', 400)
            const r = await invoker.invokeResult({ path: 'x' })
            assert.strictEqual(r.error.message, 'HTTP 400')
            assert.strictEqual(r.error.code, 'bad_request')
        })

        QUnit.test("a network failure is returned, not thrown", async function (assert) {
            const cause = new TypeError('Failed to fetch')
            globalThis.fetch = async () => { throw cause }
            const r = await invoker.invokeResult({ path: 'x' })
            assert.false(r.ok)
            assert.strictEqual(r.status, 0)
            assert.deepEqual(r.error, { status: 0, code: 'network', message: 'Failed to fetch', cause })
        })

        QUnit.test("uses custom fetcher from ServiceOptions if provided", async function (assert) {
            globalThis.fetch = async () => new Response('global', { status: 200 })
            const fetcher = async (url, opts) => new Response('{"id":"custom"}', { status: 200 })
            
            const r = await invoker.invokeResult({ path: 'x', opts: { fetcher } })
            assert.deepEqual(r, { status: 200, ok: true, payload: { id: 'custom' } })
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
            assert.deepEqual(r, { status: 200, ok: true, payload: { id: 'prep' } })
            assert.strictEqual(fetchedUrl, 'https://api.test/x?sig=test')
        })
    })

})(QUnit)
