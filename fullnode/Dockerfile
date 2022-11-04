FROM golang:1.17.8-alpine AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build -o cloudslit .


FROM alpine AS final
WORKDIR /app
COPY --from=builder /build/cloudslit /app/
COPY --from=builder /build/config.yaml /app/

ENTRYPOINT ["/app/cloudslit"]
