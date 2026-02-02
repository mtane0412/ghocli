# Makefile for gho

# Build information
VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Default target
.PHONY: all
all: build

# Build
.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o gho ./cmd/gho

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run lint (requires golangci-lint)
.PHONY: lint
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Type check
.PHONY: type-check
type-check:
	go vet ./...

# Cleanup
.PHONY: clean
clean:
	rm -f gho
	rm -f coverage.out coverage.html

# Install
.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)" ./cmd/gho

# Help
.PHONY: help
help:
	@echo "gho Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  build          - Build the gho binary"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  lint           - Run golangci-lint"
	@echo "  type-check     - Run go vet"
	@echo "  clean          - Remove build artifacts"
	@echo "  install        - Install gho to GOPATH/bin"
	@echo "  help           - Show this help message"
