#!/bin/bash -xe

export GIT_SSH_COMMAND="${GIT_SSH_COMMAND:-ssh}"

function cmd_image_download() {
  docker pull us-central1-docker.pkg.dev/i6-rs-contint/i6devhub/i6dev-full:latest
  docker tag us-central1-docker.pkg.dev/i6-rs-contint/i6devhub/i6dev-full:latest infinity6/current:dev
}

function cmd_build() {
  i6dev golang auth-meta
  ./build.sh uphead
  ./build.sh work_init
  # ./build.sh gomod_replace_on
  ./build.sh tidy
  ./build.sh codegen
  # ./build.sh gomod_replace_off
}

function cmd_test() {
  ./build.sh test
}

function cmd_release() {
  rm -v go.work go.work.sum || true 
  i6dev dever git-auth
  i6dev code github-to-ssh origin || true
  local _branch="$(git branch --show-current 2>/dev/null || echo "master")"
  git branch --set-upstream-to="origin/$_branch" "$_branch" 2>/dev/null || true
  ./build.sh release_exists || ./build.sh release
}

function cmd_deploy() {
  true
}

cd "$(dirname "$0")/../.."; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"

