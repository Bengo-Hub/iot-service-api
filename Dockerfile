# syntax=docker/dockerfile:1
# Build from service directory root: docker build -f Dockerfile -t iot-api:local .

FROM golang:1.23-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN GOTOOLCHAIN=auto go mod download
COPY . .
# Build all binaries: api, migrate, seed
RUN GOTOOLCHAIN=auto CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/iot-api ./cmd/api && \
    GOTOOLCHAIN=auto CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/iot-migrate ./cmd/migrate && \
    GOTOOLCHAIN=auto CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/iot-seed ./cmd/seed

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=builder /out/iot-api /usr/local/bin/iot-api
COPY --from=builder /out/iot-migrate /usr/local/bin/iot-migrate
COPY --from=builder /out/iot-seed /usr/local/bin/iot-seed
COPY internal/ent/migrate/migrations ./internal/ent/migrate/migrations
COPY scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh
RUN mkdir -p ./config/certs
USER app
EXPOSE 4006
ENV PORT=4006
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
