# import deploy config
dpl ?= deploy.env
include $(dpl)
export $(shell sed 's/=.*//' $(dpl))

VERSION?=v1

PACKAGES = $(shell go list ./... | grep -v 'vendor\|pkg/generated')
VERIFY_DOCKER_DEP=$(shell which docker >> /dev/null 2>&1; echo $$?)


.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


.PHONY: all
all: format test deps release

.PHONY: format
format: ## Running go fmt
	@echo "--------------------------"
	@echo "--> Running go fmt"
	@go fmt $(PACKAGES)

.PHONY: deps
deps: verify-deps clean-deps ## Dpendency verification and clean up

clean-deps: ## Cleans up unused dependencies or adds missing dependencies.
	@echo "--------------------------"
	@echo "--> Tidying up submodules"
	@go mod tidy
	@echo "--> Verifying submodules"
	@go mod verify

verify-deps: ## Running dependency check for docker
	@echo "--------------------------"
	@echo "--> Verifying docker dependency"
	@if [ $(VERIFY_DOCKER_DEP) -ne 0 ]; \
		then echo "ERROR:\t Docker is not installed. Please install docker to build the image" \
		&& exit 1; \
	fi;
	@echo "Docker dependency verified"

test: ## Run unit tests
	@echo "--------------------------"
	@echo "--> Verifying docker dependency"
	@go test -v -cover ./...


# DOCKER TASKS
# Build the container
.PHONY: build
build: ## Build the container
	@echo "--------------------------"
	@echo "--> Building the container"
	docker build -t $(EXPORTER_NAME) .

build-nc: ## Build the container without caching
	@echo "--------------------------"
	@echo "--> Building the container without caching"
	docker build --no-cache -t $(EXPORTER_NAME) .

.PHONY: release
release: build-nc push ## Make a release by building and pushing the `{version}` ans `latest` tagged image to docker

.PHONY: push
push: publish-latest publish-version ## Push the `{version}` ans `latest` tagged iamge to docker

publish-latest: tag-latest ## Publish the `latest` taged container to docker
	@echo "--------------------------"
	@echo "--> Publishing the `latest` taged container to docker"
	@echo 'publish latest to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(EXPORTER_NAME):latest

publish-version: tag-version ## Publish the `{version}` taged container to docker
	@echo "--------------------------"
	@echo "--> Publishing the `$(version)` taged container to docker"
	@echo 'publish $(VERSION) to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(EXPORTER_NAME):$(VERSION)

# Docker tagging
.PHONY: tag
tag: tag-latest tag-version ## Generate container tags for the `{version}` ans `latest` tags

tag-latest: ## Generate container `{version}` tag
	@echo "--------------------------"
	@echo 'create tag latest'
	@echo "--> Tagging container with latest"
	docker tag $(EXPORTER_NAME) $(DOCKER_REPO)/$(EXPORTER_NAME):latest

tag-version: ## Generate container `latest` tag
	@echo "--------------------------"
	@echo 'create tag $(VERSION)'
	@echo "--> Tagging container with $(VERSION)"
	docker tag $(EXPORTER_NAME) $(DOCKER_REPO)/$(EXPORTER_NAME):$(VERSION)

