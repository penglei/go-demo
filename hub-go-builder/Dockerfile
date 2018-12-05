#FROM hub.tencentyun.com/workshop/go-builder:latest
FROM golang:1.11.2-stretch

ARG REPO_DIR

ADD go.mod go.sum $REPO_DIR/

RUN go get -u golang.org/x/vgo && cd $REPO_DIR && vgo mod download

