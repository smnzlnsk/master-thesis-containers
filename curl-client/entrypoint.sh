#!/bin/bash

# Default endpoint if not provided
DEFAULT_URL="http://10.30.1.2/"

# Get the URL from environment variable or command line argument
URL="${TARGET_URL:-${1:-$DEFAULT_URL}}"

# Validate URL
if [[ ! "$URL" =~ ^https?:// ]]; then
    echo "Error: Invalid URL format. URL must start with http:// or https://"
    echo "Usage: docker run <image> [URL]"
    echo "   or: docker run -e TARGET_URL=<url> <image>"
    exit 1
fi

echo "Starting curl client for endpoint: $URL"
echo "Curling every 50ms..."
echo "Press Ctrl+C to stop"

# Counter for requests
counter=0

# Extract host from URL for logging
HOST=$(echo "$URL" | sed -E 's|^https?://([^/]+).*|\1|')

# CSV log file path
LOG_FILE="/metrics/curl_metrics.csv"

# Create CSV header if file doesn't exist
if [ ! -f "$LOG_FILE" ]; then
    echo "timestamp,host,return_code,return_body" > "$LOG_FILE"
    echo "Created CSV log file: $LOG_FILE"
fi

# Function to handle cleanup on exit
cleanup() {
    echo ""
    echo "Stopping curl client..."
    echo "Total requests sent: $counter"
    exit 0
}

# Set up signal handler for graceful shutdown
trap cleanup SIGINT SIGTERM

# Main loop - curl every 50ms (0.05 seconds)
while true; do
    counter=$((counter + 1))
    
    # Get current timestamp in microseconds from epoch
    timestamp=$(date +%s%6N)
    
    # Capture both body and status code in one variable
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" -m 5 "$URL" 2>/dev/null)
    curl_exit_code=$?
    
    if [ $curl_exit_code -eq 0 ]; then
        # Extract status code and response body
        http_code=$(echo "$response" | grep -o "HTTPSTATUS:.*" | cut -d: -f2)
        response_body=$(echo "$response" | sed 's/HTTPSTATUS:.*$//' | tr -d '\n\r' | sed 's/,/;/g')  # Replace commas to avoid CSV issues
        
        echo "[$counter] Status: $http_code | Body: $response_body"
        
        # Log to CSV: timestamp,host,return_code,return_body
        echo "$timestamp,$HOST,$http_code,\"$response_body\"" >> "$LOG_FILE"
    else
        echo "[$counter] Request failed - connection timeout or error"
        
        # Log failed request to CSV
        echo "$timestamp,$HOST,ERROR,\"Request failed\"" >> "$LOG_FILE"
    fi
    
    # Sleep for 50ms (0.05 seconds)
    sleep 0.05
done 