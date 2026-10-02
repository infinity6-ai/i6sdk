/**
 * @fileoverview Sqlite client for the SDK, running sqlite (sql.js, wasm) in the browser.
 *
 * Usage: get the file with the downloader, open it by its OPFS path and always close the
 * connection when done:
 *
 *   const ref = await downloader.update({ name, url })
 *   const conn = await sqlitez.open({ path: await fs.resolve(ref) })
 *   try {
 *     const rows = await conn.query({ query: 'SELECT * FROM t WHERE id = ?', params: [1] })
 *   } finally {
 *     await conn.close()
 *   }
 *
 * The database lives in memory: `update` changes the memory copy, the file is never written back.
 * Every method takes a single `opts` object with the real parameters inside it.
 */

import { I6Error } from "../../errors.js";

/** Must match the sql.js version in package.json. */
const SQLJS_VERSION = "1.14.2"

/** Where the sql.js wasm file is loaded from. */
const SQLJS_WASM_URL = `https://cdn.jsdelivr.net/npm/sql.js@${SQLJS_VERSION}/dist/`

/** Header every sqlite file starts with. */
const SQLITE_HEADER = "SQLite format 3\0"

/**
 * A value that can be bound to a `?` of a query.
 * @typedef {number|string|null|Uint8Array} SqliteValue
 */

/**
 * A result row: column name -> value.
 * @typedef {Record<string, SqliteValue>} SqliteRow
 */

/**
 * Loaded sql.js module, shared by every connection.
 * @type {Promise<Object>|null}
 */
let sqlJs = null

/**
 * Loads sql.js (and its wasm) once.
 * @returns {Promise<Object>} The initialized sql.js module (`SQL`).
 */
function loadSqlJs() {
    sqlJs ??= import("sql.js").then((m) => (m.default || m)({ locateFile: (file) => SQLJS_WASM_URL + file }))
    return sqlJs
}

/**
 * An open sqlite database. Get one with `sqlitez.open` and always call `close()` in a `finally`.
 */
class SqliteConn {
    /**
     * Use `sqlitez.open`, not this constructor.
     * @param {Object} db - The sql.js database.
     */
    constructor(db) {
        /** @private */
        this.db = db
    }

    /**
     * Closes the connection, freeing the database. Closing twice is fine.
     * @returns {Promise<void>}
     */
    async close() {
        if (!this.db) return
        this.db.close()
        this.db = null
    }

    /**
     * Runs a query and returns its rows.
     * @param {Object} opts
     * @param {string} opts.query - required, the SQL; use `?` for values.
     * @param {SqliteValue[]} [opts.params] - values bound, in order, to the `?` of the query (never concatenate them).
     * @returns {Promise<SqliteRow[]>} The rows (empty array when none), each one an object by column name.
     * @throws {TypeError} If opts.query is missing.
     * @throws {I6Error} With code "query" if the connection is closed or sqlite rejects the query.
     */
    async query(opts) {
        if (!opts || !opts.query) throw new TypeError("sqlitez: opts.query is required")

        const db = this.alive()
        let stmt
        try {
            stmt = db.prepare(opts.query)
            if (opts.params) stmt.bind(opts.params)
            const rows = []
            while (stmt.step()) rows.push(stmt.getAsObject())
            return rows

        } catch (cause) {
            throw new I6Error(`sqlitez: query failed: ${cause.message}`, { code: "query", cause })
        } finally {
            if (stmt) stmt.free()
        }
    }

    /**
     * Runs a statement that returns no data (INSERT, UPDATE, CREATE, ...). Only the memory copy changes.
     * @param {Object} opts
     * @param {string} opts.query - required, the SQL; use `?` for values.
     * @param {SqliteValue[]} [opts.params] - values bound, in order, to the `?` of the query.
     * @returns {Promise<void>}
     * @throws {TypeError} If opts.query is missing.
     * @throws {I6Error} With code "query" if the connection is closed or sqlite rejects the statement.
     */
    async update(opts) {
        if (!opts || !opts.query) throw new TypeError("sqlitez: opts.query is required")

        const db = this.alive()
        try {
            db.run(opts.query, opts.params)
        } catch (cause) {
            throw new I6Error(`sqlitez: update failed: ${cause.message}`, { code: "query", cause })
        }
    }

    /**
     * Runs a query stored in the `namedqueries` table of the sqlite (columns `name` and `query`).
     * @param {Object} opts
     * @param {string} opts.name - required, the `name` of the stored query.
     * @param {SqliteValue[]} [opts.params] - values bound, in order, to the `?` of the stored query.
     * @returns {Promise<SqliteRow[]>} The rows of the stored query, like `query`.
     * @throws {TypeError} If opts.name is missing.
     * @throws {I6Error} With code "query" if the connection is closed, there is no such table or named query, or the query fails.
     */
    async namedQuery(opts) {
        if (!opts || !opts.name) throw new TypeError("sqlitez: opts.name is required")

        const found = await this.query({ query: "SELECT query FROM namedqueries WHERE name = ?", params: [opts.name] })
        if (found.length === 0) {
            throw new I6Error(`sqlitez: no named query "${opts.name}"`, { code: "query" })
        }
        return this.query({ query: String(found[0].query), params: opts.params })
    }

    /**
     * @private
     * @returns {Object} The sql.js database.
     * @throws {I6Error} With code "query" if the connection is closed.
     */
    alive() {
        if (!this.db) throw new I6Error("sqlitez: the connection is closed", { code: "query" })
        return this.db
    }
}

/**
 * Sqlitez client. Use the `sqlitez` singleton.
 */
class Sqlitez {
    /**
     * Opens a sqlite file of the browser file system into memory.
     * @param {Object} opts
     * @param {string} opts.path - required, OPFS path of the file (what `fs.resolve` returns for a downloaded file).
     * @returns {Promise<SqliteConn>} The connection. The caller must `close()` it.
     * @throws {TypeError} If opts.path is missing.
     * @throws {I6Error} With code "asset" if the file cannot be read or is not a sqlite file.
     */
    async open(opts) {
        if (!opts || !opts.path) throw new TypeError("sqlitez: opts.path is required")

        let bytes
        try {
            const segments = opts.path.split("/")
            const fileName = segments.pop()
            let dir = await navigator.storage.getDirectory()
            for (const segment of segments) dir = await dir.getDirectoryHandle(segment)
            const file = await (await dir.getFileHandle(fileName)).getFile()
            bytes = new Uint8Array(await file.arrayBuffer())
        } catch (cause) {
            throw new I6Error(`sqlitez: could not read ${opts.path}: ${cause.message}`, { code: "asset", cause })
        }
        if (new TextDecoder().decode(bytes.subarray(0, 16)) !== SQLITE_HEADER) {
            throw new I6Error(`sqlitez: ${opts.path} is not a sqlite file`, { code: "asset" })
        }

        const SQL = await loadSqlJs()
        return new SqliteConn(new SQL.Database(bytes))
    }
}

/** Shared Sqlitez instance. */
const sqlitez = new Sqlitez()

export { sqlitez };
