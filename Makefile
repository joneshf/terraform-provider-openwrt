.SUFFIXES:
Makefile:;

.DEFAULT_GOAL := test

.PHONY: build
build:
	go build ./...

.PHONY: clean
clean:
	@echo clean

.PHONY: docs
docs:
	go generate ./...

.PHONY: test
test: build test-docs test-go

.PHONY: test-docs
test-docs: test-docs-up-to-date

.PHONY: test-docs-up-to-date
test-docs-up-to-date:
	./scripts/test-docs-up-to-date.sh

.PHONY: test-go
test-go:
	go test ./...
