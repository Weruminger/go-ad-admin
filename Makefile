SHELL := /bin/bash
PKG := github.com/Weruminger/go-ad-admin
BINARY := go-ad-admin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD   ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD)
.PHONY: all test run build lint cover

all: test build

build:
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o bin/$(BINARY) ./cmd/go-ad-admin

test:
	go test ./... -cover -v

run: build
	./bin/$(BINARY)


cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic && go tool cover -func=coverage.out | tail -n1
