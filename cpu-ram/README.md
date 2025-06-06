# Go Benchmarking Server

A HTTP server built in Go for generating CPU and memory load for benchmarking containerized environments.

## Project Structure

```
├── main.go         # Entry point for the application
├── benchmark/      # Benchmark task implementation
│   ├── cpu.go      # CPU load generation
│   └── memory.go   # Memory load generation
├── config/         # Configuration package
│   └── config.go   # Server configuration
├── handlers/       # HTTP handlers
│   └── handlers.go # Request handler implementations
└── README.md       # This file
```

## Running the server

```bash
# Run the server
go run main.go
```

The server will start on port 8080. You can access it at [http://localhost:8080](http://localhost:8080).

## Endpoints

### Basic endpoints
- `/` - Returns "Hello, World!"
- `/health` - Returns "Server is up and running!"
- `/status` - GET endpoint that returns the status of all benchmark tasks

### CPU Benchmark
- `/cpu/activate` - POST endpoint that starts the CPU benchmark task using all available cores
- `/cpu/activate/{n}` - POST endpoint that starts the CPU benchmark task using n cores (e.g., `/cpu/activate/2` uses 2 cores)
- `/cpu/deactivate` - POST endpoint that stops the CPU benchmark task

### Memory Benchmark
- `/memory/activate` - POST endpoint that starts the memory benchmark task with default 1GB limit
- `/memory/activate/{n}` - POST endpoint that starts the memory benchmark task with n MB limit (e.g., `/memory/activate/512` uses 512MB)
- `/memory/deactivate` - POST endpoint that stops the memory benchmark task (memory remains allocated)
- `/memory/free` - POST endpoint that explicitly releases allocated memory back to the system

### Legacy endpoints (for backward compatibility)
- `/activate` - Same as `/cpu/activate`
- `/deactivate` - Same as `/cpu/deactivate`

## Benchmarking Functionality

This application is designed to generate high CPU and memory load for benchmarking containerized environments:

### CPU Benchmark
- Starts worker goroutines to perform complex mathematical operations
- You can specify how many CPU cores to utilize (from 1 to all available cores)
- Each worker continuously executes CPU-intensive calculations involving trigonometric functions, exponentials, square roots, and other operations
- The system efficiently utilizes the specified number of CPU cores to generate load
- Status updates are printed showing the number of calculations performed
- Runs indefinitely until explicitly stopped

### Memory Benchmark
- Continuously allocates memory in 10MB blocks
- Default memory limit is 1GB (1024MB), but can be configured via the API
- You can specify a custom memory limit in MB (e.g., 512MB)
- Writes data to the allocated memory to ensure it's not optimized away
- Stops allocating more memory when the limit is reached
- Displays the total amount of allocated memory and percentage of limit used
- Allocates a new block every 500ms to provide a controlled increase in memory usage
- **Important**: Memory remains allocated even after stopping the benchmark
- Memory is only released when explicitly calling the `/memory/free` endpoint
- This allows for measuring memory pressure over extended periods

## Package Organization

- `benchmark`: Contains all resource-intensive task management:
  - `cpu.go`: CPU-intensive task implementation
  - `memory.go`: Memory-intensive task implementation
- `config`: Handles application configuration
- `handlers`: HTTP request handlers for the API endpoints

## Container Usage

This application is designed to be containerized. Example Dockerfile:

```dockerfile
FROM golang:1.19-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o benchserver

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/benchserver .
EXPOSE 8080
CMD ["./benchserver"]
```

Build and run the container:

```bash
docker build -t go-benchmark .
docker run -p 8080:8080 go-benchmark
```

## Usage Examples

Start the CPU benchmark using all available cores:
```bash
curl -X POST http://localhost:8080/cpu/activate
```

Start the CPU benchmark using 2 cores:
```bash
curl -X POST http://localhost:8080/cpu/activate/2
```

Stop the CPU benchmark:
```bash
curl -X POST http://localhost:8080/cpu/deactivate
```

Start the memory benchmark with default 1GB limit:
```bash
curl -X POST http://localhost:8080/memory/activate
```

Start the memory benchmark with a 512MB limit:
```bash
curl -X POST http://localhost:8080/memory/activate/512
```

Stop the memory benchmark (memory remains allocated):
```bash
curl -X POST http://localhost:8080/memory/deactivate
```

Explicitly release allocated memory:
```bash
curl -X POST http://localhost:8080/memory/free
```

Check benchmark status:
```bash
curl http://localhost:8080/status
```

## Memory Management Workflow Example

This example demonstrates the intended memory management workflow:

```bash
# 1. Start memory benchmark with 512MB limit
curl -X POST http://localhost:8080/memory/activate/512

# 2. Wait for it to allocate memory...

# 3. Stop the benchmark (memory remains allocated)
curl -X POST http://localhost:8080/memory/deactivate

# 4. Verify memory is still allocated
curl http://localhost:8080/status
# Output shows: Memory Benchmark: STOPPED (stopped, but still holding 512 MB of memory)

# 5. When ready to release the memory
curl -X POST http://localhost:8080/memory/free

# 6. Verify memory has been released
curl http://localhost:8080/status
``` 