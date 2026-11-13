import { I6Sdk } from '../../src/index.js'

window.I6_DEBUG = true

const $ = (id) => document.getElementById(id)

// --- log ---------------------------------------------------------------------------------------------
const consoleLog = console.log
const log = (...args) => {
  const logEl = $('log')
  logEl.textContent += args.map((a) => (typeof a === 'string' ? a : JSON.stringify(a, null, 2))).join(' ') + '\n\n'
  logEl.scrollTop = logEl.scrollHeight
}
console.log = (...args) => { consoleLog(...args); log(...args) }
$('clearlog').onclick = () => { $('log').textContent = '' }

// --- sdk setup ---------------------------------------------------------------------------------------
const sdk = new I6Sdk()
window.i6sdk = sdk

const showConfig = () => { $('sdkConfig').textContent = JSON.stringify(sdk.config(), null, 2) }
showConfig()
log('sdk config', sdk.config())

// --- ui state ----------------------------------------------------------------------------------------
const setStatus = (status) => {
  const badge = $('statusBadge')
  badge.textContent = status
  badge.className = 'badge ' + (status === 'running' ? 'running' : status.startsWith('ok') ? 'ok' : status.startsWith('error') ? 'error' : '')
}

// --- request building --------------------------------------------------------------------------------
function buildRequest() {
  const text = (id) => $(id).value.trim()
  return {
    headers: {
      x_i6_signature: text('signature')
    },
    payload: {
      signature: text('signature'),
      method: text('method'),
      path: text('path'),
      query: text('query'),
      body: text('body')
    }
  }
}

// --- run ---------------------------------------------------------------------------------------------
$('run').onclick = async () => {
  $('reqOut').textContent = ''
  $('respOut').textContent = ''
  setStatus('running')
  
  try {
    const req = buildRequest()
    $('reqOut').textContent = JSON.stringify(req, null, 2)
    
    // Call the generated API method
    const resp = await sdk.apis().userSignCheck(req)
    
    $('respOut').textContent = JSON.stringify(resp, null, 2)
    
    if (resp.ok && resp.payload?.valid) {
      setStatus('ok (valid)')
    } else if (resp.ok) {
      setStatus('error (invalid)')
    } else {
      setStatus(`error (${resp.status})`)
    }
  } catch (err) {
    log('Caught error:', err.stack || err.message || err)
    $('respOut').textContent = String(err.stack || err.message || err)
    setStatus('error (exception)')
  }
}
