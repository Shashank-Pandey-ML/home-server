# Multi-stage Dockerfile for Go microservices
# This Dockerfile is designed to be generic and reusable across multiple Go services

# Build stage
FROM golang:1.24.4-alpine AS builder

# Install necessary packages for building
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy common module first (for local dependencies)
COPY common/ ./common/

# Copy auth-service files
COPY auth-service/ ./auth-service/

# Set working directory to auth-service for build
WORKDIR /app/auth-service

# Download dependencies
RUN go mod download

# Build the application
# The binary name will be 'service' to keep it generic
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o service ./app/main.go

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/auth-service/service .

# Copy configuration files (optional - can be mounted as volume)
COPY auth-service/config.yaml ./

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port (default 8080, can be overridden)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
CMD ["./service"]
