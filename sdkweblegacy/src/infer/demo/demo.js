// Runs the stages of infer.js and prints what happens.
// Serve: python src/infer/demo/serve.py -> http://localhost:8000/src/infer/demo/index.html
import { downloader } from '../../internal/downloader/downloader.js'
import { fs } from '../../internal/fs/fs.js'
import { sqlitez } from '../../internal/sqlitez/sqlitez.js'
import { Infer } from '../infer.js'

const assets = {
  sqlite: new URL('../gen/modeldb.sqlite', import.meta.url).href,
  tflite: new URL('../gen/model.tflite', import.meta.url).href,
}

// Shows everything sent to console.log on the page too (infer.js logs with console.log)
const out = document.getElementById('log')
const consoleLog = console.log
const log = (...args) => {
  out.textContent += args.map((a) => (typeof a === 'string' ? a : JSON.stringify(a))).join(' ') + '\n'
}
console.log = (...args) => {
  consoleLog(...args)
  log(...args)
}

// Errors of the handlers below end up here
window.onunhandledrejection = (e) => log(`ERROR ${e.reason.name} [${e.reason.code}] ${e.reason.message}`)

const json = (v) => JSON.stringify(v)

// --- sqlitez helpers: print the call (sql, params) and what came back -----------------------------
async function query(conn, sql, params) {
  log(`\nquery  ${sql}  params=${json(params ?? [])}`)
  const rows = await conn.query({ query: sql, params })
  log(`  rows=${json(rows)}`)
  return rows
}

async function update(conn, sql, params) {
  log(`\nupdate ${sql}  params=${json(params ?? [])}`)
  await conn.update({ query: sql, params })
  log('  ok')
}

async function namedQuery(conn, name, params) {
  log(`\nnamedQuery name=${name}  params=${json(params ?? [])}`)
  const rows = await conn.namedQuery({ name, params })
  log(`  rows=${json(rows)}`)
  return rows
}

// Prints everything in the sqlite: every table (its CREATE statement and all its rows) and every named query
async function dump(conn) {
  const tables = await conn.query({ query: "SELECT name, sql FROM sqlite_master WHERE type = 'table' ORDER BY name" })
  log(`\ndump: ${tables.length} tables: ${tables.map((t) => t.name).join(', ')}`)

  for (const t of tables) {
    const rows = await conn.query({ query: `SELECT * FROM "${t.name}"` })
    log(`\n-- table ${t.name} (${rows.length} rows)\n${t.sql}`)
    log(rows.map(json).join('\n'))
  }

  if (tables.some((t) => t.name === 'namedqueries')) {
    log('\n-- named queries')
    for (const q of await conn.query({ query: 'SELECT name, query FROM namedqueries ORDER BY name' })) log(`${q.name}: ${q.query}`)
  } else {
    log('\n-- named queries: this sqlite has no "namedqueries" table')
  }
}

// A call that must fail: prints the error it gets
async function mustFail(what, promise) {
  log(`\n${what}`)
  await promise.then(
    () => log('  did NOT fail'),
    (e) => log(`  ${e.name}${e.code ? ` [${e.code}]` : ''}: ${e.message}`),
  )
}

document.getElementById('run').onclick = async () => {
  const opts = { model: 'fbt', params: { productIds: ['9132438814883', '9132177490083', 'unknown-id', '9132438814883'] }, assets }
  log('\n==== infer.run')
  log('opts', opts)
  await new Infer().run(opts)
}

document.getElementById('sqlitez').onclick = async () => {
  log('\n==== sqlitez')
  const ref = await downloader.update({ name: 'infer-fbt-sqlite', url: assets.sqlite })
  const path = await fs.resolve(ref)
  log(`open path=${path}`)
  const conn = await sqlitez.open({ path })

  try {
    await dump(conn)

    await query(conn, "SELECT name FROM sqlite_master WHERE type = 'table'")
    const [first] = await query(conn, 'SELECT product_id, rec_out FROM product LIMIT 1')
    await query(conn, 'SELECT product_id, rec_out FROM product WHERE product_id = ?', [first.product_id])
    await query(conn, 'SELECT product_id FROM product WHERE product_id = ?', ['no-such-id'])

    await update(conn, 'CREATE TABLE namedqueries (name TEXT, query TEXT)')
    await update(conn, 'INSERT INTO namedqueries VALUES (?, ?)', ['by_id', 'SELECT rec_out FROM product WHERE product_id = ?'])
    await namedQuery(conn, 'by_id', [first.product_id])

    await mustFail('namedQuery name=nope', conn.namedQuery({ name: 'nope' }))
    await mustFail('query SELEC', conn.query({ query: 'SELEC' }))
    await mustFail('query with no opts.query', conn.query({}))
  } finally {
    log('\nclose')
    await conn.close()
  }
  await mustFail('query after close', conn.query({ query: 'SELECT 1' }))
}

document.getElementById('clear').onclick = async () => {
  await (await navigator.storage.getDirectory()).removeEntry('i6', { recursive: true }).catch(() => {})
  log('\nopfs cache removed')
}
