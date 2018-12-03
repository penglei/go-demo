FROM hub.tencentyun.com/workshop/go-builder:latest

ARG REPO_DIR

ADD go.mod go.sum $REPO_DIR/
RUN cd $REPO_DIR && vgo mod download

