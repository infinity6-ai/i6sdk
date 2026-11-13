// Serve: python -m http.server 8000 -> http://localhost:8000/sdkweblegacy/static/reqsign/reqsign.html
import { I6Sdk } from '../../src/index.js'

window.I6_DEBUG = true

// --- Dynamic ESM Imports with Web Crypto Fallback ----------------------------------------------------
let nobleEd25519 = null
try {
  const mod = await import('@noble/curves/ed25519')
  nobleEd25519 = mod.ed25519
} catch (e) {
  console.warn('Could not load @noble/curves/ed25519 from CDN, relying on Web Crypto API:', e)
}

const $ = (id) => document.getElementById(id)

// --- Shared Application State ------------------------------------------------------------------------
const DEMO_KEY_SEED_B64 = 'UQHxCZMYdkNbbtdGt9cnfwAdOnQAkT8Z12tDXwEAPFk'
let currentSeed = null
let currentPubKey = null
let lastComputed = null

// --- Console Logger ----------------------------------------------------------------------------------
const consoleLog = console.log
const log = (...args) => {
  const logEl = $('log')
  if (!logEl) return
  logEl.textContent += args.map((a) => (typeof a === 'string' ? a : JSON.stringify(a, null, 2))).join(' ') + '\n\n'
  logEl.scrollTop = logEl.scrollHeight
}
console.log = (...args) => { consoleLog(...args); log(...args) }
$('clearlog').onclick = () => { $('log').textContent = '' }

// --- SDK Initialization ------------------------------------------------------------------------------
const sdkConfig = {
  baseUrl: "https://decsuite.sandbox.rs.infinity6.ai",
  servUrls: {
    "userz": "https://userz.sandbox.rs.infinity6.ai",
    "tfinfer": "https://tfinfer.sandbox.rs.infinity6.ai"
  }
}
const sdk = new I6Sdk(sdkConfig)
window.i6sdk = sdk

function updateSdkConfigDisplay() {
  $('sdkConfigPre').textContent = JSON.stringify(sdk.config(), null, 2)
}
updateSdkConfigDisplay()
$('btnReloadConfig').onclick = updateSdkConfigDisplay

// --- Binary & Encoding Helpers -----------------------------------------------------------------------
function b64urlEncode(buf) {
  let bin = ''
  const bytes = buf instanceof Uint8Array ? buf : new Uint8Array(buf)
  for (let i = 0; i < bytes.length; i++) {
    bin += String.fromCharCode(bytes[i])
  }
  return btoa(bin).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
}

function b64urlDecode(str) {
  str = str.trim().replace(/-/g, '+').replace(/_/g, '/')
  while (str.length % 4) str += '='
  const bin = atob(str)
  const arr = new Uint8Array(bin.length)
  for (let i = 0; i < bin.length; i++) {
    arr[i] = bin.charCodeAt(i)
  }
  return arr
}

function hexToBytes(hex) {
  hex = hex.trim()
  if (hex.length % 2 !== 0) throw new Error('Hex string must have an even number of characters')
  const bytes = new Uint8Array(hex.length / 2)
  for (let i = 0; i < bytes.length; i++) {
    bytes[i] = parseInt(hex.substr(i * 2, 2), 16)
  }
  return bytes
}

function bytesToHex(bytes) {
  const arr = bytes instanceof Uint8Array ? bytes : new Uint8Array(bytes)
  return Array.from(arr).map((b) => b.toString(16).padStart(2, '0')).join('')
}

async function sha256Hex(data) {
  if (data === '' || data === null || data === undefined) return ''
  const buf = typeof data === 'string' ? new TextEncoder().encode(data) : data
  const hash = await crypto.subtle.digest('SHA-256', buf)
  return bytesToHex(new Uint8Array(hash))
}

function queryEscape(s) {
  return encodeURIComponent(s)
    .replace(/%20/g, '+')
    .replace(/[!'()*]/g, (c) => '%' + c.charCodeAt(0).toString(16).toUpperCase())
}

function canonicalizePath(p) {
  p = (p || '').trim()
  if (!p) return '/'
  if (!p.startsWith('/')) p = '/' + p
  if (p.length > 1 && p.endsWith('/')) p = p.slice(0, -1)
  return p
}

function canonicalizeQuery(q) {
  q = (q || '').trim().replace(/^\?/, '')
  if (!q) return ''

  const params = new URLSearchParams(q)
  params.delete('i6_sig')

  const map = new Map()
  for (const [k, v] of params.entries()) {
    if (!map.has(k)) map.set(k, [])
    map.get(k).push(v)
  }

  const sortedKeys = Array.from(map.keys()).sort()
  const pairs = []
  for (const k of sortedKeys) {
    const vals = map.get(k).sort()
    const encK = queryEscape(k)
    for (const v of vals) {
      pairs.push(`${encK}=${queryEscape(v)}`)
    }
  }
  return pairs.join('&')
}

// --- Key Management & Derivation ---------------------------------------------------------------------
function parsePrivateKeyInput(raw) {
  raw = raw.trim()
  if (!raw) return null

  // 1. Try Hex
  if (/^[0-9a-fA-F]+$/.test(raw)) {
    const b = hexToBytes(raw)
    if (b.length === 32 || b.length === 64) {
      return {
        seed: b.subarray(0, 32),
        pubKey: b.length === 64 ? b.subarray(32, 64) : null,
        format: `Hex (${b.length}B)`
      }
    }
  }

  // 2. Try Base64 / Base64URL
  try {
    const b = b64urlDecode(raw)
    if (b.length === 32 || b.length === 64) {
      return {
        seed: b.subarray(0, 32),
        pubKey: b.length === 64 ? b.subarray(32, 64) : null,
        format: `Base64URL (${b.length}B)`
      }
    }
  } catch (_) {}

  throw new Error('Unrecognized private key format: Expected 32-byte seed or 64-byte key in Base64URL, Base64, or Hex.')
}

async function derivePublicKey(seed, existingPubKey) {
  if (existingPubKey && existingPubKey.length === 32) {
    return existingPubKey
  }
  if (nobleEd25519) {
    return nobleEd25519.getPublicKey(seed)
  }
  return null
}

async function updateKeypairState() {
  const rawKey = $('privKey').value.trim()
  const badge = $('keyFormatBadge')
  const pubDisplay = $('pubKeyDisplay')

  if (!rawKey) {
    currentSeed = null
    currentPubKey = null
    badge.textContent = 'No key loaded'
    badge.className = 'badge'
    pubDisplay.style.display = 'none'
    return
  }

  try {
    const parsed = parsePrivateKeyInput(rawKey)
    currentSeed = parsed.seed
    badge.textContent = `Valid ${parsed.format}`
    badge.className = 'badge ok'

    currentPubKey = await derivePublicKey(currentSeed, parsed.pubKey)
    if (currentPubKey) {
      const pubB64 = b64urlEncode(currentPubKey)
      $('pubKeyB64').textContent = pubB64
      $('pubKeyHex').textContent = bytesToHex(currentPubKey)
      pubDisplay.style.display = 'block'
      updateCliCommand(pubB64)
    } else {
      $('pubKeyB64').textContent = '(Derivation requires @noble/curves or 64B key)'
      $('pubKeyHex').textContent = ''
      pubDisplay.style.display = 'block'
      updateCliCommand('<CLIENT_PUB_KEY>')
    }
  } catch (err) {
    currentSeed = null
    currentPubKey = null
    badge.textContent = 'Invalid key format'
    badge.className = 'badge error'
    pubDisplay.style.display = 'none'
    updateCliCommand('<CLIENT_PUB_KEY>')
  }
}

function updateCliCommand(pubKeyB64) {
  const el = $('cliCreateAccessCmd')
  if (!el) return
  el.textContent = `go run ../i6userz token dataset create-access \\\n  --endpoint "https://userz.sandbox.rs.infinity6.ai" \\\n  --dataset "ds64l7qauv4l6512kf0c" \\\n  --client-pub-key "${pubKeyB64 || '<CLIENT_PUB_KEY>'}" \\\n  --grants "read,write"`
}

async function signCanonicalBytes(seed, canonicalBytes) {
  if (nobleEd25519) {
    return nobleEd25519.sign(canonicalBytes, seed)
  }

  // Fallback: Web Crypto API PKCS#8 import
  const pkcs8Prefix = new Uint8Array([0x30, 0x2e, 0x02, 0x01, 0x00, 0x30, 0x05, 0x06, 0x03, 0x2b, 0x65, 0x70, 0x04, 0x22, 0x04, 0x20])
  const pkcs8 = new Uint8Array(48)
  pkcs8.set(pkcs8Prefix, 0)
  pkcs8.set(seed, 16)

  const cryptoKey = await crypto.subtle.importKey('pkcs8', pkcs8, { name: 'Ed25519' }, false, ['sign'])
  const sigBuf = await crypto.subtle.sign({ name: 'Ed25519' }, cryptoKey, canonicalBytes)
  return new Uint8Array(sigBuf)
}

// --- Signature Computation ---------------------------------------------------------------------------
async function computeSignature() {
  const badge = $('signStatusBadge')
  const headerEl = $('headerOutput')
  const canonicalEl = $('canonicalOutput')
  const envEl = $('envelopeBreakdown')

  if (!currentSeed) {
    badge.textContent = 'Awaiting Private Key'
    badge.className = 'badge'
    headerEl.textContent = '(Enter a valid private key above to generate the signature)'
    canonicalEl.textContent = ''
    envEl.textContent = ''
    lastComputed = null
    return
  }

  try {
    const method = $('method').value.trim().toUpperCase()
    const path = canonicalizePath($('path').value)
    const query = canonicalizeQuery($('query').value)
    const bodyStr = $('body').value
    const accessToken = $('accessToken').value.trim()
    const unsignedBody = $('unsignedBody').checked

    let exp = parseInt($('expEpoch').value.trim(), 10)
    if (isNaN(exp) || exp <= 0) {
      exp = Math.floor(Date.now() / 1000) + 300
    }

    let nonceHex = $('nonceHex').value.trim().toLowerCase()
    if (nonceHex.length !== 32 || !/^[0-9a-f]+$/.test(nonceHex)) {
      nonceHex = '00000000000000000000000000000000'
    }

    const bodyHash = unsignedBody ? '' : await sha256Hex(bodyStr)
    const tokenHash = accessToken ? await sha256Hex(accessToken) : ''

    const canonicalLines = [
      method,
      path,
      query,
      bodyHash,
      tokenHash,
      String(exp),
      nonceHex
    ]
    const canonicalString = canonicalLines.join('\n')

    canonicalEl.textContent = canonicalLines
      .map((line, idx) => `[${idx + 1}] ${['METHOD', 'PATH', 'CANONICAL QUERY', 'BODY SHA-256', 'TOKEN SHA-256', 'EXPIRATION SECONDS', 'NONCE HEX'][idx]}: ${JSON.stringify(line)}`)
      .join('\n')

    const canonicalBytes = new TextEncoder().encode(canonicalString)
    const signatureBytes = await signCanonicalBytes(currentSeed, canonicalBytes)

    // Assemble 89-byte binary envelope
    const envelopeBuf = new Uint8Array(89)
    envelopeBuf[0] = 2 | (unsignedBody ? 0x80 : 0) // Type 2 (SigEnvelope) + Flags
    const dv = new DataView(envelopeBuf.buffer)
    dv.setBigUint64(1, BigInt(exp), false) // 8B BigEndian

    const nonceBytes = hexToBytes(nonceHex)
    envelopeBuf.set(nonceBytes, 9) // 16B nonce
    envelopeBuf.set(signatureBytes, 25) // 64B Ed25519 signature

    const envelopeB64 = b64urlEncode(envelopeBuf)
    const fullHeader = accessToken ? `${envelopeB64}.${accessToken}` : envelopeB64

    headerEl.textContent = fullHeader

    envEl.textContent = [
      `Length        : 89 bytes (Encoded: ${envelopeB64.length} chars Base64URL)`,
      `Byte 0 (Type) : ${envelopeBuf[0]} (Type=2, UnsignedBody=${unsignedBody})`,
      `Bytes 1..8    : Exp=${exp} (${new Date(exp * 1000).toISOString()})`,
      `Bytes 9..24   : Nonce=${nonceHex}`,
      `Bytes 25..88  : Signature (64B) = ${bytesToHex(signatureBytes).slice(0, 32)}...`,
      accessToken ? `Attached Token: ${accessToken.slice(0, 24)}...` : `Attached Token: (none)`
    ].join('\n')

    badge.textContent = 'Signature Ready'
    badge.className = 'badge ok'

    lastComputed = {
      headerValue: fullHeader,
      method,
      path,
      query,
      bodyStr,
      accessToken,
      canonicalString
    }
  } catch (err) {
    badge.textContent = 'Signing error'
    badge.className = 'badge error'
    headerEl.textContent = 'Error: ' + err.message
    log('Signing error:', err)
  }
}

// --- TTL & Nonce Helpers -----------------------------------------------------------------------------
function setTtlSeconds(seconds) {
  const exp = Math.floor(Date.now() / 1000) + seconds
  $('expEpoch').value = String(exp)
  computeSignature()
}

function updateNonceDisplay() {
  if ($('nonceZero').checked) {
    $('nonceHex').value = '00000000000000000000000000000000'
    $('nonceHex').readOnly = true
  } else {
    const rnd = new Uint8Array(16)
    crypto.getRandomValues(rnd)
    $('nonceHex').value = bytesToHex(rnd)
    $('nonceHex').readOnly = false
  }
  computeSignature()
}

// --- Event Listeners ---------------------------------------------------------------------------------
$('privKey').addEventListener('input', async () => {
  await updateKeypairState()
  computeSignature()
})

$('togglePrivKey').onclick = () => {
  const input = $('privKey')
  const btn = $('togglePrivKey')
  if (input.type === 'password') {
    input.type = 'text'
    btn.textContent = 'Hide'
  } else {
    input.type = 'password'
    btn.textContent = 'Show'
  }
}

$('genKeypair').onclick = async () => {
  try {
    if (nobleEd25519) {
      const priv = nobleEd25519.utils.randomPrivateKey()
      $('privKey').value = b64urlEncode(priv)
      $('privKey').type = 'text'
      $('togglePrivKey').textContent = 'Hide'
    } else {
      const kp = await crypto.subtle.generateKey({ name: 'Ed25519' }, true, ['sign', 'verify'])
      const privPkcs8 = await crypto.subtle.exportKey('pkcs8', kp.privateKey)
      const seed = new Uint8Array(privPkcs8).slice(16, 48)
      $('privKey').value = b64urlEncode(seed)
      $('privKey').type = 'text'
      $('togglePrivKey').textContent = 'Hide'
    }
    await updateKeypairState()
    computeSignature()
    log('Generated fresh Ed25519 keypair')
  } catch (e) {
    alert('Key generation failed: ' + e.message)
  }
}

$('loadDemoKey').onclick = async () => {
  $('privKey').value = DEMO_KEY_SEED_B64
  $('privKey').type = 'text'
  $('togglePrivKey').textContent = 'Hide'
  await updateKeypairState()
  computeSignature()
  log('Loaded sample demo key')
}

$('clearKey').onclick = async () => {
  $('privKey').value = ''
  await updateKeypairState()
  computeSignature()
}

$('copyPubB64').onclick = () => {
  navigator.clipboard.writeText($('pubKeyB64').textContent)
}
$('copyPubHex').onclick = () => {
  navigator.clipboard.writeText($('pubKeyHex').textContent)
}

$('useSdkToken').onclick = () => {
  const all = sdk.tokens().getAllDatasetTokens()
  const keys = Object.keys(all)
  if (keys.length > 0) {
    const first = all[keys[0]]
    $('accessToken').value = first.access_token || ''
    computeSignature()
    log('Loaded dataset token for dataset:', keys[0])
  } else {
    alert('No dataset tokens cached in SDK storage. Authenticate with sdk.tokens() first or paste a token manually.')
  }
}

$('formatJson').onclick = () => {
  try {
    const val = $('body').value.trim()
    if (!val) return
    $('body').value = JSON.stringify(JSON.parse(val), null, 2)
    computeSignature()
  } catch (e) {
    alert('Invalid JSON: ' + e.message)
  }
}

$('clearBody').onclick = () => {
  $('body').value = ''
  computeSignature()
}

document.querySelectorAll('.ttl-btn').forEach((btn) => {
  btn.onclick = () => setTtlSeconds(parseInt(btn.getAttribute('data-ttl'), 10))
})

$('nonceZero').onchange = updateNonceDisplay
$('nonceRandom').onchange = updateNonceDisplay

for (const id of ['method', 'path', 'query', 'body', 'accessToken', 'expEpoch', 'nonceHex', 'unsignedBody']) {
  $(id).addEventListener('input', () => computeSignature())
  $(id).addEventListener('change', () => computeSignature())
}

$('btnSign').onclick = () => computeSignature()

$('copyHeaderBtn').onclick = () => {
  const val = $('headerOutput').textContent
  if (val && !val.startsWith('(') && !val.startsWith('Error')) {
    navigator.clipboard.writeText(val)
    $('copyHeaderBtn').textContent = 'Copied!'
    setTimeout(() => { $('copyHeaderBtn').textContent = 'Copy X-I6-Signature' }, 1500)
  }
}

$('copyCurlBtn').onclick = () => {
  if (!lastComputed) return
  const url = `${sdk.config().servUrls?.['userz'] || sdk.config().baseUrl}${lastComputed.path}${lastComputed.query ? '?' + lastComputed.query : ''}`
  let cmd = `curl -X ${lastComputed.method} "${url}" \\\n  -H "X-I6-Signature: ${lastComputed.headerValue}"`
  if (lastComputed.bodyStr) {
    cmd += ` \\\n  -H "Content-Type: application/json" \\\n  -d ${JSON.stringify(lastComputed.bodyStr)}`
  }
  navigator.clipboard.writeText(cmd)
  $('copyCurlBtn').textContent = 'Copied curl!'
  setTimeout(() => { $('copyCurlBtn').textContent = 'Copy curl' }, 1500)
}

// --- signcheck API Diagnostics -----------------------------------------------------------------------
$('copyCliCmd').onclick = () => {
  const cmd = $('cliCreateAccessCmd').textContent
  if (cmd) {
    navigator.clipboard.writeText(cmd)
    $('copyCliCmd').textContent = 'Copied!'
    setTimeout(() => { $('copyCliCmd').textContent = 'Copy Command' }, 1500)
  }
}

$('btnCheckApi').onclick = async () => {
  const badge = $('checkBadge')
  const banner = $('statusBanner')
  const reqPre = $('apiReqPre')
  const respPre = $('apiRespPre')

  if (!lastComputed) {
    await computeSignature()
  }
  if (!lastComputed) {
    alert('Please enter a private key and ensure the signature is computed first.')
    return
  }

  badge.textContent = 'verifying...'
  badge.className = 'badge running'
  banner.style.display = 'none'

  const apiReq = {
    headers: {
      x_i6_signature: lastComputed.headerValue
    },
    payload: {
      signature: lastComputed.headerValue,
      method: lastComputed.method,
      path: lastComputed.path,
      query: lastComputed.query,
      body: lastComputed.bodyStr
    }
  }

  reqPre.textContent = JSON.stringify(apiReq, null, 2)
  respPre.textContent = 'Waiting for response from userapisigncheck...'

  try {
    log('Calling sdk.apis().userSignCheck(...)')
    const resp = await sdk.apis().userSignCheck(apiReq)
    log('userSignCheck response received: HTTP', resp.status, resp.text)

    // Safely extract and parse response payload
    let payload = null
    try {
      payload = resp.parse()
    } catch (_) {
      try {
        payload = JSON.parse(resp.text)
      } catch (e) {
        payload = { raw: resp.text }
      }
    }

    // Display formatted payload directly in the output pre
    respPre.textContent = JSON.stringify(payload ?? { raw: resp.text }, null, 2)

    banner.style.display = 'block'
    if (resp.ok && payload?.valid) {
      badge.textContent = 'VALID'
      badge.className = 'badge ok'
      banner.className = 'status-box valid'
      banner.innerHTML = `<strong>VALID:</strong> ${payload?.summary || 'Request signature successfully verified!'}`
    } else {
      badge.textContent = 'INVALID'
      badge.className = 'badge error'
      banner.className = 'status-box invalid'

      const errCode = payload?.error_code || (resp.ok ? 'invalid_signature' : `http_${resp.status}`)
      const summary = payload?.summary || resp.error?.message || 'Signature check failed'
      let details = `<strong>[${errCode}]</strong> ${summary}`

      if (payload?.diagnostics?.hint) {
        details += `<div style="margin-top:6px; font-size:12px; opacity:0.95"><strong>💡 Server Hint:</strong> ${payload.diagnostics.hint}</div>`
      }
      banner.innerHTML = details
    }
  } catch (err) {
    badge.textContent = 'request error'
    badge.className = 'badge error'
    respPre.textContent = String(err.stack || err.message || err)
    banner.style.display = 'block'
    banner.className = 'status-box invalid'
    banner.textContent = 'Exception invoking SDK API: ' + err.message
    log('API call exception:', err)
  }
}

// --- Initialization ----------------------------------------------------------------------------------
async function init() {
  if (!$('privKey').value) {
    $('privKey').value = DEMO_KEY_SEED_B64
    $('privKey').type = 'text'
    $('togglePrivKey').textContent = 'Hide'
  }
  setTtlSeconds(300)
  updateNonceDisplay()
  await updateKeypairState()
  await computeSignature()
}

init()
