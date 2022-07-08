FROM golang:1.17.8-alpine AS builder

ENV GO111MODULE=on
WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 go build -o ca .

FROM ubuntu:20.04

WORKDIR /root

COPY --from=builder /build/ca .
COPY --from=builder /build/conf.yml .
COPY --from=builder /build/config.json .
RUN chmod +x ca

# TLS service
CMD ["./ca", "tls"]