FROM golang:1.17.8-alpine AS builder

WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 go build -o backend ./cmd/server

FROM ubuntu:20.04

WORKDIR /backend

COPY --from=builder /build/backend .
COPY --from=builder /build/configs/config_server.toml .
RUN chmod +x backend

CMD ["./backend", "-c", "config_server.toml", "up"]