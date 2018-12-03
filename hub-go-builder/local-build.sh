#!/bin/bash

export SCRIPTDIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
export PROJECTDIR=$(cd $SCRIPTDIR/.. && pwd)

cd $SCRIPTDIR

cp $PROJECTDIR/go.sum $PROJECTDIR/go.mod .

docker build --build-arg REPO_DIR=/project -t hub.tencentyun.com/workshop/go-demo-cache-builder:latest .
#docker push hub.tencentyun.com/workshop/go-demo-cache-builder:latest

