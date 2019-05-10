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
	@install -Dm00755 bin/$(NAME) -t $(DESTDIR)/$(PREFIX)/bin

test:
	@go test ./...