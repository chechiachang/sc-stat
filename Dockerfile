# syntax=docker/dockerfile:1.9.0
FROM golang:1.22.5-alpine3.20 AS builder

ADD . $GOPATH/src/github.com/chechiachang/sc-stat
WORKDIR $GOPATH/src/github.com/chechiachang/sc-stat
RUN GOPATH=$PWD go install ./cmd/sc-stat

FROM alpine:3.20.2

LABEL org.opencontainers.image.authors="chechiachang999@gmail.com"
LABEL org.opencontainers.image.title="Sport Center Stat Fetcher"
LABEL org.opencontainers.image.source="https://github.com/chechiachang/sc-stat"
LABEL org.opencontainers.image.vendor="Che-Chia Chang"
LABEL org.opencontainers.image.base.name="alpine:3.20.2"
RUN apk add --no-cache ca-certificates git openssh

RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /go/src/github.com/chechiachang/sc-stat/bin/sc-stat /usr/local/bin/sc-stat

ARG APP_CODE_VERSION
ENV APP_CODE_VERSION=${APP_CODE_VERSION}

ENTRYPOINT ["/usr/local/bin/sc-stat"]
CMD []
