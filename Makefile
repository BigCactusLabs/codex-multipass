BINARY_NAME=codex-mp
GO_DIR=go

.PHONY: all build test clean tidy

all: build

tidy:
	cd $(GO_DIR) && go mod tidy

build:
	cd $(GO_DIR) && go build -ldflags "-X github.com/BigCactusLabs/codex-multipass/internal/app.Version=$(shell cat VERSION)" -o ../$(BINARY_NAME) cmd/codex-mp/main.go

test: build
	CODEX_MP=./$(BINARY_NAME) ./tests/smoke.sh
	CODEX_MP=./$(BINARY_NAME) ./tests/battle.sh

clean:
	rm -f $(BINARY_NAME)
	rm -f $(GO_DIR)/go.sum
