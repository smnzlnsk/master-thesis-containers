# Files Created - Prometheus to CSV Converter

This document lists all the files created for the Prometheus to CSV converter project.

## Main Script
- **`prometheus-to-csv.sh`** - The main bash script that collects Prometheus metrics and converts them to CSV format
  - Comprehensive command-line interface with multiple options
  - Robust error handling and validation
  - Support for filtering, timestamps, and custom intervals
  - Verbose logging with color-coded output

## Documentation
- **`PROMETHEUS_CSV_README.md`** - Comprehensive documentation including:
  - Feature overview and requirements
  - Installation instructions
  - Usage examples and command-line options
  - Output format specifications
  - Troubleshooting guide
  - Integration examples

## Testing and Examples
- **`example-usage.sh`** - Example script demonstrating various usage patterns
  - Shows different command-line options
  - Creates a mock Prometheus server for testing
  - Provides practical examples

- **`mock-prometheus-server.py`** - Python-based mock Prometheus server for testing
  - Generates realistic sample metrics (CPU, memory, FPS, HTTP requests)
  - Oscillating and random values for testing
  - Standard Prometheus format output

## Summary File
- **`FILES_CREATED.md`** - This file, documenting all created components

## Key Features of the Main Script

### Command-Line Options
```bash
-e, --endpoint URL        # Required: Prometheus endpoint
-o, --output FILE         # Output CSV file
-i, --interval SECONDS    # Collection interval
-d, --duration SECONDS    # Total duration
-f, --filter PATTERN     # Filter metrics by regex
-t, --no-timestamp        # Exclude timestamps
-m, --help-metrics        # Show available metrics
-v, --verbose             # Verbose output
-h, --help                # Show help
```

### Output Format
```csv
timestamp,epoch,metric_name,labels,value
2024-01-15 10:30:00,1705312200,"cpu_usage_percent","instance=""localhost""",45.67
```

### Error Handling
- Connection validation
- Dependency checking
- Timeout handling
- Signal handling (Ctrl+C)
- Data validation

## Usage Examples

### Basic Collection
```bash
./prometheus-to-csv.sh -e http://localhost:8080/metrics -v
```

### Filtered Collection
```bash
./prometheus-to-csv.sh -e http://localhost:8080/metrics -f "cpu|memory" -d 300
```

### Testing with Mock Server
```bash
# Terminal 1
python3 mock-prometheus-server.py

# Terminal 2
./prometheus-to-csv.sh -e http://localhost:8080/metrics -v -d 30
```

All files are ready to use and have been tested for functionality. 