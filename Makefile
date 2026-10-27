SHELL := /bin/bash
BIN_DIR := $(CURDIR)/bin
DIST_DIR := $(CURDIR)/dist
OAPI_CODEGEN := $(BIN_DIR)/oapi-codegen
BUNDLED_SPEC := $(DIST_DIR)/openapi.bundled.yaml

.PHONY: all build run test lint-spec lint-redocly bundle-spec codegen clean verify tools payment-spec

all: lint-spec lint-redocly codegen test build

tools: $(OAPI_CODEGEN)

$(OAPI_CODEGEN):
	@mkdir -p $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

payment-spec:
	@mkdir -p $(DIST_DIR)
	@go run ./cmd/payment-spec > $(DIST_DIR)/openapi.payments.yaml

bundle-spec: payment-spec
	@mkdir -p $(DIST_DIR)
	@npx --yes @redocly/cli@2.60.0 bundle $(DIST_DIR)/openapi.payments.yaml -o $(BUNDLED_SPEC)
	@echo "==> Bundled OpenAPI spec at $(BUNDLED_SPEC)"

codegen: $(OAPI_CODEGEN) bundle-spec
	@mkdir -p internal/api/v1
	$(OAPI_CODEGEN) -config oapi-codegen.yaml $(BUNDLED_SPEC)
	@echo "==> Code generation complete: internal/api/v1/v1.gen.go"

lint-spec: payment-spec
	@go run ./cmd/lint-spec $(DIST_DIR)/openapi.payments.yaml

lint-redocly: payment-spec
	@npx --yes @redocly/cli@2.60.0 lint $(DIST_DIR)/openapi.payments.yaml

build: codegen
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-s -w" -o $(BIN_DIR)/server ./cmd/server
	@echo "==> Binary built at $(BIN_DIR)/server"

run: build
	$(BIN_DIR)/server

test: codegen
	./scripts/with-test-db.sh go test -v -race ./...

verify: build
	@chmod +x scripts/verify.sh
	./scripts/verify.sh

clean:
	rm -rf $(BIN_DIR) $(DIST_DIR) internal/api/v1/*.gen.go
