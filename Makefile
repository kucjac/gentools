GO=go
GOMOD=$(GO) mod
GOTEST=echo "$(STEPS) Testing" && $(GO) test -count 1 -v ./...
GOVET=echo "$(STEPS) Vet" && $(GO) vet ./...
GOGENERATE=echo "$(STEPS) Generate by go:generate" && $(GO) generate -tags wireinject ./...
GOTIDY=echo "$(STEPS) Tidy modules" && $(GOMOD) tidy
GOLANGCI_LINT=golangci-lint run ./...

all: build-and-test golangci-lint

build-and-test:
	@echo "Job: Build, Vet and Test"
	@(TOTAL=3 && \
 		CNT=1 && $(GOTIDY) && \
 		CNT=2 && $(GOVET) && \
 		CNT=3 && $(GOTEST))

golangci-lint:
	@echo "Job: GolangCI Lint"
	@$(GOLANGCI_LINT)