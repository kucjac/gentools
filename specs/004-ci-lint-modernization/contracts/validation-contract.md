# Validation Contract

## Canonical local contract

With Go `1.27.1` and golangci-lint `v2.13.2` on `PATH`, from repository root:

```sh
make golangci-lint
```

The command must enforce the Go version, load `.golangci.yml`, analyze root packages matching `./...`, print findings, and exit zero only when no unresolved findings exist. The direct equivalent is `golangci-lint run ./...` after `golangci-lint config verify`.

## Hosted CI contract

The `golangci-lint` job must run on `ubuntu-latest` after checkout and Go 1.27.1 setup, select exactly v2.13.2 through the official action, use `.golangci.yml` and root scope `./...`, fail on findings with actionable annotations/logs, and remain separate from the Ubuntu/macOS/Windows `build-and-test` matrix.

## Failure and compatibility contracts

Invalid configuration, missing required tool, unsupported version, or lint finding is a non-zero validation failure. A deliberate violation must identify a source location and rule. This feature must not change exported Go APIs, runtime behavior, or existing three-OS build/test coverage; nested integration tests retain their existing commands.
