#!/bin/bash -e

eval "$(i6dev meta debug i6sdk-web I6DEV_DEBUG)"

# function cmd_env() {
#     cmd_run envconfig keys | i6dev box-config vault-export-keys
# }

function cmd_static() {
    i6dev go-platform static . static internal/stfiles
}

function cmd_update() {
    i6dev golang update
}

function cmd_run_server() {
    # [ "x$(i6dev box-config jq -r .id)" == "xlocal" ]
    # eval "$(cmd_env)"
    I6_WEB_CHECKER_ENABLED=true cmd_run server start "$@"
}

function cmd_run_testserver() {
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
    # true
    i6dev golang test_web "$@"
}

function cmd_test_visual() {
    I6_WEB_HEADLESS_DISABLED=true cmd_test_web "$@"
}

function cmd_bundle() {
    [ ! -d dist ] || rm -rf dist
    /opt/node/bin/esbuild static/sdkweb/sdkweb.js --bundle --minify --sourcemap --outfile=dist/bundle/i6sdkwev.min.js
    /opt/node/bin/esbuild "static/sdkweb/**/*.js" --minify --outdir=dist/sdkweb
}

function cmd_release() {
    cmd_static
    i6dev golang bin-compile
    i6dev golang release
    i6dev golang bin-release latest
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

