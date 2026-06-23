SHELL := /bin/sh

APP_DIR := apps/api
IMAGE ?= release-sentinel-api:local
PORT ?= 8080
APP_VERSION ?= local
COMMIT_SHA ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)

.PHONY: test
test:
	cd $(APP_DIR) && go test ./...

.PHONY: vet
vet:
	cd $(APP_DIR) && go vet ./...

.PHONY: run
run:
	cd $(APP_DIR) && APP_VERSION=$(APP_VERSION) COMMIT_SHA=$(COMMIT_SHA) HTTP_ADDR=:$(PORT) go run ./cmd/server

.PHONY: build
build:
	docker build \
		--build-arg APP_VERSION=$(APP_VERSION) \
		--build-arg COMMIT_SHA=$(COMMIT_SHA) \
		-t $(IMAGE) \
		-f $(APP_DIR)/Dockerfile .

.PHONY: compose-up
compose-up:
	docker compose up --build

.PHONY: compose-down
compose-down:
	docker compose down --remove-orphans

.PHONY: helm-template
helm-template:
	helm template release-sentinel deploy/helm/release-sentinel --values deploy/helm/release-sentinel/values.yaml

.PHONY: validate
validate: test vet
	./scripts/validate-yaml.sh

.PHONY: load-test
load-test:
	k6 run tests/load/release-validation.js

.PHONY: smoke
smoke:
	./scripts/smoke.sh
