# Prometheus to CSV Converter

A comprehensive bash script that collects metrics from Prometheus endpoints and converts them to CSV format for analysis and reporting.

## Features

- ✅ **Flexible Collection**: Collect metrics at custom intervals and durations
- ✅ **Filtering**: Filter metrics by name patterns using regex
- ✅ **Multiple Formats**: Include/exclude timestamps, customize CSV output
- ✅ **Error Handling**: Robust error handling and validation
- ✅ **Verbose Logging**: Detailed logging with color-coded output
- ✅ **Dependency Checking**: Automatic validation of required tools
- ✅ **Signal Handling**: Graceful shutdown with Ctrl+C
- ✅ **Metric Discovery**: List available metrics from endpoints

## Requirements

The script requires the following tools (automatically checked):
- `curl` - For HTTP requests
- `awk` - For text processing
- `grep` - For pattern matching
- `sort` - For sorting output

## Installation

1. Download the script:
```bash
wget https://raw.githubusercontent.com/your-repo/prometheus-to-csv.sh
# or
curl -O https://raw.githubusercontent.com/your-repo/prometheus-to-csv.sh
```

2. Make it executable:
```bash
chmod +x prometheus-to-csv.sh
```

## Usage

### Basic Syntax
```bash
./prometheus-to-csv.sh -e ENDPOINT [OPTIONS]
```

### Required Parameters
- `-e, --endpoint URL` - Prometheus metrics endpoint URL

### Optional Parameters
- `-o, --output FILE` - Output CSV file (default: metrics_TIMESTAMP.csv)
- `-i, --interval SECONDS` - Collection interval in seconds (default: 10)
- `-d, --duration SECONDS` - Total collection duration in seconds (default: infinite)
- `-f, --filter PATTERN` - Filter metrics by name pattern (regex supported)
- `-t, --no-timestamp` - Don't include timestamp column in CSV
- `-m, --help-metrics` - Show available metrics from endpoint and exit
- `-v, --verbose` - Enable verbose output
- `-h, --help` - Show help message

## Examples

### 1. Basic Collection
```bash
# Collect metrics every 10 seconds (default) indefinitely
./prometheus-to-csv.sh -e http://localhost:8080/metrics -v
```

### 2. Timed Collection
```bash
# Collect metrics every 5 seconds for 5 minutes
./prometheus-to-csv.sh -e http://localhost:8080/metrics -i 5 -d 300 -o my_metrics.csv -v
```

### 3. Filtered Collection
```bash
# Only collect CPU and memory related metrics
./prometheus-to-csv.sh -e http://localhost:8080/metrics -f "cpu|memory|ram" -v
```

### 4. Discover Available Metrics
```bash
# List all available metrics from the endpoint
./prometheus-to-csv.sh -e http://localhost:8080/metrics --help-metrics
```

### 5. Collection Without Timestamps
```bash
# Generate CSV without timestamp columns
./prometheus-to-csv.sh -e http://localhost:8080/metrics -t -o clean_metrics.csv
```

## Output Format

### With Timestamps (Default)
```csv
timestamp,epoch,metric_name,labels,value
2024-01-15 10:30:00,1705312200,"cpu_usage_percent","instance=""localhost"",job=""test""",45.67
2024-01-15 10:30:00,1705312200,"memory_usage_bytes","instance=""localhost"",job=""test""",1536
```

### Without Timestamps
```csv
metric_name,labels,value
"cpu_usage_percent","instance=""localhost"",job=""test""",45.67
"memory_usage_bytes","instance=""localhost"",job=""test""",1536
```

## Supported Metric Types

The script handles all standard Prometheus metric types:
- **Counter** - Monotonically increasing values
- **Gauge** - Values that can go up and down
- **Histogram** - Sampling observations (buckets, count, sum)
- **Summary** - Similar to histogram with quantiles

## Testing

### Using the Mock Server

The repository includes a mock Prometheus server for testing:

1. Start the mock server:
```bash
python3 mock-prometheus-server.py
```

2. In another terminal, test the script:
```bash
./prometheus-to-csv.sh -e http://localhost:8080/metrics -v -d 30
```

The mock server provides sample metrics:
- `cpu_usage_percent` - Oscillating CPU usage values
- `memory_usage_bytes` - Random memory usage values  
- `fps_current` - Oscillating FPS values
- `http_requests_total` - Counter metrics with labels
- `process_start_time_seconds` - Process start time gauge

### Real-World Testing

Test with actual Prometheus endpoints:
```bash
# Test with Prometheus itself
./prometheus-to-csv.sh -e http://localhost:9090/metrics -f "prometheus_" -d 60

# Test with Node Exporter
./prometheus-to-csv.sh -e http://localhost:9100/metrics -f "node_cpu|node_memory" -d 120

# Test with custom application
./prometheus-to-csv.sh -e http://your-app:8080/metrics -v
```

## Error Handling

The script includes comprehensive error handling:

- **Connection Errors**: Validates endpoint accessibility
- **Timeout Handling**: 30-second timeout for HTTP requests
- **Dependency Checking**: Verifies required tools are installed
- **Signal Handling**: Graceful shutdown on Ctrl+C
- **Data Validation**: Checks for empty responses

## Performance Considerations

- **Memory Usage**: Processes metrics in streaming fashion
- **Network Efficiency**: Single HTTP request per collection cycle
- **File I/O**: Appends to CSV file incrementally
- **CPU Usage**: Minimal processing overhead with awk

## Troubleshooting

### Common Issues

1. **Connection Refused**
   ```
   [ERROR] Cannot reach endpoint: http://localhost:8080/metrics
   ```
   - Verify the endpoint URL is correct
   - Ensure the service is running
   - Check firewall/network connectivity

2. **No Metrics Collected**
   ```
   [WARN] No metrics received from endpoint
   ```
   - Check if endpoint returns Prometheus format
   - Verify the `/metrics` path is correct
   - Ensure metrics are being exposed

3. **Permission Denied**
   ```
   bash: ./prometheus-to-csv.sh: Permission denied
   ```
   - Make the script executable: `chmod +x prometheus-to-csv.sh`

4. **Missing Dependencies**
   ```
   [ERROR] Missing required dependencies: curl
   ```
   - Install missing tools: `sudo apt-get install curl`

### Debug Mode

Enable verbose logging to troubleshoot issues:
```bash
./prometheus-to-csv.sh -e http://localhost:8080/metrics -v
```

## Integration Examples

### With Cron Jobs
```bash
# Add to crontab for hourly collection
0 * * * * /path/to/prometheus-to-csv.sh -e http://localhost:8080/metrics -d 3600 -o /data/metrics_$(date +\%Y\%m\%d_\%H).csv
```

### With Docker
```bash
# Run in Docker container
docker run --rm -v $(pwd):/data alpine/curl sh -c "
  apk add --no-cache bash gawk grep coreutils &&
  /data/prometheus-to-csv.sh -e http://host.docker.internal:8080/metrics -o /data/metrics.csv -d 300
"
```

### With Monitoring Scripts
```bash
#!/bin/bash
# Monitoring wrapper script
./prometheus-to-csv.sh -e http://localhost:8080/metrics -d 3600 -o hourly_metrics.csv
if [ $? -eq 0 ]; then
    echo "Metrics collection completed successfully"
    # Process CSV file, send alerts, etc.
else
    echo "Metrics collection failed" >&2
    exit 1
fi
```

## Contributing

Feel free to submit issues and enhancement requests!

## License

This script is provided as-is under the MIT License. 