# Application settings
APP_NAME := curl-client
IMAGE_NAME := $(APP_NAME)
LOCAL_REGISTRY := localhost:5000

# Version management
VERSION_FILE := VERSION
VERSION := $(shell cat $(VERSION_FILE))
IMAGE_TAG := $(IMAGE_NAME):$(VERSION)

# Default target
all: help

# Print help message
help:
	@echo "Curl Client"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make build              Build the Docker image with auto-versioning"
	@echo "  make push-registry      Push the image to local registry at $(LOCAL_REGISTRY)"
	@echo ""
	@echo "Version Commands:"
	@echo "  make version            Display current version ($(VERSION))"
	@echo "  make bump-patch         Bump patch version (0.0.x)"
	@echo "  make bump-minor         Bump minor version (0.x.0)"
	@echo "  make bump-major         Bump major version (x.0.0)"
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
	docker run -it $(IMAGE_TAG)

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
.PHONY: all help build run tag-registry push-registry version bump-patch bump-minor bump-major 