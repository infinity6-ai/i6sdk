// Serve: python static/infer-demo/serve.py -> http://localhost:8000/static/infer-demo/index.html
import { I6Sdk } from '../../src/index.js'
import { outputQueryParams } from '../../src/infer/infer.js'
import { downloader } from '../../src/internal/downloader/downloader.js'
import { fs } from '../../src/internal/fs/fs.js'
import { sqlitez } from '../../src/internal/sqlitez/sqlitez.js'
import { describeModel, loadModel, runModel } from '../../src/infer/litert.js'

// infer.js prints what it does when this is set (the log panel below shows it)
window.I6_DEBUG = true

const $ = (id) => document.getElementById(id)
// The urls are empty: the model api gives them (fill them to use local files instead, both)

// What every step leaves for the next one
const state = {}

// --- log ---------------------------------------------------------------------------------------------
const consoleLog = console.log
const log = (...args) => {
  $('log').textContent += args.map((a) => (typeof a === 'string' ? a : JSON.stringify(a))).join(' ') + '\n'
}
console.log = (...args) => { consoleLog(...args); log(...args) }
$('clearlog').onclick = () => { $('log').textContent = '' }

// --- showing values ----------------------------------------------------------------------------------
function preview(v) {
  if (ArrayBuffer.isView(v)) v = Array.from(v)
  if (Array.isArray(v)) {
    const items = !$('full').checked && v.length > 30 ? [...v.slice(0, 10), `... ${v.length - 10} more (${v.length} in total)`] : v
    return items.map(preview)
  }
  if (v && typeof v === 'object') return Object.fromEntries(Object.entries(v).map(([k, x]) => [k, preview(x)]))
  return v
}
const text = (v) => (typeof v === 'string' ? v : JSON.stringify(preview(v), null, 1))

function setStatus(card, status) {
  const badge = card.querySelector('.badge')
  badge.textContent = status
  badge.className = `badge ${status}`
}

// Runs one step: shows its input and output, or the error
async function runStep(n, fn) {
  const card = $(`step${n}`)
  card.querySelector('.code').textContent = fn.toString()
  setStatus(card, 'running')
  card.querySelector('.in').textContent = ''
  card.querySelector('.out').textContent = ''
  const sqlBox = card.querySelector('.sql')
  if (sqlBox) sqlBox.textContent = ''
  log(`\n==== step ${n}`)
  try {
    const { input, output, sql } = await fn()
    if (sql) card.querySelector('.sql').textContent = sql
    card.querySelector('.in').textContent = text(input)
    card.querySelector('.out').textContent = text(output)
    setStatus(card, 'ok')
  } catch (e) {
    card.querySelector('.out').textContent = `${e.name}${e.code ? ` [${e.code}]` : ''}: ${e.message}`
    setStatus(card, 'error')
    log(`ERROR ${e.name} [${e.code}] ${e.message}`)
  }
}

const need = (...keys) => {
  const missing = keys.filter((k) => !(k in state))
  if (missing.length) throw new Error(`run the previous steps first (missing: ${missing.join(', ')})`)
}

// --- what the SDK user passes in ---------------------------------------------------------------------
const config = () => {
  const text = (id) => $(id).value.trim()
  const opts = {
    modelName: text('modelName'),
    modelVersion: text('modelVersion'),
    storeId: text('storeId'),
    dataset: text('dataset'),
    params: { cart_skus: text('sku_ids').split(/[\s,]+/).filter(Boolean).map((sku_id) => ({ sku_id })) },
  }
  if (text('sqliteUrl')) opts.sqliteUrl = text('sqliteUrl')
  if (text('tfliteUrl')) opts.tfliteUrl = text('tfliteUrl')
  return opts
}

// The request of sdk.funcs().productRelevanceFashionFbt, from the config card. The funcs method sends all of `input` as
// the model params, and has no sqliteUrl / tfliteUrl / limit: those fields of the card are not used by step 7.
function funcsRequest() {
  const opts = config()
  return {
    input: { dataset: opts.dataset, model_name: opts.modelName, ...opts.params },
    partitions: { store_id: opts.storeId, version: opts.modelVersion },
  }
}

// The code of the SDK user, written with the values of the config card. Shown at the top of the page and run by step 7
function usage() {
  return `import { I6Sdk } from './i6sdk-web.js'

const result = await new I6Sdk().funcs().productRelevanceFashionFbt(${JSON.stringify(funcsRequest(), null, 2)})`
}
for (const id of ['modelName', 'modelVersion', 'storeId', 'dataset', 'sku_ids', 'sqliteUrl', 'tfliteUrl']) $(id).addEventListener('input', () => { $('usage').textContent = usage() })

// Opens the sqlite of step 1, runs a named query and always closes it. Returns the result and the SQL of the stored query.
async function namedJsonQuery(name, params) {
  const conn = await sqlitez.open({ path: state.dbPath })
  try {
    const [stored] = await conn.query({ query: 'SELECT query FROM namedqueries WHERE name = ?', params: [name] })
    return { ret: await conn.namedJsonQuery({ name, params }), sql: stored?.query.trim() }
  } finally {
    await conn.close()
  }
}


// --- the steps ---------------------------------------------------------------------------------------
const steps = {
  1: () => runStep(1, async () => {
    const opts = config()
    // Like run(): the urls of the model api, unless both are filled
    const urls = await new I6Sdk().infer().assetUrls(opts)
    const output = {}
    const refs = {}
    for (const kind of ['sqlite', 'tflite']) {
      const name = `infer-${opts.modelName}-${opts.storeId}-${opts.modelVersion}-${kind}`
      const before = await fs.get({ name })
      refs[kind] = await downloader.update({ name, url: urls[kind] })
      output[kind] = { result: before && before.version === refs[kind].version ? 'cache hit' : 'downloaded', file: refs[kind] }
    }
    state.dbPath = await fs.resolve(refs.sqlite)
    state.tfliteRef = refs.tflite
    output.sqlite.path = state.dbPath
    return { input: urls, output }
  }),

  2: () => runStep(2, async () => {
    need('tfliteRef')
    state.compiled = await loadModel({ ref: state.tfliteRef })
    return {
      input: { ref: state.tfliteRef },
      output: describeModel(state.compiled),
    }
  }),

  3: () => runStep(3, async () => {
    need('dbPath')
    const { modelName, params } = config()
    const name = `model-${modelName}-input`
    const { ret, sql } = await namedJsonQuery(name, [params])
    state.ctx = ret
    return {
      sql,
      input: { namedQuery: name, params: [params] },
      output: state.ctx,
    }
  }),

  4: () => runStep(4, async () => {
    need('ctx')
    state.inputs = state.ctx
    return { input: { ctx: state.ctx }, output: state.inputs }
  }),

  5: () => runStep(5, async () => {
    need('compiled', 'inputs')
    const outputs = await runModel(state.compiled, state.inputs)
    state.outputs = {}
    const output = {}
    for (const [name, values] of Object.entries(outputs)) {
      state.outputs[name] = Array.from(values)
      output[name] = { length: values.length, values: state.outputs[name] }
    }
    const scores = state.outputs.signature_probs
    if (scores) output.signature_probs.best5 = scores.map((score, i) => ({ i, score })).sort((a, b) => b.score - a.score).slice(0, 5)
    return { input: state.inputs, output }
  }),

  6: () => runStep(6, async () => {
    need('ctx', 'outputs')
    const opts = config()
    const name = `model-${opts.modelName}-output`
    const queryParams = outputQueryParams(state.outputs, opts)
    const { ret, sql } = await namedJsonQuery(name, queryParams)
    state.result = ret
    return { sql, input: { namedQuery: name, params: queryParams }, output: state.result }
  }),

  7: () => runStep(7, async () => {
    const req = funcsRequest()
    const result = await new I6Sdk().funcs().productRelevanceFashionFbt(req)
    return { input: req, output: result }
  }),
}

$('usage').textContent = usage()

for (const button of document.querySelectorAll('button[data-step]')) {
  button.onclick = steps[button.dataset.step]
}

// --- tools -------------------------------------------------------------------------------------------
// Redirects to the login page when nobody is signed in, and comes back here
$('login').onclick = async () => {
  try {
    const user = await new I6Sdk().auth().requireUser()
    if (user) log('logged in as', user)
  } catch (e) {
    log(`ERROR login: ${e.message}`)
  }
}

$('clear').onclick = async () => {
  await (await navigator.storage.getDirectory()).removeEntry('i6', { recursive: true }).catch(() => {})
  for (const k of Object.keys(state)) delete state[k]
  log('\nopfs cache removed, the steps start from the beginning')
}

// Everything in the sqlite: every table (CREATE statement and rows) and every named query
$('dump').onclick = async () => {
  need('dbPath')
  const conn = await sqlitez.open({ path: state.dbPath })
  try {
    const tables = await conn.query({ query: "SELECT name, sql FROM sqlite_master WHERE type = 'table' ORDER BY name" })
    log(`\ndump: ${tables.length} tables: ${tables.map((t) => t.name).join(', ')}`)
    for (const t of tables) {
      const rows = await conn.query({ query: `SELECT * FROM "${t.name}" LIMIT 20` })
      log(`\n-- table ${t.name} (first ${rows.length} rows)\n${t.sql}`)
      log(rows.map((r) => JSON.stringify(r)).join('\n'))
    }
  } finally {
    await conn.close()
  }
}

window.onunhandledrejection = (e) => log(`ERROR ${e.reason.name} [${e.reason.code}] ${e.reason.message}`)
