# Application settings
APP_NAME := fps-benchmarking
IMAGE_NAME := $(APP_NAME)
CONTAINER_NAME := $(APP_NAME)
SOCKET_PATH := /metrics
LOCAL_REGISTRY := localhost:5000
#LOCAL_REGISTRY := 131.159.25.135:5000

# Version management
VERSION_FILE := VERSION
VERSION := $(shell cat $(VERSION_FILE))
IMAGE_TAG := $(IMAGE_NAME):$(VERSION)

# Default target
all: help

# Print help message
help:
	@echo "FPS Benchmarking Container"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make build              Build the Docker image with auto-versioning"
	@echo "  make run                Run the Docker container"
	@echo "  make stop               Stop the Docker container"
	@echo "  make restart            Restart the Docker container"
	@echo "  make logs               Show container logs"
	@echo "  make push-registry      Push the image to local registry at $(LOCAL_REGISTRY)"
	@echo ""
	@echo "Version Commands:"
	@echo "  make version            Display current version ($(VERSION))"
	@echo "  make bump-patch         Bump patch version (0.0.x)"
	@echo "  make bump-minor         Bump minor version (0.x.0)"
	@echo "  make bump-major         Bump major version (x.0.0)"
	@echo ""
	@echo "Development Commands:"
	@echo "  make build-local        Build the application locally"
	@echo "  make run-local          Run the application locally"
	@echo ""

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
	docker run -d --name $(CONTAINER_NAME) \
		-v $(SOCKET_PATH):$(SOCKET_PATH) \
		$(IMAGE_TAG)
	@echo "FPS benchmark container running (version $(VERSION))"

stop:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

restart: stop run

logs:
	docker logs -f $(CONTAINER_NAME)

# Build and push to local registry
tag-registry: build
	docker tag $(IMAGE_TAG) $(LOCAL_REGISTRY)/$(IMAGE_TAG)
	docker tag $(IMAGE_TAG) $(LOCAL_REGISTRY)/$(IMAGE_NAME):latest
	@echo "Tagged image for local registry"

push-registry: tag-registry
	docker push $(LOCAL_REGISTRY)/$(IMAGE_TAG)
	docker push $(LOCAL_REGISTRY)/$(IMAGE_NAME):latest
	@echo "Pushed $(IMAGE_TAG) to local registry at $(LOCAL_REGISTRY)"

# Development commands
build-local: bump-patch
	go build -o fps-benchmark -ldflags "-X main.buildVersion=$$(cat $(VERSION_FILE))"

run-local:
	./fps-benchmark

# Create Docker Compose file
docker-compose.yml:
	@echo "version: '3'" > docker-compose.yml
	@echo "" >> docker-compose.yml
	@echo "services:" >> docker-compose.yml
	@echo "  fps-benchmark:" >> docker-compose.yml
	@echo "    build: ." >> docker-compose.yml
	@echo "    container_name: $(CONTAINER_NAME)" >> docker-compose.yml
	@echo "    volumes:" >> docker-compose.yml
	@echo "      - $(SOCKET_PATH):$(SOCKET_PATH)" >> docker-compose.yml
	@echo "Docker Compose file created: docker-compose.yml"

# Docker Compose commands
up: docker-compose.yml build
	VERSION=$$(cat $(VERSION_FILE)) docker-compose up -d
	@echo "FPS benchmark container running with Docker Compose (version $$(cat $(VERSION_FILE)))"

down:
	docker-compose down

# Clean up
clean:
	rm -f fps-benchmark

.PHONY: all help build run stop restart logs tag-registry push-registry build-local run-local docker-compose.yml up down clean version bump-patch bump-minor bump-major 
