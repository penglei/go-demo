#FROM golang:1.11.2-stretch as builder
ARG GIT_REPO_DIR=/go/src/github.com/qcloud2018/go-demo
FROM hub.tencentyun.com/workshop/go-demo-cache-builder:latest as builder
ARG GIT_REPO_DIR
WORKDIR $GIT_REPO_DIR
COPY . $GIT_REPO_DIR
RUN vgo build -v -a -o /go-demo ./cmd/*.go

FROM debian:stretch
EXPOSE 8080
CMD ["/go-demo"]
ADD nsswitch.conf    /etc/
ARG GIT_REPO_DIR

COPY --from=builder /go-demo /
COPY --from=builder $GIT_REPO_DIR/migration/versions $GIT_REPO_DIR/migration/versions

