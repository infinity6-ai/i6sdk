#!/bin/bash -xe

# export GIT_SSH_COMMAND="${GIT_SSH_COMMAND:-ssh}"

# function cmd_build() {
#   ./build.sh work_init
#   ./build.sh tidy
#   rm -v go.work go.work.sum || true
#   ./build.sh gomod_upgrade
#   ./build.sh update
#   mkdir -p gen
#   if ! i6dev code is-clean; then
#     touch gen/dever.flag
#     ./build.sh version-inc
#   fi 
# }

# function cmd_test() {
#   ./meta/ci/ci.sh test
# }

# function cmd_release() {
#   ./meta/ci/ci.sh release
# }

# cd "$(dirname "$0")/../.."; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"

true
