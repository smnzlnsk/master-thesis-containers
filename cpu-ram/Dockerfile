# Build stage
FROM golang:1.19-alpine AS builder

# Set the working directory
WORKDIR /app

# Accept build arguments for version
ARG VERSION=0.0.1

# Copy go.mod and go.sum files
COPY go.mod ./

# Copy the source code
COPY . .

# Build the application with version information
RUN go build -ldflags "-X main.buildVersion=${VERSION}" -o benchserver

# Final stage
FROM alpine:latest

# Accept build arguments for version (to use in labels)
ARG VERSION=0.0.1

# Add container metadata
LABEL version="${VERSION}" \
      maintainer="Benchmarking Container Team" \
      description="CPU-RAM benchmarking container"

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory in the container
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/benchserver .

# Expose the port the app runs on
EXPOSE 80

# Command to run the executable
CMD ["./benchserver"] 