FROM golang:1.11.2-alpine3.8 as builder
WORKDIR /go/src/github.com/qcloud2018/go-demo
COPY . /go/src/github.com/qcloud2018/go-demo
RUN go build -v -a -o /go-demo

FROM alpine:3.8
EXPOSE 8080
CMD ["/go-demo"]
ADD nsswitch.conf    /etc/

COPY --from=builder /go-demo /


