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
    cmd_run server start "$@"
}

function cmd_test() {
    cmd_test_unit
    cmd_test_web
}

function cmd_test_unit() {
    GOWORK=off go test -run '^TestUnit.*$' ./... "$@"
}

function cmd_test_web() {
    GOWORK=off go test -run '^TestWeb.*$' ./... "$@"
}

function cmd_test_visual() {
    BRWJSZ_DISABLE_HEADLESS=true cmd_test_web "$@"
}

function cmd_npm_install() {
    npm install
}

# Generates src/apigenz (gitignored): one JS client per API, imported by src/ingest and friends.
# Must run before anything that resolves src/**, since a fresh checkout has no apigenz.
function cmd_generate() {
    go run . generate
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
    npx esbuild src/main/index.js --bundle --sourcemap --outfile=dist/bundle/i6sdk-web.js
    npx esbuild src/main/index.js --bundle --minify --sourcemap --outfile=dist/bundle/i6sdk-web.min.js
    npx esbuild "src/main/**/*.js" --minify --outdir=dist/module-min
    npx esbuild "src/main/**/*.js" --outdir=dist/module
}

function cmd_release() {
    cmd_bundle
    mkdir -p ../../i6pages/releases/sdkweb
    rsync -av --delete dist/ ../../i6pages/releases/sdkweb/latest/
    rsync -av --delete src/site/tokens/ ../../i6pages/tokens/ --exclude 'target'
    # rsync -av --delete ../sdkgo/ ../../i6pages/sdkgo/ --exclude 'target'
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

