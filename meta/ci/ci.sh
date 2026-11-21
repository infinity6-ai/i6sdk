#!/bin/bash -xe

function cmd_image_download() {
  docker pull us-central1-docker.pkg.dev/i6-rs-contint/i6devhub/i6dev-full:latest
  docker tag us-central1-docker.pkg.dev/i6-rs-contint/i6devhub/i6dev-full:latest infinity6/current:dev
}

function cmd_build() {
  i6dev golang auth-meta
  ./build.sh update
}

function cmd_test() {
#   sudo mv /usr/bin/google-chrome /usr/bin/google-chrome.original
#   echo '#!/bin/bash -xe
#   exec /usr/bin/google-chrome.original --no-sandbox "$@" 
#   ' | sudo tee /usr/bin/google-chrome
#   sudo chmod +x /usr/bin/google-chrome
  ./build.sh test
}

function cmd_release() {
  # ./build.sh release_exists || ./build.sh release
  false it is not working
}

function cmd_deploy() {
  true
}

cd "$(dirname "$0")/../.."; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"

