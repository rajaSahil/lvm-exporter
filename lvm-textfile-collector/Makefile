# import deploy config
dpl ?= deploy.env
include $(dpl)
export $(shell sed 's/=.*//' $(dpl))

VERSION?=v1

VERIFY_DOCKER_DEP=$(shell which docker >> /dev/null 2>&1; echo $$?)

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: all
all: deps release

.PHONY: deps
deps: verify-deps ## Dpendency verification and clean up

verify-deps: ## Running dependency check for docker
	@echo "--------------------------"
	@echo "--> Verifying docker dependency"
	@if [ $(VERIFY_DOCKER_DEP) -ne 0 ]; \
		then echo "ERROR:\t Docker is not installed. Please install docker to build the image" \
		&& exit 1; \
	fi;
	@echo "Docker dependency verified"

# DOCKER TASKS
# Build the container
.PHONY: build
build: ## Build the container
	@echo "--------------------------"
	@echo "--> Building the container"
	docker build -t $(LVM_TEXTFILE_COLLECTOR) .

build-nc: ## Build the container without caching
	@echo "--------------------------"
	@echo "--> Building the container without caching"
	docker build --no-cache -t $(LVM_TEXTFILE_COLLECTOR) .

.PHONY: release
release: build-nc push ## Make a release by building and pushing the `{version}` ans `latest` tagged image to docker

.PHONY: push
push: publish-latest publish-version ## Push the `{version}` ans `latest` tagged iamge to docker

publish-latest: tag-latest ## Publish the `latest` taged container to docker
	@echo "--------------------------"
	@echo "Publishing the `latest` taged container to docker"
	@echo 'publish latest to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(LVM_TEXTFILE_COLLECTOR):latest

publish-version: tag-version ## Publish the `{version}` taged container to docker
	@echo "--------------------------"
	@echo "Publishing the `$(version)` taged container to docker"
	@echo 'publish $(VERSION) to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(LVM_TEXTFILE_COLLECTOR):$(VERSION)

# Docker tagging
.PHONY: tag
tag: tag-latest tag-version ## Generate container tags for the `{version}` ans `latest` tags

tag-latest: ## Generate container `{version}` tag
	@echo "--------------------------"
	@echo 'create tag latest'
	@echo "--> Tagging container with latest"
	docker tag $(LVM_TEXTFILE_COLLECTOR) $(DOCKER_REPO)/$(LVM_TEXTFILE_COLLECTOR):latest

tag-version: ## Generate container `latest` tag
	@echo "--------------------------"
	@echo 'create tag $(VERSION)'
	@echo "--> Tagging container with $(VERSION)"
	docker tag $(LVM_TEXTFILE_COLLECTOR) $(DOCKER_REPO)/$(LVM_TEXTFILE_COLLECTOR):$(VERSION)

