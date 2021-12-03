.DEFAULT_GOAL := help

COMPOSE := docker-compose
COMPOSE_FILES := -f docker-compose.yml
DOCKER_COMPOSE := $(COMPOSE) $(COMPOSE_FILES)
CGO_ENABLED=0
GO_BUILD_FLAGS=-ldflags "-extldflags '-static'"
GO_BUILD := CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS)

########################### сборка образа с docker-compose

.PHONY: run
run: ## docker-compose up
	$(DOCKER_COMPOSE) up

########################### документация makefile

.PHONY: help
help: ## show this help
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

########################### используется в docker файлах для сборки проектов

.PHONY: build-bot
build-bot: ## compile telegram bot
	$(GO_BUILD) -o bin/app ./cmd/app

.PHONY: build-downloader
build-downloader: ## compile downloader
	$(GO_BUILD) -o bin/download ./cmd/download
