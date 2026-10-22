// Serve: python static/infer-demo/serve.py -> http://localhost:8000/static/infer-demo/index.html
import { I6Sdk } from '../../src/index.js'

// infer.js prints what it does when this is set (the log panel below shows it)
window.I6_DEBUG = true

const $ = (id) => document.getElementById(id)

// --- log ---------------------------------------------------------------------------------------------
const consoleLog = console.log
const log = (...args) => {
  $('log').textContent += args.map((a) => (typeof a === 'string' ? a : JSON.stringify(a))).join(' ') + '\n'
}
console.log = (...args) => { consoleLog(...args); log(...args) }
$('clearlog').onclick = () => { $('log').textContent = '' }

// --- session -----------------------------------------------------------------------------------------
// What the SDK user passes to `new I6Sdk(...)`. The session lives as long as the tab; the user is kept across visits.
const stored = (storage, key) => {
  try {
    return storage.getItem(key) ?? (storage.setItem(key, crypto.randomUUID()), storage.getItem(key))
  } catch { return crypto.randomUUID() }
}
const sdkConfig = { baseUrl: "https://decsuite.sandbox.rs.infinity6.ai", sessionId: stored(sessionStorage, 'i6demo.sessionId'), userId: stored(localStorage, 'i6demo.userId') }
// The config the SDK really runs with: the defaults (baseUrl) merged with ours
// (the demo builds the SDK without ids, then gives them with setConfig, which can be called anytime)
const sdk = new I6Sdk()
sdk.setConfig(sdkConfig)
window.i6sdk = sdk // try sdk.setConfig({ userId: 'x' }) in the console, then run
const showConfig = () => { $('sdkConfig').textContent = JSON.stringify(sdk.config(), null, 2) }
showConfig()
log('sdk config', sdk.config())

// What the SDK user passes as selectResults: picks the skus that are shown. Shown in the usage card as is.
// Here: best rated results first, no sku that is already in the cart, no sku twice, at most 4, with their slot in the carousel.
// The SDK adds the rating of each sku. Without selectResults, the SDK shows one random sku.
async function selectResults(req, { results }) {
  const inCart = new Set(req.input.skus.map((sku) => sku.sku_id))
  const shown = new Set()
  for (const result of [...results].sort((a, b) => b.rating - a.rating)) {
    for (const { sku_id } of result.skus) {
      if (!inCart.has(sku_id)) shown.add(sku_id)
    }
  }
  return {
    results: [...shown].slice(0, 4).map((sku_id, slot) => ({ sku_id, enriched: { slot } })),
  }
}

// --- what the SDK user passes in ---------------------------------------------------------------------
// The request of sdk.funcs().productRelevanceFashionFbtV1, from the config card.
function funcsRequest() {
  const text = (id) => $(id).value.trim()
  return {
    input: {
      dataset: text('dataset'),
      model_name: text('modelName'),
      skus: text('sku_ids').split(/[\s,]+/).filter(Boolean).map((sku_id) => ({ sku_id })),
    },
    partitions: { store_id: text('storeId'), version: text('modelVersion') },
  }
}

// The code of the SDK user, written with the values of the config card. Shown at the top of the page and run by the button
function usage() {
  return `import { I6Sdk } from './i6sdk-web.js'

const sdk = new I6Sdk()
sdk.setConfig(${JSON.stringify(sdkConfig)}) // anytime; overwrites only the given keys

${selectResults}

const result = await sdk.funcs().productRelevanceFashionFbtV1(${JSON.stringify(funcsRequest(), null, 2)}, selectResults)`
}
for (const id of ['modelName', 'modelVersion', 'storeId', 'dataset', 'sku_ids']) $(id).addEventListener('input', () => { $('usage').textContent = usage() })
$('usage').textContent = usage()

// --- run ---------------------------------------------------------------------------------------------
const setStatus = (status) => {
  const badge = $('step').querySelector('.badge')
  badge.textContent = status
  badge.className = `badge ${status}`
}

$('run').onclick = async () => {
  const card = $('step')
  const req = funcsRequest()
  setStatus('running')
  card.querySelector('.in').textContent = JSON.stringify(req, null, 1)
  card.querySelector('.out').textContent = ''
  try {
    const result = await sdk.funcs().productRelevanceFashionFbtV1(req, selectResults)
    // const result = await sdk.funcs().productRelevanceFashionFbtV1(req)
    card.querySelector('.out').textContent = JSON.stringify(result, null, 1)
    setStatus('ok')
  } catch (e) {
    card.querySelector('.out').textContent = `${e.name} [${e.code}] ${e.message}`
    setStatus('error')
  }
}

// --- dash --------------------------------------------------------------------------------------------
function dashRequest() {
  const text = (id) => $(id).value.trim()
  return {
    input: { dataset: text('dashDataset'), params: { min_year_month: text('minYearMonth'), max_year_month: text('maxYearMonth') } },
    partitions: { store_id: text('dashStoreId'), version: text('dashVersion') },
  }
}

const dashUsage = () => `const result = await new I6Sdk().funcs().${$('dashFunc').value}(${JSON.stringify(dashRequest(), null, 2)})`
const dashIds = ['dashFunc', 'dashDataset', 'dashStoreId', 'dashVersion', 'minYearMonth', 'maxYearMonth']
for (const id of dashIds) $(id).addEventListener('input', () => { $('dashUsage').textContent = dashUsage() })
$('dashUsage').textContent = dashUsage()

const setDashStatus = (status) => {
  $('dashBadge').textContent = status
  $('dashBadge').className = `badge ${status}`
}

$('dashRun').onclick = async () => {
  const req = dashRequest()
  setDashStatus('running')
  $('dashIn').textContent = JSON.stringify(req, null, 1)
  $('dashOut').textContent = ''
  try {
    const result = await sdk.funcs()[$('dashFunc').value](req)
    $('dashOut').textContent = JSON.stringify(result, null, 1)
    setDashStatus('ok')
  } catch (e) {
    $('dashOut').textContent = `${e.name} [${e.code}] ${e.message}`
    setDashStatus('error')
  }
}

// --- tools -------------------------------------------------------------------------------------------
// Redirects to the login page when nobody is signed in, and comes back here
$('login').onclick = async () => {
  try {
    const user = await sdk.auth().requireUser()
    if (user) log('logged in as', user)
  } catch (e) {
    log(`ERROR login: ${e.message}`)
  }
}

$('clear').onclick = async () => {
  await (await navigator.storage.getDirectory()).removeEntry('i6', { recursive: true }).catch(() => { })
  log('\nopfs cache removed')
}

window.onunhandledrejection = (e) => log(`ERROR ${e.reason.name} [${e.reason.code}] ${e.reason.message}`)
