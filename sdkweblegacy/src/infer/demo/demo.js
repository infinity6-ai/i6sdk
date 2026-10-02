// Runs the stages of infer.js one by one and prints what happens. Served with:
//   cd sdkweblegacy && python -m http.server 8000   ->   http://localhost:8000/src/infer/demo/index.html
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

document.getElementById('run').onclick = run
document.getElementById('clear').onclick = async () => {
  await (await navigator.storage.getDirectory()).removeEntry('i6', { recursive: true }).catch(() => {})
  log('opfs cache removed')
}
