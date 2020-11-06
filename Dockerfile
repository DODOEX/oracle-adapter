FROM golang:1.15-alpine as builder

RUN apk add --no-cache make curl git gcc musl-dev linux-headers

ADD . /go/src/github.com/Dominator008/asset-price-oracle-adapter
RUN cd /go/src/github.com/Dominator008/asset-price-oracle-adapter && make build

# Copy into a second stage container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/Dominator008/asset-price-oracle-adapter/asset-price-oracle-adapter /usr/local/bin/

EXPOSE 8080
CMD ["asset-price-oracle-adapter"]