# stage 1: build src code to binary
FROM golang:1.16-buster as builder
ENV GO11MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR $GOPATH/src/github.com/lushenle/sendmail
COPY . .

RUN make

# stage 2: use alpine as base image
FROM alpine:3.10
LABEL maintainer="Shnele Lu <lushenle@gmail.com>" \
    app=sendmail \
    version=v1.0

RUN apk update && \
    apk --no-cache add tzdata ca-certificates && \
    cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

COPY --from=builder /go/src/github.com/lushenle/sendmail/app /
ENV MAIL_CONFIG=/config.json
CMD ["/app"]
