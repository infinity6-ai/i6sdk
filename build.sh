#!/bin/bash -e

[ "$I6DEV_DEBUG" != "true" ] || eval "$(i6dev meta debug i6sdk-build I6DEV_DEBUG)"

function cmd_comps_list() {
  echo "versionz"
  echo "sdkcorez"
  echo "apiuserz"
  echo "apidatasetz"
  echo "apiingestz"
  echo "apipipez"
  echo "apitfz"
  echo "apicollectionz"
  echo "mold"
  echo "forgez"
  echo "sdkgo"
  echo "sdkweb"
  echo "sdkpy"
}

function cmd_comps_exec() {
  local _k=""
  cmd_comps_list | while read _k; do
    (cd "$_k" && "$@")
  done
}

function cmd_comps_run() {
  local _k=""
  cmd_comps_list | while read _k; do
    ./comp.sh "$_k" "$@"
  done
}

function cmd_work_init() {
  [ ! -f go.work.sum ] || rm go.work.sum
  [ ! -f go.work ] || rm go.work
  cat <<EOF > go.work
go 1.26.1

use (
$(cmd_comps_list | sed 's/^/\t.\//')
)
EOF
}

function cmd_clean() {
  cmd_comps_run clean "$@"
}

function cmd_go() {
  cmd_comps_run go "$@"
}

function cmd_uphead() {
  cmd_comps_run uphead "$@"
}

function cmd_update() {
  cmd_comps_run update "$@"
}

function cmd_codegen() {
  cmd_comps_run codegen "$@"
}

function cmd_test() {
  cmd_comps_run test "$@"
}

function cmd_gomod_upgrade() {
  cmd_comps_run gomod_upgrade "$@"
}

function cmd_test_web() {
  cmd_comps_run test_web "$@"
}

function cmd_fmt() {
  cmd_comps_run fmt "$@"
}

function cmd_force_delete_version() {
  local _version="${1?'_version'}"
  [ ! -z "$_version" ]
  cmd_comps_list | while read _k; do
    git tag -d "$_k/$_version" || true
    git push --delete origin "$_k/$_version" || true
  done
  git tag -d "$_version" || true
  git push --delete origin "$_version" || true
}

function cmd_release_exists() {
  local _version="${1:-$(./comp.sh versionz run version)}"
  [ ! -z "$_version" ]
  git fetch --tags origin 2>/dev/null || git fetch --tags 2>/dev/null || true
  git rev-parse -q --verify "refs/tags/$_version" >/dev/null 2>&1 || \
    git rev-parse -q --verify "refs/tags/$_version^{commit}" >/dev/null 2>&1 || \
    git tag -l | grep -q "$_version$" || \
    [ -n "$(git ls-remote --tags origin "refs/tags/$_version" 2>/dev/null)" ]
}

function cmd_release() {
  if find . -name '*.go' | xargs grep go.code.infinity6.ai; then false i6sdk cannot use private deps; fi
  if find . -name 'go.mod' | xargs grep go.code.infinity6.ai; then false i6sdk cannot use private deps; fi
  [ -z "$(git status -s "$@")" ]
  if git rev-parse --verify '@{u}' >/dev/null 2>&1; then
    [ "x0" == "x$(git rev-list --count '@{u}..HEAD')" ]
  fi
  local _version="$(./comp.sh versionz run version)"
  [ ! -z "$_version" ]
  if cmd_release_exists "$_version" >/dev/null; then
    echo "release already exists: $_version" 1>&2
    false
  fi
  GOWORK=off cmd_comps_run release "$_version"
  git tag "$_version"
  git push origin "$_version"
}

function cmd_version-inc() {
  ./comp.sh versionz run inc
  ./comp.sh versionz run version
}

function cmd_version() {
  ./comp.sh versionz run version
}

function cmd_infer-demo() {
  python ./sdkweb/src/site/infer-demo/serve.py
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
