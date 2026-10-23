SHELL := /bin/bash
BIN_DIR := $(CURDIR)/bin
OAPI_CODEGEN := $(BIN_DIR)/oapi-codegen

.PHONY: all build run test lint-spec lint-redocly codegen clean verify tools

all: lint-spec codegen test build

tools: $(OAPI_CODEGEN)

$(OAPI_CODEGEN):
	@mkdir -p $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

codegen: $(OAPI_CODEGEN)
	@mkdir -p internal/api/v1
	$(OAPI_CODEGEN) -config oapi-codegen.yaml api/openapi.yaml
	@echo "==> Code generation complete: internal/api/v1/v1.gen.go"

lint-spec:
	@go run ./cmd/lint-spec api/openapi.yaml

lint-redocly:
	@npx --yes @redocly/cli lint api/openapi.yaml

build: codegen
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-s -w" -o $(BIN_DIR)/server ./cmd/server
	@echo "==> Binary built at $(BIN_DIR)/server"

run: build
	$(BIN_DIR)/server

test: codegen
	go test -v -race ./...

verify: build
	@chmod +x scripts/verify.sh
	./scripts/verify.sh

clean:
	rm -rf $(BIN_DIR) internal/api/v1/*.gen.go
