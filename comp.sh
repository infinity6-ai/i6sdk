#!/bin/bash -e

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

function cmd_tidy() {
  go mod tidy
}

function cmd_update() {
  cmd_tidy
  cmd_codegen
}

function cmd_gomod_upgrade() {
  i6dev golang gomod-update-i6
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

function cmd_i6gox_set_version() {
  local _version="${1?'_version'}"
  sed -i "s/\(github.com\/infinity6-ai\/gox\/\w\+\)\s\+v[0-9][a-zA-Z0-9\.\-_]\+/\1 $_version/g" go.mod
  go mod tidy
}

function cmd_i6sdk_set_version() {
  local _version="${1?'_version'}"
  sed -i "s/\(github.com\/infinity6-ai\/i6sdk\/\w\+\)\s\+v[0-9][a-zA-Z0-9\._-]\+/\1 $_version/g" go.mod
  go mod tidy
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

function cmd_release_exists() {
  local _version="${1:-$(../comp.sh versionz run version)}"
  [ ! -z "$_version" ]
  git fetch --tags origin 2>/dev/null || git fetch --tags 2>/dev/null || true
  git rev-parse -q --verify "refs/tags/$_comp/$_version" >/dev/null 2>&1 || \
    git rev-parse -q --verify "refs/tags/$_comp/$_version^{commit}" >/dev/null 2>&1 || \
    git tag -l | grep -q "^$_comp/$_version$" || \
    [ -n "$(git ls-remote --tags origin "refs/tags/$_comp/$_version" 2>/dev/null)" ]
}

function cmd_release() {
  local _version="${1?'_version'}"
  [ ! -z "$_version" ]
  [ -z "$(git status -s "$@")" ]
  if git rev-parse --verify '@{u}' >/dev/null 2>&1; then
    [ "x0" == "x$(git rev-list --count '@{u}..HEAD')" ]
  fi
  local _tag="$_comp/$_version"
  if cmd_release_exists "$_version"; then
    echo "tag already exists: $_tag" 1>&2
  else
    cmd_i6sdk_set_version "$_version"
    # local _go_base_path="$(_go_base_path)"
    # local _go_dep_mod=""
    # GOWORK=off go list -mod=readonly -m "$_go_base_path/..." | \
    #   grep "^$_go_base_path/.*\ v" | \
    #   cut -d' ' -f1 | while read _go_dep_mod; do
    #   go mod edit -require="${_go_dep_mod}@${_version}"
    # done
    # go mod tidy
    [ "x$(go run . version)" == "x$_version" ]
    if [ -n "$(git status -s)" ]; then
      git add .
      git commit -m "release: $_tag"
      git push origin HEAD 2>/dev/null || git push || true
    fi
    git tag "$_tag"
    git push origin "$_tag"
  fi
}

function cmd_gomod_replace_on() {
  local _go_base_path="$(_go_base_path)"
  local _other=""
  for _other in $(ls -d ../*/ 2>/dev/null | sed 's|^\.\./||; s|/$||'); do
    local _mod="${_go_base_path}/$_other"
    if [ "$_other" != "$_comp" ] && [ -f "../$_other/go.mod" ]; then
      go mod edit -replace="${_mod}=../$_other"
    fi
  done
}

function cmd_gomod_replace_off() {
  local _go_base_path="$(_go_base_path)"
  local _other=""
  for _other in $(ls -d ../*/ 2>/dev/null | sed 's|^\.\./||; s|/$||'); do
    local _mod="${_go_base_path}/$_other"
    if [ "$_other" != "$_comp" ]; then
      go mod edit -dropreplace="${_mod}"
    fi
  done
}

function _go_base_path() {
  local _p="$(head -n 1 go.mod 2>/dev/null | awk '{print $2}' | rev | cut -d'/' -f2- | rev)"
  if [ -n "$_p" ]; then
    echo "$_p"
  else
    GOWORK=off go list -m | rev | cut -d'/' -f2- | rev
  fi
}

_comp="${1?"comp is required"}"; shift; 
_cmd="${1?"cmd is required"}"; shift; 
cd "$(dirname "$0")/$_comp"; 
"cmd_${_cmd}" "$@"
