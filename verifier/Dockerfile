FROM golang:1.17.8-alpine AS builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build -o verifier .


FROM alpine AS final
WORKDIR /app
COPY --from=builder /build/verifier /app/
COPY --from=builder /build/config.yaml /app/

ENTRYPOINT ["/app/verifier"]
