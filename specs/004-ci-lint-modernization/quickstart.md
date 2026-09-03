# Quickstart: CI and Lint Modernization

## Prerequisites

Use Go `1.27.1` from `PATH`, as declared in `.go-version`. Install the exact binary release `golangci-lint v2.13.2` using the official installer:

```sh
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v2.13.2
```

Confirm the versions:

```sh
go version
golangci-lint version
```

## Validate configuration and local lint

```sh
golangci-lint config verify --config .golangci.yml
make golangci-lint
```

Expected result: configuration verification succeeds and the canonical root lint command exits zero on a clean revision. Findings must include a location and rule and cause a non-zero result.

## Validate failure behavior

In a temporary uncommitted change, introduce a violation covered by an enabled rule and run `make golangci-lint`. Expected result: non-zero status and an actionable rule/location. Revert the temporary change afterward.

## Validate preserved coverage

```sh
make build-and-test
make test-integration
make verify-docs
```

The hosted workflow must additionally show successful build-and-test runs on Ubuntu, macOS, and Windows, plus a separately reported Ubuntu lint job using v2.13.2. See [the validation contract](contracts/validation-contract.md) for the exact evidence requirements.
