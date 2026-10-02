// Runs the stages of infer.js one by one and prints what happens. Served with:
//   cd sdkweblegacy && python -m http.server 8000   ->   http://localhost:8000/src/infer/demo/index.html
import { downloader } from '../../internal/downloader/downloader.js'
import { fs } from '../../internal/fs/fs.js'
import { sqlitez } from '../../internal/sqlitez/sqlitez.js'
import { setLog, downloadAssets, loadLiteRt, loadModel } from '../infer.js'

const assets = {
  sqlite: new URL('../gen/modeldb.sqlite', import.meta.url).href,
  tflite: new URL('../gen/model.tflite', import.meta.url).href,
  // litertWasm: new URL('../../../node_modules/@litertjs/core/wasm/', import.meta.url).href,
}

const out = document.getElementById('log')
const t0 = performance.now()

function log(msg, data) {
  const time = ((performance.now() - t0) / 1000).toFixed(2).padStart(6)
  out.insertAdjacentHTML('beforeend', `<div><span class="t">${time}s</span> ${escape(msg)}${data === undefined ? '' : ' ' + escape(JSON.stringify(data))}</div>`)
  console.log(msg, data)
}
const stage = (msg) => out.insertAdjacentHTML('beforeend', `<div class="stage">== ${escape(msg)}</div>`)
const escape = (s) => String(s).replace(/[&<]/g, (c) => ({ '&': '&amp;', '<': '&lt;' }[c]))

async function run() {
  try {
    stage('1. model files and runtimes')
    setLog(log)
    const { dbBytes, modelBytes } = await downloadAssets('fbt', assets)
    const litert = await loadLiteRt(assets.litertWasm)
    const compiled = await loadModel(litert, modelBytes)
    log('step 1 done', { dbBytes: dbBytes.length, modelBytes: modelBytes.length })
  } catch (e) {
    out.insertAdjacentHTML('beforeend', `<div class="err">${escape(e.name)} [${e.code}] ${escape(e.message)}</div>`)
    console.error(e)
  }
}

// Step A: exercises every method of sqlitez on the model sqlite
async function testSqlitez() {
  try {
    stage('A. sqlitez')
    const ref = await downloader.update({ name: 'infer-fbt-sqlite', url: assets.sqlite })
    const path = await fs.resolve(ref)
    log('opening', path)
    const conn = await sqlitez.open({ path })
    try {
      log('tables', await conn.query({ query: "SELECT name FROM sqlite_master WHERE type = 'table' ORDER BY name" }))
      log('count(product)', await conn.query({ query: 'SELECT COUNT(*) AS n FROM product' }))
      const [first] = await conn.query({ query: 'SELECT sku_id, product_id, rec_out FROM product LIMIT 1' })
      log('first product', first)
      log('bind (?) by product_id', await conn.query({ query: 'SELECT sku_id, rec_out FROM product WHERE product_id = ?', params: [first.product_id] }))
      log('query with no rows', await conn.query({ query: 'SELECT * FROM product WHERE sku_id = ?', params: ['nope'] }))

      // update changes the memory copy only: build a namedqueries table (the test sqlite has none)
      await conn.update({ query: 'CREATE TABLE namedqueries (name TEXT PRIMARY KEY, query TEXT)' })
      await conn.update({
        query: 'INSERT INTO namedqueries (name, query) VALUES (?, ?)',
        params: ['product_by_id', 'SELECT sku_id, rec_out FROM product WHERE product_id = ?'],
      })
      log('namedQuery product_by_id', await conn.namedQuery({ name: 'product_by_id', params: [first.product_id] }))
      await conn.namedQuery({ name: 'missing' }).catch((e) => log(`namedQuery missing -> ${e.name} [${e.code}] ${e.message}`))
      await conn.query({ query: 'SELEC nonsense' }).catch((e) => log(`bad sql -> ${e.name} [${e.code}] ${e.message}`))
      await conn.query({}).catch((e) => log(`no query -> ${e.name} ${e.message}`))
    } finally {
      await conn.close()
      log('closed')
    }
    await conn.query({ query: 'SELECT 1' }).catch((e) => log(`after close -> ${e.name} [${e.code}] ${e.message}`))

    // the file was not touched: a new connection has no namedqueries table
    const again = await sqlitez.open({ path })
    try {
      await again.namedQuery({ name: 'product_by_id' }).catch((e) => log(`new conn, update was memory only -> [${e.code}] ${e.message}`))
    } finally {
      await again.close()
    }
    await sqlitez.open({ path: 'i6/fs/nothing/blob.bin' }).catch((e) => log(`open missing file -> [${e.code}] ${e.message}`))
    log('step A done')
  } catch (e) {
    out.insertAdjacentHTML('beforeend', `<div class="err">${escape(e.name)} [${e.code}] ${escape(e.message)}</div>`)
    console.error(e)
  }
}

document.getElementById('run').onclick = run
document.getElementById('sqlitez').onclick = testSqlitez
document.getElementById('clear').onclick = async () => {
  await (await navigator.storage.getDirectory()).removeEntry('i6', { recursive: true }).catch(() => {})
  log('opfs cache removed')
}
