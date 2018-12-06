FROM golang:1.11.2-alpine3.8
#FROM golang:1.11.2-stretch

LABEL TencentHubComponent='{"description": "Golang Golint Analysis, 用以对Golang编写的程序进行语法规则和代码风格检测.","input":[{"name": "GIT_CLONE_URL", "desc":"源代码地址"},{"name":"GIT_REF","desc":"git 提交对象，可以是一个git branch, git tag 或者git commit ID, 默认值master"},{"name":"GIT_REPO_DIR","desc":"git clone到本地的目录"}], "output":[]}'

WORKDIR /

ENTRYPOINT /entrypoint.sh

RUN apk add --no-cache git bash openssh-client && go get -u github.com/golang/lint/golint && mkdir /root/.ssh && chmod 0700 /root/.ssh

ADD component-base/libs.sh /component-base/
ADD ./*.sh /

