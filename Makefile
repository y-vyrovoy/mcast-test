export GO111MODULE=on

SHELL=/bin/bash
#IMAGE_TAG := $(shell git rev-parse HEAD)
#DOCKER_COMPOSE = IMAGE_TAG=$(IMAGE_TAG) docker-compose -f docker-composition/default.yml -f docker-composition/system-test-mask.yml
#DOCKER_REPO = nexus.tools.devopenocean.studio

.PHONY: all
all: deps deps_check lint test build

.PHONY: ci
ci: lint test

.PHONY: deps
deps:
	go mod download
	go mod vendor

.PHONY: deps_check
deps_check:
	go mod verify

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -o artifacts/svc .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -mod=vendor -count=1 -v -cover ./...

.PHONY: dockerise
dockerise:
	docker build -t "${DOCKER_REPO}/frontend-app-renderer:${IMAGE_TAG}" .