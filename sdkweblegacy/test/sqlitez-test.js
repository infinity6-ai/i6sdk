((QUnit) => {

    QUnit.module("i6sdk-web sqlitez", function (hooks) {

        const dirName = `test-sqlitez-${Date.now()}`
        const sqlJsWasmUrl = 'https://cdn.jsdelivr.net/npm/sql.js@1.14.2/dist/'

        /** @type {import('../src/internal/sqlitez/sqlitez.js').sqlitez} */
        let sqlitez
        let I6Error
        let conn

        /**
         * Writes a file inside the OPFS directory of the test.
         * @param {string} fileName
         * @param {Uint8Array|string} content
         * @returns {Promise<string>} The OPFS path of the file.
         */
        const write = async (fileName, content) => {
            const dir = await (await navigator.storage.getDirectory()).getDirectoryHandle(dirName, { create: true })
            const writable = await (await dir.getFileHandle(fileName, { create: true })).createWritable()
            await writable.write(content)
            await writable.close()
            return `${dirName}/${fileName}`
        }

        /**
         * Builds a sqlite file with sql.js and writes it to OPFS.
         * @param {string[]} statements - SQL run, in order, to fill the database.
         * @returns {Promise<string>} The OPFS path of the file.
         */
        const writeSqlite = async (statements) => {
            const SQL = await (await import('sql.js')).default({ locateFile: (file) => sqlJsWasmUrl + file })
            const db = new SQL.Database()
            db.run('CREATE TABLE _init (x)') // an empty database exports as 0 bytes, which is not a sqlite file
            for (const statement of statements) db.run(statement)
            const bytes = db.export()
            db.close()
            return write('db.sqlite', bytes)
        }

        /**
         * Asserts that a promise is rejected with an I6Error of the given code.
         * @param {Assert} assert
         * @param {Promise} promise
         * @param {string} code
         * @param {RegExp} [message]
         */
        const rejectsWithCode = async (assert, promise, code, message = /./) => {
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
            ({ sqlitez } = await import('../src/internal/sqlitez/sqlitez.js'));
            ({ I6Error } = await import('../src/errors.js'))
            conn = null
        })

        hooks.afterEach(async function () {
            if (conn) await conn.close()
            try {
                await (await navigator.storage.getDirectory()).removeEntry(dirName, { recursive: true })
            } catch (e) { }
        })

        // ---- open

        QUnit.test("open requires path", async function (assert) {
            await assert.rejects(sqlitez.open(), TypeError)
            await assert.rejects(sqlitez.open({}), /opts.path/)
        })

        QUnit.test("open missing file is an asset error", async function (assert) {
            await rejectsWithCode(assert, sqlitez.open({ path: `${dirName}/nope.sqlite` }), 'asset', /could not read/)
        })

        QUnit.test("open file that is not sqlite is an asset error", async function (assert) {
            const path = await write('db.sqlite', 'this is not a sqlite file at all')
            await rejectsWithCode(assert, sqlitez.open({ path }), 'asset', /not a sqlite file/)
        })

        QUnit.test("open returns a working connection", async function (assert) {
            const path = await writeSqlite(['CREATE TABLE t (id INTEGER PRIMARY KEY, name TEXT)', "INSERT INTO t VALUES (1, 'a')"])
            conn = await sqlitez.open({ path })
            assert.deepEqual(await conn.query({ query: 'SELECT * FROM t' }), [{ id: 1, name: 'a' }])
        })

        // ---- query

        QUnit.test("query requires query", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite([]) })
            await assert.rejects(conn.query(), TypeError)
            await assert.rejects(conn.query({}), /opts.query/)
        })

        QUnit.test("query returns rows as objects, empty array when none", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite([
                'CREATE TABLE t (id INTEGER, name TEXT, score REAL, note TEXT)',
                "INSERT INTO t VALUES (1, 'a', 1.5, NULL), (2, 'b', 2.5, 'x')",
            ]) })
            assert.deepEqual(await conn.query({ query: 'SELECT * FROM t ORDER BY id' }), [
                { id: 1, name: 'a', score: 1.5, note: null },
                { id: 2, name: 'b', score: 2.5, note: 'x' },
            ])
            assert.deepEqual(await conn.query({ query: 'SELECT * FROM t WHERE id = 99' }), [])
        })

        QUnit.test("query binds params in order", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite([
                'CREATE TABLE t (id INTEGER, name TEXT)',
                "INSERT INTO t VALUES (1, 'a'), (2, 'b'), (3, 'c')",
            ]) })
            const rows = await conn.query({ query: 'SELECT id FROM t WHERE id > ? AND name != ? ORDER BY id', params: [1, 'c'] })
            assert.deepEqual(rows, [{ id: 2 }])
        })

        QUnit.test("query params are values, never SQL", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(['CREATE TABLE t (name TEXT)', "INSERT INTO t VALUES ('a')"]) })
            assert.deepEqual(await conn.query({ query: 'SELECT * FROM t WHERE name = ?', params: ["a' OR '1'='1"] }), [])
            assert.deepEqual(await conn.query({ query: 'SELECT count(*) AS n FROM t' }), [{ n: 1 }])
        })

        QUnit.test("query with bad SQL is a query error", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite([]) })
            await rejectsWithCode(assert, conn.query({ query: 'SELEC nothing' }), 'query', /query failed/)
            await rejectsWithCode(assert, conn.query({ query: 'SELECT * FROM missing_table' }), 'query', /no such table/)
        })

        QUnit.test("query with too many params is a query error", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite([]) })
            await rejectsWithCode(assert, conn.query({ query: 'SELECT ?', params: [1, 2] }), 'query')
        })

        QUnit.test("connection is usable after a failed query", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite([]) })
            await assert.rejects(conn.query({ query: 'SELEC nothing' }))
            assert.deepEqual(await conn.query({ query: 'SELECT 1 AS one' }), [{ one: 1 }])
        })

        // ---- update

        QUnit.test("update requires query", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite([]) })
            await assert.rejects(conn.update(), TypeError)
            await assert.rejects(conn.update({}), /opts.query/)
        })

        QUnit.test("update changes the connection and returns nothing", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(['CREATE TABLE t (id INTEGER, name TEXT)']) })
            assert.strictEqual(await conn.update({ query: 'INSERT INTO t VALUES (?, ?)', params: [1, 'a'] }), undefined)
            await conn.update({ query: 'UPDATE t SET name = ? WHERE id = ?', params: ['b', 1] })
            assert.deepEqual(await conn.query({ query: 'SELECT * FROM t' }), [{ id: 1, name: 'b' }])
        })

        QUnit.test("update is memory only, a new connection does not see it", async function (assert) {
            const path = await writeSqlite(['CREATE TABLE t (id INTEGER)'])
            conn = await sqlitez.open({ path })
            await conn.update({ query: 'INSERT INTO t VALUES (1)' })
            await conn.close()

            conn = await sqlitez.open({ path })
            assert.deepEqual(await conn.query({ query: 'SELECT * FROM t' }), [])
        })

        QUnit.test("update with bad SQL is a query error", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite([]) })
            await rejectsWithCode(assert, conn.update({ query: 'INSERT INTO missing_table VALUES (1)' }), 'query', /update failed/)
        })

        // ---- namedQuery

        const namedQueries = (...rows) => [
            'CREATE TABLE namedqueries (name TEXT PRIMARY KEY, query TEXT)',
            'CREATE TABLE t (id INTEGER, name TEXT)',
            "INSERT INTO t VALUES (1, 'a'), (2, 'b')",
            ...rows.map(([name, query]) => `INSERT INTO namedqueries VALUES ('${name}', '${query.replace(/'/g, "''")}')`),
        ]

        QUnit.test("namedQuery requires name", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries()) })
            await assert.rejects(conn.namedQuery(), TypeError)
            await assert.rejects(conn.namedQuery({}), /opts.name/)
        })

        QUnit.test("namedQuery runs the stored query with params", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries(
                ['all', 'SELECT * FROM t ORDER BY id'],
                ['by_id', 'SELECT name FROM t WHERE id = ?'],
            )) })
            assert.deepEqual(await conn.namedQuery({ name: 'all' }), [{ id: 1, name: 'a' }, { id: 2, name: 'b' }])
            assert.deepEqual(await conn.namedQuery({ name: 'by_id', params: [2] }), [{ name: 'b' }])
        })

        QUnit.test("namedQuery with unknown name is a query error", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries(['all', 'SELECT 1'])) })
            await rejectsWithCode(assert, conn.namedQuery({ name: 'nope' }), 'query', /no named query "nope"/)
        })

        QUnit.test("namedQuery without namedqueries table is a query error", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(['CREATE TABLE t (id INTEGER)']) })
            await rejectsWithCode(assert, conn.namedQuery({ name: 'all' }), 'query', /namedqueries/)
        })

        QUnit.test("namedQuery with a broken stored query is a query error", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries(['broken', 'SELEC nothing'])) })
            await rejectsWithCode(assert, conn.namedQuery({ name: 'broken' }), 'query')
        })

        // ---- namedJsonQuery

        QUnit.test("namedJsonQuery requires name", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries()) })
            await assert.rejects(conn.namedJsonQuery(), TypeError)
            await assert.rejects(conn.namedJsonQuery({}), /opts.name/)
            await assert.rejects(conn.namedJsonQuery({ name: 'x', params: 'nope' }), /array/)
        })

        QUnit.test("namedJsonQuery binds params as JSON, in order, and parses ret", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries(
                ['echo', "SELECT json_object('got', json(?1), 'then', json(?2)) AS ret"],
                ['sum', "SELECT json_object('total', sum(value)) AS ret FROM json_each(json_extract(?1, '$.numbers'))"],
            )) })
            assert.deepEqual(await conn.namedJsonQuery({ name: 'echo', params: [{ a: [1, 2], b: 'x' }, 'y'] }), { got: { a: [1, 2], b: 'x' }, then: 'y' })
            assert.deepEqual(await conn.namedJsonQuery({ name: 'sum', params: [{ numbers: [1, 2, 3] }] }), { total: 6 })
        })

        QUnit.test("namedJsonQuery binds a null param as JSON null", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries(
                ['is_null', "SELECT json_object('type', json_type(?)) AS ret"],
            )) })
            assert.deepEqual(await conn.namedJsonQuery({ name: 'is_null', params: [null] }), { type: 'null' })
        })

        QUnit.test("namedJsonQuery returns null when there is no row or ret is null", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries(
                ['no_row', "SELECT '1' AS ret WHERE length(?) < 0"],
                ['null_ret', 'SELECT NULL AS ret WHERE ? IS NOT NULL'],
            )) })
            assert.strictEqual(await conn.namedJsonQuery({ name: 'no_row' }), null)
            assert.strictEqual(await conn.namedJsonQuery({ name: 'null_ret' }), null)
        })

        QUnit.test("namedJsonQuery uses the first row only", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries(
                ['two_rows', "SELECT '1' AS ret WHERE ? IS NOT NULL UNION ALL SELECT '2'"],
            )) })
            assert.strictEqual(await conn.namedJsonQuery({ name: 'two_rows', params: [null] }), 1)
        })

        QUnit.test("namedJsonQuery with unknown name is a query error", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries()) })
            await rejectsWithCode(assert, conn.namedJsonQuery({ name: 'nope' }), 'query', /no named query/)
        })

        QUnit.test("namedJsonQuery without ret column is a query error", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries(['no_ret', "SELECT ? AS other"])) })
            await rejectsWithCode(assert, conn.namedJsonQuery({ name: 'no_ret' }), 'query', /no "ret" column/)
        })

        QUnit.test("namedJsonQuery with invalid JSON in ret is a query error", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries(['bad_json', "SELECT 'not json' || substr(?, 1, 0) AS ret"])) })
            await rejectsWithCode(assert, conn.namedJsonQuery({ name: 'bad_json', params: [null] }), 'query', /not valid JSON/)
        })

        QUnit.test("namedJsonQuery with a stored query without ? for a given param is a query error", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries(['no_bind', "SELECT '1' AS ret"])) })
            await rejectsWithCode(assert, conn.namedJsonQuery({ name: 'no_bind', params: [1] }), 'query', /out of range/)
        })

        // ---- close

        QUnit.test("close twice is fine", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite([]) })
            await conn.close()
            await conn.close()
            assert.true(true)
        })

        QUnit.test("every method fails with a query error after close", async function (assert) {
            conn = await sqlitez.open({ path: await writeSqlite(namedQueries(['all', "SELECT '1' AS ret"])) })
            await conn.close()
            await rejectsWithCode(assert, conn.query({ query: 'SELECT 1' }), 'query', /closed/)
            await rejectsWithCode(assert, conn.update({ query: 'CREATE TABLE x (id INTEGER)' }), 'query', /closed/)
            await rejectsWithCode(assert, conn.namedQuery({ name: 'all' }), 'query', /closed/)
            await rejectsWithCode(assert, conn.namedJsonQuery({ name: 'all' }), 'query', /closed/)
        })

        QUnit.test("connections are independent", async function (assert) {
            const path = await writeSqlite(['CREATE TABLE t (id INTEGER)'])
            conn = await sqlitez.open({ path })
            const other = await sqlitez.open({ path })
            try {
                await conn.update({ query: 'INSERT INTO t VALUES (1)' })
                await conn.close()
                assert.deepEqual(await other.query({ query: 'SELECT * FROM t' }), [])
            } finally {
                await other.close()
            }
        })

    });

})(QUnit)
