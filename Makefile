.SUFFIXES:
Makefile:;

ACCEPTANCE_TEST_BUILD_CONSTRAINT := acceptance.test
ACCEPTANCE_TEST_DOCKER_COMPOSE_FILE := ./docker-compose.acceptance-test.yaml
CACHE_DIRECTORY := .cache

.DEFAULT_GOAL := test

.PHONY: build
build:
	go build ./...

.PHONY: bump-patch-version
bump-patch-version:
	./scripts/bump-patch-version.sh

.PHONY: clean
clean: clean-acceptance-test-server clean-cache-directory

.PHONY: clean-acceptance-test-server
clean-acceptance-test-server:
	docker compose --file $(ACCEPTANCE_TEST_DOCKER_COMPOSE_FILE) down --remove-orphans --rmi all --volumes

.PHONY: clean-cache-directory
clean-cache-directory:
	rm -fr $(CACHE_DIRECTORY)

.PHONY: docs
docs: install
	go generate ./...

.PHONY: install
install:
	go install ./...

.PHONY: release
release:
	goreleaser release --rm-dist

.PHONY: start-acceptance-test-server
start-acceptance-test-server:
	docker compose --file $(ACCEPTANCE_TEST_DOCKER_COMPOSE_FILE) up --build --remove-orphans --wait

.PHONY: test
test: build test-docs test-go

.PHONY: test-docs
test-docs: test-docs-up-to-date

.PHONY: test-docs-up-to-date
test-docs-up-to-date:
	./scripts/test-docs-up-to-date.sh

.PHONY: test-go
test-go: test-go-unit-test test-go-acceptance-test

.PHONY: test-go-acceptance-test
test-go-acceptance-test:
	TF_ACC=1 go test -tags=$(ACCEPTANCE_TEST_BUILD_CONSTRAINT) ./...

.PHONY: test-go-unit-test
test-go-unit-test:
	go test ./...
