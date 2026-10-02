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

// Runs one demo step, printing the error if it fails
async function step(name, fn) {
  log(`\n== ${name}`)
  try {
    await fn()
  } catch (e) {
    log(`ERROR ${e.name} [${e.code}] ${e.message}`)
  }
}

// Prints the error of a call that is expected to fail
const fails = (what, promise) => promise.then(() => log(`${what}: did NOT fail`), (e) => log(`${what} -> [${e.code}] ${e.message}`))

// Until every stage is implemented, run() stops with an error at the first one that is not
document.getElementById('run').onclick = () => step('infer.run', () => new Infer().run({
  model: 'fbt',
  params: { productIds: ['9132438814883'] },
  assets,
}))

document.getElementById('sqlitez').onclick = () => step('A. sqlitez', async () => {
  const ref = await downloader.update({ name: 'infer-fbt-sqlite', url: assets.sqlite })
  const conn = await sqlitez.open({ path: await fs.resolve(ref) })
  try {
    log('tables', await conn.query({ query: "SELECT name FROM sqlite_master WHERE type = 'table'" }))
    const [first] = await conn.query({ query: 'SELECT product_id, rec_out FROM product LIMIT 1' })
    log('first product', first)
    log('bind', await conn.query({ query: 'SELECT rec_out FROM product WHERE product_id = ?', params: [first.product_id] }))

    // update changes memory only; the test sqlite has no namedqueries table, so we make one
    await conn.update({ query: 'CREATE TABLE namedqueries (name TEXT, query TEXT)' })
    await conn.update({ query: 'INSERT INTO namedqueries VALUES (?, ?)', params: ['by_id', 'SELECT rec_out FROM product WHERE product_id = ?'] })
    log('namedQuery', await conn.namedQuery({ name: 'by_id', params: [first.product_id] }))

    await fails('missing named query', conn.namedQuery({ name: 'nope' }))
    await fails('bad sql', conn.query({ query: 'SELEC' }))
  } finally {
    await conn.close()
  }
  await fails('query after close', conn.query({ query: 'SELECT 1' }))
})

document.getElementById('clear').onclick = async () => {
  await (await navigator.storage.getDirectory()).removeEntry('i6', { recursive: true }).catch(() => {})
  log('\nopfs cache removed')
}
