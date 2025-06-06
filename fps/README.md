# Benchmarking Container

This repository contains containerized benchmarking applications for system performance testing:

- **FPS Benchmarking**: Simulates frames per second (FPS) metrics using a sine wave pattern and exports them via OpenTelemetry

## Prerequisites

- Docker
- Docker Compose (optional, for easier management)
- An OpenTelemetry collector configured to listen on a Unix domain socket (if you want to collect metrics)
- Local Docker registry (optional, if you want to push images to a registry)

## Quick Start

### Using Docker

1. Build the Docker image:
   ```bash
   make build
   ```

2. Run the container:
   ```bash
   make run
   ```

3. View logs:
   ```bash
   make logs
   ```

4. Stop the container:
   ```bash
   make stop
   ```

### Using Docker Compose

1. Start the container:
   ```bash
   make up
   ```

2. Stop the container:
   ```bash
   make down
   ```

### Using with Local Registry

1. Build and push the image to a local registry at localhost:5000:
   ```bash
   make push-registry
   ```

2. Pull and run from the registry:
   ```bash
   docker pull localhost:5000/fps-benchmarking
   docker run -d --name fps-benchmarking -v /metrics:/metrics localhost:5000/fps-benchmarking
   ```

## Configuration

The FPS benchmarking application can be configured with the following environment variables:

- `OTEL_SOCKET_PATH`: Path to the Unix domain socket for OpenTelemetry metrics (default: `/metrics`)
- `HOSTNAME`: Container hostname for metrics identification

## Metrics

The FPS benchmarking application exports metrics in Prometheus text format directly on the Unix socket:

- **Socket location**: `$OTEL_SOCKET_PATH/$HOSTNAME`
- **Protocol**: HTTP over Unix domain socket
- **Metrics endpoint**: `/metrics` (standard Prometheus endpoint)

### Accessing Metrics

You can access the metrics using any HTTP client that supports Unix domain sockets. For example:

```bash
# Using curl with Unix socket support
curl --unix-socket /metrics/fps-host http://localhost/metrics
```

### Metrics Details

The application exports the following metrics:

- `service.fps`: The simulated frames per second (gauge)
  - Includes a `state` label with values `"in"` and `"out"`
  - FPS follows a sine wave pattern between 30 and 120 FPS with a 2-minute cycle

## Development

To build and run the application locally without Docker:

```bash
make build-local
make run-local
```

## Directory Structure

- `fps/`: FPS benchmarking application source code
- `Dockerfile`: Container definition for the benchmarking application
- `Makefile`: Build and management commands
- `docker-compose.yml`: Docker Compose configuration 