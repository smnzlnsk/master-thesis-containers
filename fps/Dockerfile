# Build stage
FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

# Accept build arguments for version
ARG VERSION=0.0.1

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code and version file
COPY . .

# Build the application with version information
RUN go build -ldflags "-X main.buildVersion=${VERSION}" -o fps-benchmark

# Final stage
FROM alpine:latest

# Accept build arguments for version (to use in labels)
ARG VERSION=0.0.1

# Add container metadata
LABEL version="${VERSION}" \
      maintainer="Benchmarking Container Team" \
      description="FPS monitoring container"

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory in the container
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/fps-benchmark .

# Create directory for metrics socket
RUN mkdir -p /metrics

# Set environment variables
ENV OTEL_SOCKET_PATH=/metrics

# Expose port 80 for the hostname HTTP server
EXPOSE 80

# Set the default command to run the application
CMD ["./fps-benchmark"] 