.PHONY: build install

NAME := stack-updater
ifeq ($(PREFIX),)
    PREFIX := $(HOME)/.local
endif


build:
	@mkdir -p bin
	@go build -o $(CURDIR)/bin/$(NAME) ./main.go 

install:
	@install -D bin/* -t $(PREFIX)/bin

test:
	go test ./...