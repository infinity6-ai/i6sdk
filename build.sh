#!/bin/bash -e

eval "$(i6dev meta debug i6gox-build I6DEV_DEBUG)"

function cmd_comp_list() {
  find . -maxdepth 2 -name go.mod | cut -d'/' -f2
}

function cmd_comps_run() {
  local _k=""
  cmd_comp_list | while read _k; do
    ./"$_k"/build.sh "$@"
  done
}

function cmd_clean() {
  local _k=""
  find -name target | while read _k; do
    rm -rf "$_k"
  done
}


function cmd_update() {
  i6dev golang auth
  I6DEV_GOLANG_AUTH_DISABLED="true" cmd_comps_run update "$@"
}

function cmd_test() {
  cmd_comps_run test "$@"
}

function cmd_test_remote() {
  cmd_comps_run test_remote "$@"
}

function cmd_test_all() {
  cmd_comps_run test_all "$@"
}

function cmd_release() {
  cmd_comps_run release "$@"
}

function cmd_update_i6() {
  cmd_comps_run update_i6 "$@"
}

function cmd_version-inc() {
  i6dev dever version-inc sdk-web/version/version.go
}


cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"

