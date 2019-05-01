.PHONY: generate build install test

NAME := stack-updater
ifeq ($(PREFIX),)
    PREFIX := $(HOME)/.local
endif

.DEFAULT_GOAL := build

generate:
	@go generate ./...

build:
	@mkdir -p bin
	@go build -o $(CURDIR)/bin/$(NAME) ./main.go 

install:
	@install -D bin/* -t $(PREFIX)/bin

test:
	@go test ./...