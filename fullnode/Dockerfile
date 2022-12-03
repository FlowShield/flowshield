FROM golang:1.17.8-alpine AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build -o flowshield .


FROM alpine AS final
WORKDIR /app
COPY --from=builder /build/flowshield /app/
COPY --from=builder /build/config.yaml /app/

ENTRYPOINT ["/app/flowshield"]
