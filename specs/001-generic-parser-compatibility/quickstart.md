# Quickstart: Validate Generic Parser Compatibility

## Prerequisites

- A Go toolchain installed and selected explicitly for each validation run.
- Network access for first-time module downloads, or a warmed module cache.
- A checkout at the repository root.

The `.go-version` file selects Go 1.27.1. Run each command serially with one
build worker to stay within the available memory budget.

## Run root validation

From the repository root, run:

```sh
GOMAXPROCS=1 GOFLAGS=-p=1 GODEBUG=gotypesalias=1 go test -count=1 ./...
```

Expected result: every package test passes, including generic aliases and
generic-method coverage.

## Run integration validation

The repository has a nested module that root `./...` does not include. Run it
separately with Go 1.27.1:

```sh
cd internal/integration
GOMAXPROCS=1 GOFLAGS=-p=1 GODEBUG=gotypesalias=1 go test -count=1 ./...
```

Expected result: the protobuf integration test loads generated declarations
through the updated root module without regression.

## Run static validation

```sh
go fmt ./...
git diff --exit-code -- '*.go'
go vet ./...
```

Run the project linter using Go 1.27.1 as selected by the workflow.

## CI verification

Push the branch or open a pull request and inspect the job names in
[.github/workflows/ci.yml](../../../.github/workflows/ci.yml). Confirm build
and test results for Go 1.27.1 on every operating system; confirm the separate
lint job identifies its toolchain and operating system.

Expected result: formatting leaves no Git diff, and static checks exit successfully.

## Updating the selected toolchain

Update `.go-version`, the root `go` directive, and the pinned workflow version
together. Retain all three operating systems and add a focused fixture only
when the newly selected Go version adds syntax or `go/types` behavior. The job
name reports the selected version, OS, and validation category.

## Recorded local evidence

Validated on 2026-09-02 with Go 1.27.1, `GOROOT` unset, `GOMAXPROCS=1`,
`GOFLAGS=-p=1`, `GODEBUG=gotypesalias=1`, and a 6 GiB virtual-memory cap:

```text
go vet ./...                         PASS
go test -count=1 ./...               PASS
(cd internal/integration && go vet ./...)               PASS
(cd internal/integration && go test -count=1 ./...)     PASS
bash -n scripts/verify-ci-matrix.sh                      PASS
scripts/verify-ci-matrix.sh                              PASS
```

The cap is intentional: it confirms the generic-alias loader no longer grows
without bound when an alias target cannot be resolved.

For expected entities and failure behavior, see
[specs/001-generic-parser-compatibility/contracts/public-type-model.md](contracts/public-type-model.md)
and [specs/001-generic-parser-compatibility/data-model.md](data-model.md).
