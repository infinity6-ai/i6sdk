((QUnit) => {

    QUnit.module("i6sdk-web dash", function (hooks) {

        const sqlJsWasmUrl = 'https://cdn.jsdelivr.net/npm/sql.js@1.14.2/dist/'
        const dashName = 'product-relevance-fashion-dash-funnel'
        const sqliteUrl = 'https://dash.test/sqlite.bin'
        const base = {
            dashName,
            dataset: 'dectest',
            partitions: { store_id: 's1', version: 'v1' },
            params: { dataset: 'dectest', params: { min_year_month: '2024-01' } },
        }

        let Dash
        let realFetch
        let etag
        let sqliteBytes
        let fetches

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

        /** Builds a dash sqlite whose named query echoes ?1 and a marker. */
        const buildSqlite = async (marker) => {
            const SQL = await (await import('sql.js')).default({ locateFile: (file) => sqlJsWasmUrl + file })
            const db = new SQL.Database()
            db.run('CREATE TABLE namedqueries (name TEXT PRIMARY KEY, query TEXT)')
            const query = `SELECT json_object('marker', '${marker}', 'params', json(?1)) AS ret`
            db.run('INSERT INTO namedqueries VALUES (?, ?)', [dashName, query])
            const bytes = db.export()
            db.close()
            return bytes
        }

        // The sqlite is served by a fetch stub, because a test page has no server giving it an etag
        hooks.beforeEach(async function () {
            ({ Dash } = await import('../src/dash/dash.js'));
            const { fs } = await import('../src/internal/fs/fs.js')
            await fs.clean()
            await (await navigator.storage.getDirectory()).removeEntry('i6', { recursive: true }).catch(() => { })

            etag = '"e1"'
            sqliteBytes = await buildSqlite('first')
            fetches = []
            realFetch = globalThis.fetch
            globalThis.fetch = async (url, init) => {
                if (url !== sqliteUrl) return realFetch(url, init)
                const method = init?.method ?? 'GET'
                fetches.push(method)
                const headers = { etag }
                return method === 'HEAD' ? new Response(null, { headers }) : new Response(sqliteBytes, { headers })
            }
        })

        hooks.afterEach(function () {
            globalThis.fetch = realFetch
        })

        QUnit.test("query validates its opts", async function (assert) {
            const dash = new Dash()
            const ok = { ...base, sqliteUrl }
            await assertRejects(assert, dash.query(), /opts.dashName/)
            for (const key of ['dashName', 'dataset']) {
                await assertRejects(assert, dash.query({ ...ok, [key]: '' }), new RegExp(`opts.${key}`))
            }
            for (const partitions of [undefined, 'x', ['a'], { store_id: 1 }]) {
                await assertRejects(assert, dash.query({ ...ok, partitions }), /opts.partitions/)
            }
        })

        QUnit.test("query runs the named query of the downloaded sqlite", async function (assert) {
            const result = await new Dash().query({ ...base, sqliteUrl })
            assert.deepEqual(result, { marker: 'first', params: base.params })
            assert.deepEqual(fetches, ['HEAD', 'GET'], 'downloaded once')
        })

        QUnit.test("query binds no parameter when there are no params", async function (assert) {
            // the stock query binds ?1, so this one uses a query that binds nothing
            const SQL = await (await import('sql.js')).default({ locateFile: (file) => sqlJsWasmUrl + file })
            const noParams = new SQL.Database()
            noParams.run('CREATE TABLE namedqueries (name TEXT PRIMARY KEY, query TEXT)')
            noParams.run('INSERT INTO namedqueries VALUES (?, ?)', [dashName, "SELECT json_object('marker', 'none') AS ret"])
            sqliteBytes = noParams.export()
            noParams.close()

            const result = await new Dash().query({ ...base, params: undefined, sqliteUrl })
            assert.deepEqual(result, { marker: 'none' })
        })

        QUnit.test("query downloads again only when the etag changes", async function (assert) {
            const dash = new Dash()
            await dash.query({ ...base, sqliteUrl })

            fetches.length = 0
            assert.strictEqual((await dash.query({ ...base, sqliteUrl })).marker, 'first')
            assert.deepEqual(fetches, ['HEAD'], 'same etag, cache hit')

            etag = '"e2"'
            sqliteBytes = await buildSqlite('second')
            fetches.length = 0
            assert.strictEqual((await dash.query({ ...base, sqliteUrl })).marker, 'second')
            assert.deepEqual(fetches, ['HEAD', 'GET'], 'new etag, downloaded')
        })

        QUnit.test("query reports an asset that cannot be downloaded", async function (assert) {
            const opts = { ...base, sqliteUrl: new URL('nope.sqlite', location.href).href }
            await assertRejects(assert, new Dash().query(opts), /nope\.sqlite/)
        })

        QUnit.test("query reports the http status and cause of a missing asset", async function (assert) {
            const opts = { ...base, sqliteUrl: new URL('nope.sqlite', location.href).href }
            try {
                await new Dash().query(opts)
                assert.ok(false, 'should have been rejected')
            } catch (e) {
                assert.strictEqual(e.status, 404)
            }
        })

        QUnit.test("query reports a sqlite without the named query", async function (assert) {
            sqliteBytes = await buildSqlite('x')
            await assertRejects(assert, new Dash().query({ ...base, dashName: 'other', sqliteUrl }), /no named query "other"/)
        })

        QUnit.test("query finds the sqlite with the signed url api", async function (assert) {
            const requests = []
            const optsSeen = []
            const services = { apis: () => ({ dszSignedUrl: async (req, opts) => {
                requests.push(req)
                optsSeen.push(opts)
                return { ok: true, status: 200, parse: () => ({ urls: [{ url: sqliteUrl }] }) }
            } }) }

            const result = await new Dash(services).query(base)
            assert.strictEqual(result.marker, 'first')
            assert.deepEqual(requests, [{
                params: { dataset: 'dectest' },
                payload: { paths: [{ method: 'GET', path: 'i6tb://results.relevance-fashion-store-dash/store_id=s1/sqlite.bin' }] },
            }])

            const serviceOpts = { datasetToken: 'tok-1' }
            await new Dash(services).query({ ...base, serviceOpts })
            assert.deepEqual(optsSeen, [undefined, serviceOpts], 'serviceOpts given to the api call')

            const noSqlite = { apis: () => ({ dszSignedUrl: async () => ({ ok: true, status: 200, parse: () => ({ urls: [] }) }) }) }
            await assertRejects(assert, new Dash(noSqlite).query(base), /no "sqlite.bin" file/)
        })
    })

})(QUnit)
