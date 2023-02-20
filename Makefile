.SUFFIXES:
Makefile:;

.DEFAULT_GOAL := test

.PHONY: build
build:
	@echo build

.PHONY: clean
clean:
	@echo clean

.PHONY: test
test: build
	@echo test
