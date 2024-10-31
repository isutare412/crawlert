ifneq (,$(wildcard ./.env.local))
    include .env.local
    export $(shell sed 's/=.*//' ./.env.local)
endif

TAG ?= latest
IMAGE ?= redshoore/crawlert:$(TAG)

DOCKER_USER ?= <docker_hub_username>
DOCKER_PASSWORD ?= <docker_hub_secret>

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: run
run: ## Run Crawlert.
	go run ./cmd/...

.PHONY: mocks
mocks: mockery ## Generate mock implementations.
	$(MOCKERY) --config mockery.yaml

##@ Build

.PHONY: image-build
image-build: ## Build docker image.
	docker build -f ./Dockerfile -t $(IMAGE) .

.PHONY: image-push
image-push: ## Push docker image.
	echo $(DOCKER_PASSWORD) | docker login -u $(DOCKER_USER) --password-stdin
	docker push $(IMAGE)

##@ Tool Dependencies

## Location to install dependencies to.
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

## Tool Binaries
MOCKERY ?= $(LOCALBIN)/mockery

.PHONY: mockery
mockery: $(MOCKERY) ## Install mockery if necessary.
$(MOCKERY): $(LOCALBIN)
	@test -s $(MOCKERY) || \
	GOBIN=$(LOCALBIN) go install github.com/vektra/mockery/v2@latest
