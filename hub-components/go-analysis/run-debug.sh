#!/bin/bash

# --entrypoint bash -i
docker run --rm -t \
  -e GIT_REF='feature/framework' \
  -e GIT_CLONE_URL='https://github.com/qcloud2018/go-demo.git' \
  -e GIT_REPO_DIR=/go/src/github.com/qcloud2018/go-demo \
  hub.tencentyun.com/workshop/go-analysis

