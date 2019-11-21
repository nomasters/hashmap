FROM golang:1.13.4-alpine3.10 AS builder

ARG CGO_ENABLED=0

RUN apk update && apk add ca-certificates

WORKDIR /app
COPY . .

RUN go install

FROM alpine:3.10

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/bin/hashmap /usr/bin/hashmap
RUN addgroup -S hashmap && adduser -S hashmap -G hashmap
WORKDIR /home/hashmap

USER hashmap