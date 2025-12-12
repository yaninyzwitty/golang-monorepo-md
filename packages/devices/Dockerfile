# Stage 1: Build
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory to workspace root
WORKDIR /workspace

COPY go.work go.work.sum ./

# Copy gen module (generated protobufs)
COPY gen/ ./gen/

COPY packages/shared/ ./packages/shared/
# Copy gateway-service module (needed for go.work)
COPY packages/gateway/ ./packages/gateway/
# Copy product-service module
COPY packages/devices/ ./packages/devices/

WORKDIR /workspace/packages/devices
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o devices ./main.go


# Stage 2: Runtime
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /workspace/packages/devices/devices .

# Copy config file if needed
COPY --from=builder /workspace/packages/devices/config.yaml ./config.yaml

# Expose gRPC port
EXPOSE 50051

# Run the service
CMD ["./devices"]