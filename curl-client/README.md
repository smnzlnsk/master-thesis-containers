# Curl Client Docker Container

A simple Docker container that continuously curls a given endpoint every 50ms for benchmarking and testing purposes.

## Features

- Lightweight Alpine Linux base image
- Configurable target URL via environment variable or command line argument
- Curls endpoint every 50ms (20 requests per second)
- Shows response status codes and response bodies
- CSV logging to `/metrics/curl_metrics.csv` with timestamp, host, return code, and return body
- Graceful shutdown with request counter
- Built-in error handling and validation

## Building the Container

```bash
docker build -t curl-client .
```

## Usage

### Option 1: Using Environment Variable with Volume Mount
```bash
docker run -v $(pwd)/metrics:/metrics -e TARGET_URL="https://your-endpoint.com/api" curl-client
```

### Option 2: Using Command Line Argument with Volume Mount
```bash
docker run -v $(pwd)/metrics:/metrics curl-client "https://your-endpoint.com/api"
```

### Option 3: Using Default URL (httpbin.org) with Volume Mount
```bash
docker run -v $(pwd)/metrics:/metrics curl-client
```

**Note**: The `-v $(pwd)/metrics:/metrics` mount is required to persist the CSV log file on the host system.

## Example Output

```
Starting curl client for endpoint: https://httpbin.org/get
Curling every 50ms...
Press Ctrl+C to stop
Created CSV log file: /metrics/curl_metrics.csv
[1] Status: 200 | Body: {"args":{},"headers":{"Host":"httpbin.org"},"origin":"1.2.3.4","url":"https://httpbin.org/get"}
[2] Status: 200 | Body: {"args":{},"headers":{"Host":"httpbin.org"},"origin":"1.2.3.4","url":"https://httpbin.org/get"}
[3] Status: 200 | Body: {"args":{},"headers":{"Host":"httpbin.org"},"origin":"1.2.3.4","url":"https://httpbin.org/get"}
...
```

## Stopping the Container

Press `Ctrl+C` or send a SIGTERM signal to gracefully stop the container. It will display the total number of requests sent.

## CSV Logging

The container logs all requests to `/metrics/curl_metrics.csv` with the following format:

```csv
timestamp,host,return_code,return_body
1705316445123456,httpbin.org,200,"{"args":{},"headers":{"Host":"httpbin.org"}}"
1705316445173789,httpbin.org,200,"Simple response text"
1705316445223012,example.com,ERROR,"Request failed"
```

- **timestamp**: Microseconds from epoch (Unix timestamp)
- **host**: Extracted hostname from the target URL
- **return_code**: HTTP status code or "ERROR" for failed requests
- **return_body**: Response body (single line, commas replaced with semicolons to avoid CSV conflicts)

## Configuration

- **Frequency**: Currently set to 50ms intervals (20 RPS)
- **Timeout**: 5 seconds per request
- **Default URL**: http://httpbin.org/get (for testing)
- **Log File**: `/metrics/curl_metrics.csv` (requires volume mount)

## Use Cases

- Load testing endpoints
- Network connectivity monitoring
- API response time monitoring
- Benchmarking web services 