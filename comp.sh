#!/bin/bash -e

export GIT_SSH_COMMAND="${GIT_SSH_COMMAND:-ssh}"

[ "$I6DEV_DEBUG" != "true" ] || eval "$(i6dev meta debug i6sdk-comp I6DEV_DEBUG)"

function cmd_clean() {
  local _k=""
  find -name target | while read _k; do
    rm -rf "$_k"
  done
}

function cmd_go() {
  go "$@"
}

function cmd_update() {
  go mod tidy
  cmd_codegen
}

function cmd_run() {
  go run . "$@"
}

function cmd_codegen() {
  go generate ./...
}

function cmd_test_unit() {
  go test -run '^TestUnit.*$' ./... "$@" 
  [ ! -f cmd_unit_test.sh ] || ./cmd_unit_test.sh all
}

function cmd_test_web() {
  go test -run '^TestWeb.*$' ./... "$@" 
}

function cmd_test_remote() {
  go test -run '^TestRemote.*$' ./... "$@" 
}

function cmd_test_manual() {
  go test -run '^TestManual.*$' ./... "$@" 
}

function cmd_test() {
  cmd_test_unit
  cmd_test_web
}

function cmd_fmt() {
  go fmt ./...
}

function cmd_uphead() {
  local _go_base_path="$(_go_base_path)"
  local _go_dep_mod=""
  GOWORK=off go list -mod=readonly -m "$_go_base_path/..." | \
    grep "^$_go_base_path/.*\ v" | \
    cut -d' ' -f1 | while read _go_dep_mod; do
    go get "${_go_dep_mod}@main"
  done
}

function cmd_release() {
  local _version="${1?'_version'}"
  [ ! -z "$_version" ]
  [ -z "$(git status -s "$@")" ]
  if git rev-parse --verify '@{u}' >/dev/null 2>&1; then
    [ "x0" == "x$(git rev-list --count '@{u}..HEAD')" ]
  fi
  local _tag="$_comp/$_version"
  if git show-ref --tags "$_tag" --quiet; then
    echo "tag already exists: $_tag" 1>&2
    false
  fi
  local _go_base_path="$(_go_base_path)"
  local _go_dep_mod=""
  GOWORK=off go list -mod=readonly -m "$_go_base_path/..." | \
    grep "^$_go_base_path/.*\ v" | \
    cut -d' ' -f1 | while read _go_dep_mod; do
    go mod edit -require="${_go_dep_mod}@${_version}"
  done
  go mod tidy
  [ "x$(go run . version)" == "x$_version" ]
  if [ -n "$(git status -s)" ]; then
    git add .
    git commit -m "release: $_tag"
    git push origin HEAD 2>/dev/null || git push || true
  fi
  git tag "$_tag"
  git push origin "$_tag"
}

function _go_base_path() {
  GOWORK=off go list -m | rev | cut -d'/' -f2- | rev
}

_comp="${1?"comp is required"}"; shift; 
_cmd="${1?"cmd is required"}"; shift; 
cd "$(dirname "$0")/$_comp"; 
"cmd_${_cmd}" "$@"
