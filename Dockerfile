# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder
WORKDIR /src
# Copy shared auth-client first (needed for replace directive)
# Build context should be from workspace root: docker build -f iot-service/Dockerfile -t iot-service:local .
COPY shared/auth-client /shared/auth-client
COPY iot-service/go.mod iot-service/go.sum ./
RUN GOTOOLCHAIN=auto go mod download
COPY iot-service .

RUN GOTOOLCHAIN=auto CGO_ENABLED=0 go build -o /out/iot ./cmd/api

FROM alpine:3.20
RUN addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=builder /out/iot /app/service
# TLS certificates directory (optional, can be mounted as volume)
RUN mkdir -p ./config/certs
USER app
EXPOSE 4006
ENV PORT=4006
ENTRYPOINT ["/app/service"]

