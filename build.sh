#!/bin/bash -e

eval "$(i6dev meta debug sdkgo-build I6DEV_DEBUG)"

function cmd_env() {
  cmd_run envconfig keys | i6dev box-config vault-export-keys
}

function cmd_clean() {
  local _k=""
  find -name target | while read _k; do
    rm -rf "$_k"
  done
}

function cmd_update() {
  i6dev golang auth
  I6DEV_GOLANG_AUTH_DISABLED="true" i6dev golang update
  I6DEV_GOLANG_AUTH_DISABLED="true" ./sdkweb/build.sh update
}

function cmd_test() {
  i6dev golang test_all "$@"
  ./sdkweb/build.sh test "$@"
}

function cmd_test_remote() {
  ./sdkweb/build.sh test_remote "$@"
}

function cmd_test_all() {
  ./sdkweb/build.sh test_all "$@"
}

function cmd_vendor() {
  i6dev golang vendor
  ./sdkweb/build.sh vendor
}

function cmd_release() {
  i6dev golang bin-compile
  i6dev golang release
  i6dev golang bin-release latest
  ./sdkweb/build.sh release
}

function cmd_update_i6() {
  i6dev golang gomod-update-i6
  ./sdkweb/build.sh update_i6
}

function cmd_run() {
  go run . "$@"
}

function cmd_version-inc() {
  i6dev dever version-inc sdkgo/version/version.go
  i6dev dever version-inc sdkweb/version/version.go
}

function cmd_version() {
  go run version/main/main.go
}

cd "$(dirname "$0")";
[ ! -f .venv/bin/activate ] || source .venv/bin/activate
_cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
