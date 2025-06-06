# Use Alpine Linux as base image for minimal size
FROM alpine:latest

# Install curl, bash, and GNU coreutils for proper date command
RUN apk add --no-cache curl bash coreutils

# Copy the entrypoint script
COPY entrypoint.sh ./entrypoint.sh

# Make the script executable
RUN chmod +x ./entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["./entrypoint.sh"]
