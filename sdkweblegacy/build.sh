#!/bin/bash -e

export GOWORK=off

eval "$(i6dev meta debug sdkweb I6DEV_DEBUG)"

function cmd_static() {
    i6dev go-platform static . static internal/stfiles
}

function cmd_update() {
    i6dev golang update
}

function cmd_run_server() {
    # [ "x$(i6dev box-config jq -r .id)" == "xlocal" ]
    # eval "$(cmd_env)"
    cmd_generate
    I6_WEB_CHECKER_ENABLED=true cmd_run server start "$@"
}

function cmd_run_testserver() {
    cmd_generate
    cmd_run appstack testserver
}

function cmd_test() {
    cmd_test_unit
    cmd_test_web
}

function cmd_test_unit() {
    i6dev golang test_unit "$@"
}

function cmd_test_web() {
    i6dev golang test_web "$@"
}

function cmd_test_visual() {
    I6_WEB_HEADLESS_DISABLED=true cmd_test_web "$@"
}

function cmd_npm_install() {
    npm install
}

# Generates src/apigenz (gitignored): one JS client per API, imported by src/ingest and friends.
# Must run before anything that resolves src/**, since a fresh checkout has no apigenz.
function cmd_generate() {
    (cd .. && go run . js generate)
}

function cmd_bundle() {
    [ ! -d dist ] || rm -rf dist
    cmd_generate
    cmd_npm_install
    
    # 1. Generate the standard .d.ts files into dist/types
    npx tsc
    
    # 2. Stitch them all into one giant file for the LLM 
    # (mkdir ensures the bundle folder exists before we create the file)
    mkdir -p dist/bundle
    find dist/types -name "*.d.ts" -exec cat {} + > dist/bundle/i6sdk-web.d.ts

    # 3. Build the JavaScript
    npx esbuild src/index.js --bundle --sourcemap --outfile=dist/bundle/i6sdk-web.js
    npx esbuild src/index.js --bundle --minify --sourcemap --outfile=dist/bundle/i6sdk-web.min.js
    npx esbuild "src/**/*.js" --minify --outdir=dist/module-min
    npx esbuild "src/**/*.js" --outdir=dist/module
}

function cmd_release() {
    cmd_bundle
    mkdir -p ../../i6pages/releases/sdkweb
    rsync -av --delete dist/ ../../i6pages/releases/sdkweb/latest/
    rsync -av --delete static/tokens/ ../../i6pages/tokens/ --exclude 'target'
    # rsync -av --delete ../sdkgo/ ../../i6pages/sdkgo/ --exclude 'target'
    cmd_release-infer-demo
}

# Publishes static/infer-demo to i6pages/infer-demo, running against the released bundle instead of the source
function cmd_release-infer-demo() {
    local src=static/infer-demo dst=../../i6pages/infer-demo
    local tag='<script src="../releases/sdkweb/latest/bundle/i6sdk-web.js"></script>'
    mkdir -p "$dst"
    rm -f "$dst"/i6sdk-web.js "$dst"/i6sdk-web.js.map

    sed -e "1,2c\\
// Serve this repo's root: python -m http.server 8000 -> http://localhost:8000/infer-demo/\\
// The SDK comes from the published bundle (index.html), not from source: this page tests what the CDN serves.\\
const { I6Sdk } = window" \
        -e "s#return \`import { I6Sdk } from './i6sdk-web.js'#return \`// <script src=\"https://cdn.infinity6.ai/releases/sdkweb/latest/bundle/i6sdk-web.js\"></script>  (sets window.I6Sdk)#" \
        "$src/demo.js" > "$dst/demo.js"

    sed -e '/<!-- Demo only/,/<\/script>/c\
  <!-- The published SDK bundle of this repo (sql.js is inside it); sets window.I6Sdk -->\
  '"$tag" \
        "$src/index.html" > "$dst/index.html"

    sed -e "s#import { I6Sdk } from '[^']*'#const { I6Sdk } = window#" \
        -e "s#^\( *\)<script type=\"module\">#\1$tag\n\1<script type=\"module\">#" \
        "$src/logincallback.html" > "$dst/logincallback.html"

    # Fail loudly if a rewrite stopped matching
    grep -q '^const { I6Sdk } = window' "$dst/demo.js" && ! grep -q "src/index.js" "$dst/demo.js" \
        && ! grep -q importmap "$dst/index.html" && grep -q "i6sdk-web.js" "$dst/index.html" \
        && grep -q "i6sdk-web.js" "$dst/logincallback.html" && ! grep -q "src/index.js" "$dst/logincallback.html" \
        || { echo "infer-demo release: a rewrite did not match, check $src" >&2; return 1; }
}

function cmd_vendor() {
    i6dev golang vendor
}


function cmd_run() {
    go run main.go "$@"
}

function cmd_update_i6() {
  i6dev golang gomod-update-i6
}

function cmd_release-describe() {
  i6dev golang release-describe
}

cd "$(dirname "$0")"; 
[ ! -f .venv/bin/activate ] || source .venv/bin/activate
_cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"

