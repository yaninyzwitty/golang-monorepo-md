
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

# Set workspace root
WORKDIR /workspace

# Copy go.work files first (better caching)
COPY go.work go.work.sum ./

# Copy all modules used by the gateway service
COPY gen/ ./gen/
COPY packages/shared/ ./packages/shared/
COPY packages/gateway/ ./packages/gateway/
COPY packages/devices/ ./packages/devices/

# Build the gateway binary
WORKDIR /workspace/packages/gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./main.go



FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy the built binary
COPY --from=builder /workspace/packages/gateway/gateway .

# Copy config file (if using one)
COPY --from=builder /workspace/packages/gateway/config.yaml ./config.yaml

# Expose gRPC/HTTP ports as needed
EXPOSE 8080

# Run the service
CMD ["./gateway"]
