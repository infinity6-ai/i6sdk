// Demo of the infer SDK: the config of an SDK user, one button per step (inputs and outputs shown) and
// a button that runs everything through `new I6Sdk().infer().run(opts)`.
// Steps 1-6 use the same internal pieces as infer.js (downloader, sqlitez, LiteRT) so each stage can be seen alone.
// Serve: python src/infer/demo/serve.py -> http://localhost:8000/src/infer/demo/index.html
import { I6Sdk } from '../../index.js'
import { downloader } from '../../internal/downloader/downloader.js'
import { fs } from '../../internal/fs/fs.js'
import { sqlitez } from '../../internal/sqlitez/sqlitez.js'

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
// Long arrays are cut (first 10 + how many more) unless "show full arrays" is checked
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
    params: { productIds: [...new Set(productIds)] },
    assets: { sqlite: $('sqlite').value.trim(), tflite: $('tflite').value.trim() },
  }
}

// The caller's job (not infer.js): turn the result of the input query into the FBT model inputs.
// Last items' rec_in vectors, one masked row (2), then padding rows (0); pad_mask is 1 up to the masked row.
const NUM_ITEMS = 10
const VECTOR_LENGTH = 7
const buildFbtInputs = (productIds) => ({ products }) => {
  const byId = Object.fromEntries(products.map((p) => [p.product_id, p]))
  const vectors = productIds
    .map((id) => byId[id]?.rec_in?.split('|').map(Number))
    .filter((v) => v && v.length === VECTOR_LENGTH && v.every(Number.isFinite))
    .slice(-(NUM_ITEMS - 1))
  const rows = [...vectors, Array(VECTOR_LENGTH).fill(2)]
  while (rows.length < NUM_ITEMS) rows.push(Array(VECTOR_LENGTH).fill(0))
  return {
    fields: { data: rows.flat(), shape: [1, NUM_ITEMS, VECTOR_LENGTH] },
    pad_mask: { data: rows.map((_, i) => (i <= vectors.length ? 1 : 0)), shape: [1, NUM_ITEMS] },
  }
}

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

const describe = (d) => ({ name: d.name, dtype: d.dtype, shape: Array.from(d.shape) })
const DTYPE_CTORS = { float32: Float32Array, int32: Int32Array, uint8: Uint8Array }

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
    const litert = await import('@litertjs/core')
    const loading = litert.getGlobalLiteRtPromise()
    if (loading) await loading
    else await litert.loadLiteRt('https://cdn.jsdelivr.net/npm/@litertjs/core@2.5.3/wasm/')
    if (state.compiled) state.compiled.delete()
    state.litert = litert
    state.compiled = await litert.loadAndCompile(state.modelBytes, { accelerator: 'wasm' })
    return {
      input: { tflite: `${state.modelBytes.length} bytes`, accelerator: 'wasm' },
      output: { inputs: state.compiled.getInputDetails().map(describe), outputs: state.compiled.getOutputDetails().map(describe) },
    }
  }),

  3: () => runStep(3, async () => {
    need('dbPath')
    const { model, params } = config()
    const name = `model-${model}-input`
    const { ret, sql } = await namedJsonQuery(name, params)
    state.ctx = ret
    return {
      sql,
      input: { namedQuery: name, param: params },
      output: { ...state.ctx, 'storeMask: ones': state.ctx.storeMask.filter((m) => m === 1).length, 'products: found': state.ctx.products.length },
    }
  }),

  4: () => runStep(4, async () => {
    need('ctx')
    state.inputs = buildFbtInputs(config().productIds)(state.ctx)
    return { input: { products: state.ctx.products, storeMask: `${state.ctx.storeMask.length} values` }, output: state.inputs }
  }),

  5: () => runStep(5, async () => {
    need('compiled', 'inputs')
    const tensors = {}
    const described = {}
    let raw = {}
    try {
      for (const d of state.compiled.getInputDetails()) {
        const input = state.inputs[d.name]
        if (!input) throw new Error(`step 4 did not return "${d.name}"`)
        const shape = input.shape ?? Array.from(d.shape, (x) => (x === -1 ? 1 : x))
        tensors[d.name] = new state.litert.Tensor(new DTYPE_CTORS[d.dtype](input.data), shape)
        described[d.name] = { dtype: d.dtype, shape, data: input.data }
      }
      raw = await state.compiled.run(tensors)
      state.outputs = {}
      const output = {}
      for (const [name, tensor] of Object.entries(raw)) {
        const values = Array.from(await tensor.data())
        state.outputs[name] = values
        output[name] = { length: values.length, values }
      }
      // The best 5 raw scores, so they can be checked at a glance
      const scores = state.outputs.signature_probs
      if (scores) output.signature_probs.best5 = scores.map((score, i) => ({ i, score })).sort((a, b) => b.score - a.score).slice(0, 5)
      return { input: described, output }
    } finally {
      Object.values(raw).forEach((t) => t.delete())
      Object.values(tensors).forEach((t) => t.delete())
    }
  }),

  6: () => runStep(6, async () => {
    need('ctx', 'outputs')
    const { model, params } = config()
    const name = `model-${model}-output`
    const param = { params, ctx: state.ctx, outputs: state.outputs }
    const { ret, sql } = await namedJsonQuery(name, param)
    state.result = ret
    return { sql, input: { namedQuery: name, param }, output: state.result }
  }),

  7: () => runStep(7, async () => {
    const { model, productIds, params, assets } = config()
    const opts = { model, params, assets, buildInputs: buildFbtInputs(productIds) }
    const result = await new I6Sdk().infer().run(opts)
    const same = state.result ? JSON.stringify(result) === JSON.stringify(state.result) : 'n/a (run step 6 to compare)'
    return { input: { ...opts, buildInputs: '[function]' }, output: { 'same as step 6': same, result } }
  }),
}

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
