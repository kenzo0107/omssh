NAME := omssh
VERSION := $(shell gobump show -r cmd/omssh)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := "-X main.version=$(VERSION) -X main.revision=$(REVISION)"

export GO111MODULE=on

## Install dependencies
.PHONY: deps
deps:
	go get -v -d

## Setup
.PHONY: devel-deps
devel-deps: deps
	GO111MODULE=off go get \
		golang.org/x/lint/golint \
		honnef.co/go/tools/staticcheck \
		github.com/kisielk/errcheck \
		golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow \
		github.com/securego/gosec/cmd/gosec \
		github.com/motemen/gobump/cmd/gobump \
		github.com/Songmu/make2help/cmd/make2help

## Run tests
.PHONY: test
test: deps
	go test -v -count=1 -cover ./...

.PHONY: cov
cov:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out

## Lint
.PHONY: lint
lint: devel-deps
	go vet ./...
	go vet -vettool=$(shell which shadow)
	staticcheck ./...
	errcheck ./...
	# exclude G106: Audit the use of ssh.InsecureIgnoreHostKey
	gosec -quiet -exclude=G106 ./... 
	golint -set_exit_status ./...

## build binaries ex. make bin/omssh
bin/%: cmd/%/main.go deps
	@go build -ldflags $(LDFLAGS) -o $@ $<

## build binary
.PHONY: build
build: bin/${NAME}

.PHONY: install
install:
	@go install ./cmd/omssh

## Show help
.PHONY: help
help:
	@make2help $(MAKEFILE_LIST)

## clean package
.PHONY: clean
clean:
	rm -f bin/${NAME}

.PHONY: release
release:
	@git tag v$(VERSION)
	@git push --tags
	goreleaser --rm-dist

