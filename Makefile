# Makefile for gho

# ビルド情報
VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# デフォルトターゲット
.PHONY: all
all: build

# ビルド
.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o gho ./cmd/gho

# テスト実行
.PHONY: test
test:
	go test -v ./...

# カバレッジ付きテスト
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint実行（golangci-lintが必要）
.PHONY: lint
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# 型チェック
.PHONY: type-check
type-check:
	go vet ./...

# クリーンアップ
.PHONY: clean
clean:
	rm -f gho
	rm -f coverage.out coverage.html

# インストール
.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)" ./cmd/gho

# ヘルプ
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
