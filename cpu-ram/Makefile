# Application settings
APP_NAME := go-benchmark
PORT := 8080
IMAGE_NAME := $(APP_NAME)
HOST_URL := http://localhost:$(PORT)
LOCAL_REGISTRY := localhost:5000

# Version management
VERSION_FILE := VERSION
VERSION := $(shell cat $(VERSION_FILE))
IMAGE_TAG := $(IMAGE_NAME):$(VERSION)

# Default target
all: help

# Print help message
help:
	@echo "Go Benchmarking Container"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make build              Build the Docker image with auto-versioning"
	@echo "  make run                Run the Docker container"
	@echo "  make stop               Stop the Docker container"
	@echo "  make restart            Restart the Docker container"
	@echo "  make push-registry      Push the image to local registry at $(LOCAL_REGISTRY)"
	@echo ""
	@echo "Version Commands:"
	@echo "  make version            Display current version ($(VERSION))"
	@echo "  make bump-patch         Bump patch version (0.0.x)"
	@echo "  make bump-minor         Bump minor version (0.x.0)"
	@echo "  make bump-major         Bump major version (x.0.0)"
	@echo ""
	@echo "Docker Compose Commands:"
	@echo "  make up                 Start the container with Docker Compose"
	@echo "  make down               Stop the container with Docker Compose"
	@echo "  make logs               Show container logs"
	@echo ""
	@echo "Benchmark Commands:"
	@echo "  make status             Show the benchmark status"
	@echo "  make cpu-start          Start CPU benchmark using all cores"
	@echo "  make cpu-start-N        Start CPU benchmark using N cores (e.g., make cpu-start-2)"
	@echo "  make cpu-stop           Stop CPU benchmark"
	@echo "  make mem-start          Start memory benchmark with default limit (1GB)"
	@echo "  make mem-start-N        Start memory benchmark with N MB (e.g., make mem-start-512)"
	@echo "  make mem-stop           Stop memory benchmark (keeps memory allocated)"
	@echo "  make mem-free           Free allocated memory"
	@echo ""
	@echo "Complex Examples:"
	@echo "  make cpu-mem-test       Run both CPU and memory benchmarks"
	@echo "  make stop-clean         Stop all benchmarks and free memory"
	@echo "  make cpu-test-512       Test CPU benchmark with 512MB memory allocation"

# Version management
version:
	@echo "Current version: $(VERSION)"

bump-patch:
	@echo $$(awk -F. '{print $$1"."$$2"."$$3+1}' $(VERSION_FILE)) > $(VERSION_FILE)
	@echo "Bumped patch version to $$(cat $(VERSION_FILE))"

bump-minor:
	@echo $$(awk -F. '{print $$1"."$$2+1".0"}' $(VERSION_FILE)) > $(VERSION_FILE)
	@echo "Bumped minor version to $$(cat $(VERSION_FILE))"

bump-major:
	@echo $$(awk -F. '{print $$1+1".0.0"}' $(VERSION_FILE)) > $(VERSION_FILE)
	@echo "Bumped major version to $$(cat $(VERSION_FILE))"

# Docker commands
build: bump-patch
	@echo "Building image $(IMAGE_TAG)"
	docker build -t $(IMAGE_TAG) -t $(IMAGE_NAME):latest -f Dockerfile --build-arg VERSION=$$(cat $(VERSION_FILE)) .
	@echo "Built version $$(cat $(VERSION_FILE))"

run:
	docker run -d --name $(APP_NAME) -p $(PORT):8080 $(IMAGE_TAG)
	@echo "Server running at $(HOST_URL) (version $(VERSION))"

stop:
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

restart: stop run

# Docker Compose commands
up: docker-compose.yml
	VERSION=$$(cat $(VERSION_FILE)) docker-compose up -d
	@echo "Server running at $(HOST_URL) (version $$(cat $(VERSION_FILE)))"

down:
	docker-compose down

logs:
	docker-compose logs -f

# Development commands
dev:
	go run main.go

# Benchmark commands
status:
	@curl -s $(HOST_URL)/status
	@echo ""

cpu-start:
	@curl -s -X POST $(HOST_URL)/cpu/activate
	@echo ""

cpu-stop:
	@curl -s -X POST $(HOST_URL)/cpu/deactivate
	@echo ""

mem-start:
	@curl -s -X POST $(HOST_URL)/memory/activate
	@echo ""

mem-stop:
	@curl -s -X POST $(HOST_URL)/memory/deactivate
	@echo ""

mem-free:
	@curl -s -X POST $(HOST_URL)/memory/free
	@echo ""

# Dynamic targets for custom core count
cpu-start-%:
	@curl -s -X POST $(HOST_URL)/cpu/activate/$*
	@echo ""

# Dynamic targets for custom memory limit
mem-start-%:
	@curl -s -X POST $(HOST_URL)/memory/activate/$*
	@echo ""

# Complex targets
cpu-mem-test:
	@echo "Starting CPU benchmark..."
	@make cpu-start
	@echo "Starting memory benchmark..."
	@make mem-start
	@echo "Both benchmarks running. Check status with 'make status'"

stop-clean:
	@echo "Stopping CPU benchmark..."
	@make cpu-stop
	@echo "Stopping memory benchmark..."
	@make mem-stop
	@echo "Freeing memory..."
	@make mem-free
	@echo "All benchmarks stopped and memory freed"

# Custom testing scenarios
cpu-test-%:
	@echo "Running CPU benchmark with $* MB memory allocation"
	@make cpu-start
	@make mem-start-$*
	@echo "Test running. Check status with 'make status'"

high-load:
	@echo "Running high-load test with all CPU cores and 80% of memory limit..."
	@make cpu-start
	@make mem-start-800
	@echo "High load test running. Check status with 'make status'"

# Wait for the server to be up
wait-for-server:
	@echo "Waiting for server to be up..."
	@while ! curl -s $(HOST_URL)/health > /dev/null; do \
		sleep 1; \
	done
	@echo "Server is up!"

# Registry commands
tag-registry: build
	docker tag $(IMAGE_TAG) $(LOCAL_REGISTRY)/$(IMAGE_TAG)
	docker tag $(IMAGE_TAG) $(LOCAL_REGISTRY)/$(IMAGE_NAME):latest
	@echo "Tagged image for local registry"

push-registry: tag-registry
	docker push $(LOCAL_REGISTRY)/$(IMAGE_TAG)
	docker push $(LOCAL_REGISTRY)/$(IMAGE_NAME):latest
	@echo "Pushed $(IMAGE_TAG) to local registry at $(LOCAL_REGISTRY)"

# Automatically generated target to match any parameters
.PHONY: all help build run stop restart up down logs dev \
        status cpu-start cpu-stop mem-start mem-stop mem-free \
        cpu-mem-test stop-clean high-load wait-for-server \
        tag-registry push-registry version bump-patch bump-minor bump-major 