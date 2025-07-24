FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN go build -o app ./
RUN wget -q -O /grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/latest/download/grpc_health_probe-linux-amd64 && \
    chmod +x /grpc_health_probe

FROM alpine
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /grpc_health_probe /bin/grpc_health_probe

ENTRYPOINT ["./app"]
