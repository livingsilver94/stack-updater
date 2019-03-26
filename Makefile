.PHONY: build install

NAME := stack-updater
ifeq ($(PREFIX),)
    PREFIX := $(HOME)/.local
endif


build:
	@mkdir -p bin
	@GOBIN=$(CURDIR)/bin go install cmd/*

install:
	@install -D $(CURDIR)/bin/$(NAME) -t $(PREFIX)/bin
