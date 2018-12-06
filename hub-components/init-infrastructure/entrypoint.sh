#!/bin/bash

#Example:
#
# docker run --rm \
#   -e GIT_CLONE_URL='https://github.com/qcloud2018/go-demo.git' \
#   -e GIT_REPO_DIR=/go/src/github.com/qcloud2018/go-demo \
#   hub.tencentyun.com/workshop/go-analysis
#

set -e
export SCRIPTDIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

source $SCRIPTDIR/component-base/libs.sh

do_task() {
}

do_task
