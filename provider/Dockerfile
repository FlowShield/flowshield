FROM golang:1.17.8-alpine AS builder

WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 go build -o provider ./cmd

FROM ubuntu:20.04

WORKDIR /root

COPY --from=builder /build/backend .
COPY --from=builder /build/configs .
RUN chmod +x provider

CMD ["./provider", "-c", "./configs/config.toml"]