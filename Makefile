NAME := omssh
DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
VERSION := $(gobump show -r)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := "-X main.revision=$(REVISION)"

export GO111MODULE=on

## Install dependencies

.PHONY: deps
deps:
	go get -v -d

# 開発に必要な依存をインストールする ## Setup
.PHONY: deps
devel-deps: deps
	GO111MODULE=off go get \
		github.com/golang/lint/golint \
		github.com/motemen/gobump/cmd/gobump \
		github.com/Songmu/make2help/cmd/make2help

# テストを実行する
## Run tests
.PHONY: test
test: deps
	go test ./...

## Lint
.PHONY: lint
lint: devel-deps
	go vet ./...
	golint -set_exit_status ./...

## build binaries ex. make bin/omssh
bin/%: *.go cmd/%/main.go deps
	cd cmd/omssh && go build -ldflags "-s -w -X main.version=${GIT_VER} -X main.buildDate=${DATE}" -gcflags="-trimpath=${PWD}"

install: build
	install cmd/omssh/omssh ${GOPATH}/bin

## build binary
.PHONY: build
build: bin/omssh

## Show help
.PHONY: help
help:
	@make2help $(MAKEFILE_LIST)

clean:
	rm -f cmd/omssh/omssh
	rm -f pkg/*
