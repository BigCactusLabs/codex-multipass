BINARY_NAME=codex-mp
GO_DIR=go
GO_ENV=GOCACHE=$(CURDIR)/.gocache GOMODCACHE=$(CURDIR)/.gomodcache
GOLANGCI_LINT_VERSION=v2.12.1
SHELLCHECK=shellcheck

.PHONY: all build test clean tidy fmt lint shell-lint go-lint unit-test

all: build

tidy:
	cd $(GO_DIR) && $(GO_ENV) go mod tidy

build:
	cd $(GO_DIR) && $(GO_ENV) go build -ldflags "-X github.com/BigCactusLabs/codex-multipass/internal/app.Version=$(shell cat VERSION)" -o ../$(BINARY_NAME) cmd/codex-mp/main.go

test: build
	CODEX_MP=./$(BINARY_NAME) ./tests/smoke.sh
	CODEX_MP=./$(BINARY_NAME) ./tests/battle.sh
	CODEX_MP=./$(BINARY_NAME) ./tests/concurrency_test.sh
	CODEX_MP=./$(BINARY_NAME) ./tests/corrupt_storage_test.sh

clean:
	rm -f $(BINARY_NAME)

fmt:
	cd $(GO_DIR) && gofmt -w .

lint: shell-lint go-lint

shell-lint:
	$(SHELLCHECK) bash/codex-switch scripts/*.sh tests/*.sh

go-lint:
	cd $(GO_DIR) && $(GO_ENV) go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run

unit-test:
	cd $(GO_DIR) && $(GO_ENV) go test ./...
