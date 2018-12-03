#!/bin/bash

export SCRIPTDIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
export PROJECTDIR=$(cd $SCRIPTDIR/.. && pwd)

cd "$SCRIPTDIR"

image_tag=hub.tencentyun.com/workshop/go-demo-cache-builder:latest

docker build --build-arg REPO_DIR=/project -t $image_tag "$PROJECTDIR"
#docker push $image_tag

