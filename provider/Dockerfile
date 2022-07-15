FROM golang:1.17.8-alpine AS builder

WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 go build -o provider ./cmd

FROM ubuntu:20.04

WORKDIR /root

RUN apt update && apt install -y curl

COPY --from=builder /build/backend .
COPY --from=builder /build/configs ./configs
RUN chmod +x provider

CMD ["./provider"]