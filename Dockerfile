FROM golang:1.22.5-alpine3.20 AS builder
ADD . $GOPATH/src/github.com/chechiachang/sc-stat
WORKDIR $GOPATH/src/github.com/chechiachang/sc-stat
RUN GOPATH=$PWD go install ./cmd/sc-stat

FROM alpine:3.20.2
RUN apk add --no-cache ca-certificates
RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /go/src/github.com/chechiachang/sc-stat/bin/sc-stat /usr/local/bin/sc-stat

ENTRYPOINT ["/usr/local/bin/sc-stat"]
CMD []
