#!/bin/bash

# Example usage of prometheus-to-csv.sh script
# This script demonstrates various ways to use the Prometheus to CSV converter

echo "=== Prometheus to CSV Converter - Example Usage ==="
echo

# Check if the main script exists
if [[ ! -f "prometheus-to-csv.sh" ]]; then
    echo "Error: prometheus-to-csv.sh not found in current directory"
    exit 1
fi

echo "1. Show help:"
echo "./prometheus-to-csv.sh --help"
echo
./prometheus-to-csv.sh --help
echo

echo "2. Basic usage (replace with your actual endpoint):"
echo "./prometheus-to-csv.sh -e http://localhost:8080/metrics -v"
echo

echo "3. Collect metrics for 60 seconds with 5-second intervals:"
echo "./prometheus-to-csv.sh -e http://localhost:8080/metrics -i 5 -d 60 -o sample_metrics.csv -v"
echo

echo "4. Filter specific metrics (CPU and memory related):"
echo "./prometheus-to-csv.sh -e http://localhost:8080/metrics -f \"cpu|memory|ram\" -v"
echo

echo "5. Show available metrics from endpoint:"
echo "./prometheus-to-csv.sh -e http://localhost:8080/metrics --help-metrics"
echo

echo "6. Collect without timestamps:"
echo "./prometheus-to-csv.sh -e http://localhost:8080/metrics -t -o no_timestamp_metrics.csv"
echo

echo "=== Example with mock Prometheus endpoint ==="
echo

# Create a simple mock Prometheus endpoint for testing
cat << 'EOF' > mock-prometheus-server.py
#!/usr/bin/env python3
"""
Simple mock Prometheus metrics server for testing
"""
import http.server
import socketserver
import time
import random
import math

class PrometheusHandler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/metrics':
            self.send_response(200)
            self.send_header('Content-type', 'text/plain; version=0.0.4; charset=utf-8')
            self.end_headers()
            
            # Generate mock metrics
            timestamp = time.time()
            cpu_usage = 50 + 30 * math.sin(timestamp / 10)  # Oscillating CPU usage
            memory_usage = 1024 + 512 * random.random()     # Random memory usage
            fps_value = 60 + 10 * math.sin(timestamp / 5)   # Oscillating FPS
            
            metrics = f"""# HELP cpu_usage_percent CPU usage percentage
# TYPE cpu_usage_percent gauge
cpu_usage_percent{{instance="localhost",job="test"}} {cpu_usage:.2f}

# HELP memory_usage_bytes Memory usage in bytes
# TYPE memory_usage_bytes gauge
memory_usage_bytes{{instance="localhost",job="test"}} {memory_usage:.0f}

# HELP fps_current Current frames per second
# TYPE fps_current gauge
fps_current{{instance="localhost",job="test"}} {fps_value:.2f}

# HELP http_requests_total Total HTTP requests
# TYPE http_requests_total counter
http_requests_total{{method="GET",status="200"}} {int(timestamp) % 1000}
http_requests_total{{method="POST",status="200"}} {int(timestamp) % 500}
http_requests_total{{method="GET",status="404"}} {int(timestamp) % 100}

# HELP process_start_time_seconds Start time of the process since unix epoch
# TYPE process_start_time_seconds gauge
process_start_time_seconds {timestamp - 3600}
"""
            self.wfile.write(metrics.encode())
        else:
            self.send_response(404)
            self.end_headers()
    
    def log_message(self, format, *args):
        # Suppress default logging
        pass

PORT = 8080

print(f"Starting mock Prometheus server on port {PORT}")
print("Metrics available at: http://localhost:8080/metrics")
print("Press Ctrl+C to stop")

try:
    with socketserver.TCPServer(("", PORT), PrometheusHandler) as httpd:
        httpd.serve_forever()
except KeyboardInterrupt:
    print("\nServer stopped")
EOF

chmod +x mock-prometheus-server.py

echo "Created mock-prometheus-server.py for testing"
echo
echo "To test the script:"
echo "1. Start the mock server: python3 mock-prometheus-server.py"
echo "2. In another terminal, run: ./prometheus-to-csv.sh -e http://localhost:8080/metrics -v -d 30"
echo
echo "The mock server provides sample metrics including:"
echo "- cpu_usage_percent (oscillating values)"
echo "- memory_usage_bytes (random values)"
echo "- fps_current (oscillating FPS values)"
echo "- http_requests_total (counter metrics)"
echo "- process_start_time_seconds (gauge metric)" 