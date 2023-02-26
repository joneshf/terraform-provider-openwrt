.SUFFIXES:
Makefile:;

ACCEPTANCE_TEST_BUILD_CONSTRAINT := acceptance.test
ACCEPTANCE_TEST_DOCKER_COMPOSE_FILE := lucirpc/docker-compose.acceptance-test.yaml

.DEFAULT_GOAL := test

.PHONY: build
build:
	go build ./...

.PHONY: clean
clean: clean-acceptance-test-server

.PHONY: clean-acceptance-test-server
clean-acceptance-test-server:
	docker compose --file $(ACCEPTANCE_TEST_DOCKER_COMPOSE_FILE) down --remove-orphans --rmi all --volumes

.PHONY: docs
docs:
	go generate ./...

.PHONY: install
install:
	go install ./...

.PHONY: start-acceptance-test-server
start-acceptance-test-server:
	docker compose --file $(ACCEPTANCE_TEST_DOCKER_COMPOSE_FILE) up --build --detach --remove-orphans

.PHONY: test
test: build start-acceptance-test-server test-docs test-go

.PHONY: test-docs
test-docs: test-docs-up-to-date

.PHONY: test-docs-up-to-date
test-docs-up-to-date:
	./scripts/test-docs-up-to-date.sh

.PHONY: test-go
test-go: test-go-unit-test test-go-acceptance-test

.PHONY: test-go-acceptance-test
test-go-acceptance-test:
	go test -tags=$(ACCEPTANCE_TEST_BUILD_CONSTRAINT) ./...

.PHONY: test-go-unit-test
test-go-unit-test:
	go test ./...
