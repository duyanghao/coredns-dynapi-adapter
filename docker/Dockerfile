# Building stage
FROM golang:1.13-alpine3.10

RUN apk add --update make

WORKDIR /go/src/github.com/duyanghao/coredns-dynapi-adapter

COPY . /go/src/github.com/duyanghao/coredns-dynapi-adapter

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GO111MODULE=on

ENV TIMEZONE "Asia/Shanghai"

RUN make build
RUN mv build/pkg/cmd/coredns-dynapi-adapter/coredns-dynapi-adapter /go/bin

# Production stage
FROM alpine:3.10
WORKDIR /go/bin

LABEL maintainer="duyanghao <1294057873@qq.com>"

# copy the go binaries from the building stage
COPY --from=0 /go/bin /go/bin

COPY etc/coredns-dynapi-adapter.yaml /etc/coredns-dynapi-adapter/coredns-dynapi-adapter.yaml

EXPOSE 8081

ENTRYPOINT ["./coredns-dynapi-adapter"]
