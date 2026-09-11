#!/bin/bash -e

export GOWORK=off

eval "$(i6dev meta debug sdkgo-build I6DEV_DEBUG)"

function cmd_env() {
  cmd_run envconfig keys | i6dev box-config vault-export-keys
}

function cmd_update() {
  i6dev golang update
}

function cmd_test() {
  i6dev golang test_all "$@"
}

function cmd_release() {
  i6dev golang bin-compile
  i6dev golang release
  i6dev golang bin-release latest
}

function cmd_run() {
  go run . "$@"
}

function cmd_version-inc() {
  i6dev dever version-inc version/version.go
}

function cmd_version() {
  go run version/main/main.go
}

cd "$(dirname "$0")";
[ ! -f .venv/bin/activate ] || source .venv/bin/activate
_cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
