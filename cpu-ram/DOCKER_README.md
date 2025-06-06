# Go Benchmarking Container

This is a containerized version of the Go Benchmarking application designed for generating CPU and memory load to benchmark containerized environments.

## Prerequisites

- Docker
- Docker Compose (optional, for easier management)
- Make (optional, for convenience commands)

## Quick Start

### Using Make (Recommended)

This project includes a Makefile for convenience. To see all available commands:

```bash
make help
```

#### Start the container with Docker directly:

```bash
# Build the Docker image
make build

# Run the container
make run

# Check if the server is up
curl http://localhost:8080/health
```

#### Start the container with Docker Compose:

```bash
# Start with Docker Compose
make up

# Check logs
make logs
```

### Running benchmarks

The Makefile includes various convenient targets to interact with the benchmark API:

```bash
# Check the current status of all benchmarks
make status

# Start CPU benchmark using all cores
make cpu-start

# Start CPU benchmark using 2 cores
make cpu-start-2

# Start memory benchmark with 512MB limit
make mem-start-512

# Stop CPU benchmark
make cpu-stop

# Stop memory benchmark (memory remains allocated)
make mem-stop

# Explicitly free allocated memory
make mem-free
```

### Complex Benchmark Scenarios

```bash
# Run both CPU and memory benchmarks
make cpu-mem-test

# Stop all benchmarks and free memory
make stop-clean

# Run CPU benchmark with 512MB memory allocation
make cpu-test-512

# Run high load test (all CPU cores + 800MB memory)
make high-load
```

## Manual Docker Usage

If you don't want to use the Makefile:

```bash
# Build the image
docker build -t go-benchmark .

# Run the container
docker run -d --name go-benchmark -p 8080:8080 go-benchmark

# Stop and remove the container
docker stop go-benchmark
docker rm go-benchmark
```

## Manual Docker Compose Usage

```bash
# Start container
docker-compose up -d

# View logs
docker-compose logs -f

# Stop container
docker-compose down
```

## Benchmarking API Endpoints

### Basic endpoints
- `GET /` - Returns "Hello, World!"
- `GET /health` - Returns "Server is up and running!"
- `GET /status` - Returns the status of all benchmark tasks

### CPU Benchmark
- `POST /cpu/activate` - Starts the CPU benchmark task using all available cores
- `POST /cpu/activate/{n}` - Starts the CPU benchmark task using n cores (e.g., `/cpu/activate/2` uses 2 cores)
- `POST /cpu/deactivate` - Stops the CPU benchmark task

### Memory Benchmark
- `POST /memory/activate` - Starts the memory benchmark task with default 1GB limit
- `POST /memory/activate/{n}` - Starts the memory benchmark task with n MB limit (e.g., `/memory/activate/512` uses 512MB)
- `POST /memory/deactivate` - Stops the memory benchmark task (memory remains allocated)
- `POST /memory/free` - Explicitly releases allocated memory back to the system

## Container Resource Limits

To run with resource constraints, you can modify the `docker-compose.yml` file by uncommenting and adjusting the corresponding settings:

```yaml
services:
  benchmark:
    # CPU limits
    cpu_count: 2          # Number of CPUs
    cpus: 2.0             # Portion of CPU resources (2 CPUs)
    cpu_shares: 1024      # Relative CPU share (1024 = 1 CPU)
    
    # Memory limits
    mem_limit: 2g         # Maximum amount of memory the container can use
    memswap_limit: 2g     # Total amount of memory + swap
```

Alternatively, if using Docker directly:

```bash
docker run -d --name go-benchmark \
  -p 8080:8080 \
  --cpus=2 \
  --memory=2g \
  go-benchmark
```

## Benchmarking Workflow Example

```bash
# 1. Start the container
make up

# 2. Wait for the server to be up
make wait-for-server

# 3. Start CPU benchmark with 2 cores
make cpu-start-2

# 4. Start memory benchmark with 512MB
make mem-start-512

# 5. Check status
make status

# 6. Stop CPU benchmark
make cpu-stop

# 7. Wait and observe memory still allocated
make status

# 8. Free the memory
make mem-free

# 9. Stop the container
make down
```

## Troubleshooting

- If the container doesn't start, check the logs with `make logs` or `docker-compose logs`
- If endpoints are not responding, ensure the container is running with `docker ps`
- If memory isn't being released, explicitly call the memory free endpoint: `curl -X POST http://localhost:8080/memory/free` 