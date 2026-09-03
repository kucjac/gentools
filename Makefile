GO_VERSION := $(strip $(shell cat .go-version 2>/dev/null))
GOLANGCI_LINT_VERSION := v2.13.2
GO ?= go
GO_CMD = env -u GOROOT -u GOTOOLCHAIN $(GO)
GOMOD=$(GO_CMD) mod
GOTEST=echo "$(STEPS) Testing" && $(GO_CMD) test  -v ./...
GOVET=echo "$(STEPS) Vet" && $(GO_CMD) vet ./...
GOGENERATE=echo "$(STEPS) Generate by go:generate" && $(GO_CMD) generate -tags wireinject ./...
GOTIDY=echo "$(STEPS) Tidy modules" && $(GOMOD) tidy
GOLANGCI_LINT=echo "Golang CI Lint $(GOLANGCI_LINT_VERSION)..." && golangci-lint run --config .golangci.yml ./...

all: build-and-test golangci-lint

check-go:
	@if [ -z "$(GO_VERSION)" ]; then \
		echo "Missing Go version in .go-version" >&2; \
		exit 1; \
	fi
	@if ! command -v "$(GO)" >/dev/null 2>&1; then \
		echo "Go is required on PATH; install or configure Go $(GO_VERSION)" >&2; \
		exit 1; \
	fi
	@actual_version="$$($(GO_CMD) version 2>/dev/null | awk '{print $$3}')"; \
	expected_version="go$(GO_VERSION)"; \
	if [ "$$actual_version" != "$$expected_version" ]; then \
		echo "Go $$expected_version is required, found $${actual_version:-unknown} via $(GO)" >&2; \
		exit 1; \
	fi

test: check-go
	@$(GOTEST)

build-and-test: check-go
	@echo "Job: Build, Vet and Test"
	@(TOTAL=3 && \
 		CNT=1 && $(GOTIDY) && \
 		CNT=2 && $(GOVET) && \
 		CNT=3 && $(GOTEST))

golangci-lint: check-go
	@echo "Job: GolangCI Lint"
	@command -v golangci-lint >/dev/null 2>&1 || (echo "golangci-lint $(GOLANGCI_LINT_VERSION) is required on PATH" >&2; exit 1)
	@actual_version="$$(golangci-lint version 2>/dev/null | awk '/has version/{print "v"$$4; exit}')"; \
	if [ "$$actual_version" != "$(GOLANGCI_LINT_VERSION)" ]; then \
		echo "golangci-lint $(GOLANGCI_LINT_VERSION) is required, found $${actual_version:-unknown}" >&2; \
		exit 1; \
	fi
	@$(GOLANGCI_LINT)

test-integration: check-go
	@GODEBUG=gotypesalias=1 scripts/test-integration.sh all

test-gaps: check-go
	@GODEBUG=gotypesalias=1 scripts/assess-test-gaps.sh

verify-docs: check-go
	@GODEBUG=gotypesalias=1 scripts/verify-docs.sh

test-quality: test test-integration test-gaps verify-docs
