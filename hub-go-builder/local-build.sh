#!/bin/bash

export SCRIPTDIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
export PROJECTDIR=$(cd $SCRIPTDIR/.. && pwd)

image_tag=hub.tencentyun.com/workshop/go-demo-cache-builder:latest
project_container_dir=/project

docker build --build-arg REPO_DIR=$project_container_dir -t $image_tag -f $SCRIPTDIR/Dockerfile "$PROJECTDIR"
#docker push $image_tag

