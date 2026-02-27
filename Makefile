BINARY := unifai
BUILD_DIR := bin
GOCACHE ?= /tmp/go-build
GOMODCACHE ?= /tmp/go-mod

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
	-X unifai/internal/version.Version=$(VERSION) \
	-X unifai/internal/version.Commit=$(COMMIT) \
	-X unifai/internal/version.BuildDate=$(BUILD_DATE)

.PHONY: tidy fmt test build run snapshot-release release

tidy:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go mod tidy

fmt:
	gofmt -w cmd internal

test:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go test ./...

build:
	mkdir -p $(BUILD_DIR)
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) ./cmd/unifai

run:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go run ./cmd/unifai

snapshot-release:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean
