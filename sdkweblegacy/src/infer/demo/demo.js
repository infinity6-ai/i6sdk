// Serve: python src/infer/demo/serve.py -> http://localhost:8000/src/infer/demo/index.html
import { I6Sdk } from '../../index.js'
import { downloader } from '../../internal/downloader/downloader.js'
import { fs } from '../../internal/fs/fs.js'
import { sqlitez } from '../../internal/sqlitez/sqlitez.js'
import { describeModel, loadModel, runModel } from '../infer.js'

const $ = (id) => document.getElementById(id)
$('sqlite').value = new URL('../gen/modeldb.sqlite', import.meta.url).href
$('tflite').value = new URL('../gen/model.tflite', import.meta.url).href

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
  const productIds = $('productIds').value.split(/[\s,]+/).filter(Boolean)
  return {
    model: $('model').value.trim(),
    productIds,
    assets: { sqlite: $('sqlite').value.trim(), tflite: $('tflite').value.trim() },
  }
}

// The code of the SDK user, written with the values of the config card. Shown at the top of the page and run by step 7
function usage() {
  const { model, productIds, assets } = config()
  return `import { I6Sdk } from './i6sdk-web.js'

const result = await new I6Sdk().infer().run({
  model: ${JSON.stringify(model)},
  params: { productIds: ${JSON.stringify(productIds)} },
  assets: ${JSON.stringify(assets, null, 4).replace(/\n/g, '\n  ')},
})
// params go to the stored query model-${model}-input, which returns the model inputs as plain rows
// (the SDK pads them to the model's shape and creates the tensors); the model outputs go to
// model-${model}-output, and result is what that query returns.

// Optional: build the inputs yourself instead of the stored query, from what the input query returned:
//   buildInputs: (ctx) => ({ fields: { data: [[...], [...]], fill: 0 }, pad_mask: { data: [1, 1], fill: 0 } })`
}
for (const id of ['model', 'productIds', 'sqlite', 'tflite']) $(id).addEventListener('input', () => { $('usage').textContent = usage() })

// Opens the sqlite of step 1, runs a named query and always closes it. Returns the result and the SQL of the stored query.
async function namedJsonQuery(name, param) {
  const conn = await sqlitez.open({ path: state.dbPath })
  try {
    const [stored] = await conn.query({ query: 'SELECT query FROM namedqueries WHERE name = ?', params: [name] })
    return { ret: await conn.namedJsonQuery({ name, param }), sql: stored?.query.trim() }
  } finally {
    await conn.close()
  }
}


// --- the steps ---------------------------------------------------------------------------------------
const steps = {
  1: () => runStep(1, async () => {
    const { model, assets } = config()
    const output = {}
    const refs = {}
    for (const kind of ['sqlite', 'tflite']) {
      const name = `infer-${model}-${kind}`
      const before = await fs.get({ name })
      refs[kind] = await downloader.update({ name, url: assets[kind] })
      output[kind] = { result: before && before.version === refs[kind].version ? 'cache hit' : 'downloaded', file: refs[kind] }
    }
    state.dbPath = await fs.resolve(refs.sqlite)
    state.modelBytes = await fs.read(refs.tflite)
    output.sqlite.path = state.dbPath
    output.tflite.bytes = state.modelBytes.length
    output.tflite.magic = new TextDecoder().decode(state.modelBytes.subarray(4, 8))
    return { input: { model, assets }, output }
  }),

  2: () => runStep(2, async () => {
    need('modelBytes')
    state.compiled = await loadModel(state.modelBytes)
    return {
      input: { tflite: `${state.modelBytes.length} bytes`, accelerator: 'wasm' },
      output: describeModel(state.compiled),
    }
  }),

  3: () => runStep(3, async () => {
    need('dbPath')
    const { model, productIds } = config()
    const params = { productIds }
    const name = `model-${model}-input`
    const { ret, sql } = await namedJsonQuery(name, params)
    state.ctx = ret
    return {
      sql,
      input: { namedQuery: name, param: params },
      output: { ...state.ctx, 'storeMask: ones': state.ctx.storeMask.filter((m) => m === 1).length },
    }
  }),

  4: () => runStep(4, async () => {
    need('ctx')
    state.inputs = state.ctx.inputs
    return { input: { ctx: { inputs: 'see output', storeMask: `${state.ctx.storeMask.length} values` } }, output: state.inputs }
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
    const { model, productIds } = config()
    const params = { productIds }
    const name = `model-${model}-output`
    const param = { params, ctx: state.ctx, outputs: state.outputs }
    const { ret, sql } = await namedJsonQuery(name, param)
    state.result = ret
    return { sql, input: { namedQuery: name, param }, output: state.result }
  }),

  7: () => runStep(7, async () => {
    const { model, productIds, assets } = config()
    const opts = { model, params: { productIds }, assets }
    const result = await new I6Sdk().infer().run(opts)
    return { input: opts, output: result }
  }),
}

$('usage').textContent = usage()

for (const button of document.querySelectorAll('button[data-step]')) {
  button.onclick = steps[button.dataset.step]
}

// --- tools -------------------------------------------------------------------------------------------
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
