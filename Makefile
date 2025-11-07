SHELL := /bin/bash

BINARY := go-ad-admin
PKG := github.com/Weruminger/go-ad-admin

.PHONY: all test run build lint cover

all: test build

test:
	go test ./... -cover -covermode=atomic

build:
	CGO_ENABLED=0 go build -o bin/$(BINARY) ./cmd/go-ad-admin

run:
	go run ./cmd/go-ad-admin

cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic && go tool cover -func=coverage.out | tail -n1
