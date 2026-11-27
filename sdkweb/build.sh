#!/bin/bash -e

function cmd_npm_install() {
    npm install
}

function cmd_bundle() {
    [ ! -d dist ] || rm -rf dist
    go generate ./...
    cmd_npm_install

    # 1. Generate the standard .d.ts files into dist/types
    npx tsc

    # 2. Stitch them all into one giant file for IDEs & LLMs
    mkdir -p dist/bundle
    find dist/types -name "*.d.ts" -exec cat {} + > dist/bundle/i6sdkweb.d.ts
    cp dist/bundle/i6sdkweb.d.ts dist/bundle/i6sdk-web.d.ts

    # 3. Build the JavaScript bundles and modules
    npx esbuild src/main/index.js --bundle --sourcemap --outfile=dist/bundle/i6sdkweb.js
    npx esbuild src/main/index.js --bundle --minify --sourcemap --outfile=dist/bundle/i6sdkweb.min.js
    cp dist/bundle/i6sdkweb.js dist/bundle/i6sdk-web.js
    cp dist/bundle/i6sdkweb.js.map dist/bundle/i6sdk-web.js.map
    cp dist/bundle/i6sdkweb.min.js dist/bundle/i6sdk-web.min.js
    cp dist/bundle/i6sdkweb.min.js.map dist/bundle/i6sdk-web.min.js.map
    npx esbuild "src/main/**/*.js" --minify --outdir=dist/module-min
    npx esbuild "src/main/**/*.js" --outdir=dist/module
}

function cmd_release_prepare() {
    true
}

function cmd_release() {
    cmd_bundle
}

function cmd_test_visual() {
    BRWJSZ_DISABLE_HEADLESS=true go test -run '^TestWeb.*$' ./... "$@"
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
