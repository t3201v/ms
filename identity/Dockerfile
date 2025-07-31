FROM golang:1.24.5-alpine AS builder
WORKDIR /app
ENV GOBIN=/app/bin
COPY go.* .
RUN go mod download
COPY . .
RUN go install ./cmd/***

RUN #wget -q -O /grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/latest/download/grpc_health_probe-linux-amd64 && \
    chmod +x ./grpc_health_probe

FROM alpine
WORKDIR /app
COPY --from=builder /app/bin .
COPY --from=builder /app/grpc_health_probe /bin/grpc_health_probe
