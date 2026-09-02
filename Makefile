GO ?= go
BIN := bin/openlane-mcp

.PHONY: all build test test-access vet fmt vuln check

all: check build

build:
	@mkdir -p bin
	$(GO) build -o $(BIN) ./cmd/openlane-mcp

test:
	$(GO) test ./...

test-access:
	@bash scripts/test-access.sh

vet:
	$(GO) vet ./...

fmt:
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)

vuln:
	$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

check: fmt vet test
