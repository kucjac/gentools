# Data Model: CI and Lint Modernization

This feature has no runtime domain entities. It defines repository quality artifacts and their relationships.

## Toolchain declaration

| Field | Value/rule | Validation |
|---|---|---|
| Go version | `.go-version`, currently `1.27.1` | `make check-go`; CI setup-go |
| Linter version | one exact value, planned `v2.13.2` | local version output and CI action input |
| Source scope | root module `./...` | same command locally and in CI |
| Configuration | `.golangci.yml` in v2 schema | `golangci-lint config verify` |

## Lint policy

The policy consists of the v2 configuration version, enabled linters, settings, exclusions, timeout, output/annotation behavior, and exit status. Every retained setting must be accepted by the selected binary; unresolved findings must produce a non-zero process status.

## Validation evidence

Evidence is an observation, not application data. A valid record contains tool versions, configuration verification, the exact root command and result, deliberate-failure rule/location and non-zero result, hosted job/version/scope/status, and confirmation that the three-OS matrix still ran.

## Relationships and boundaries

`Toolchain declaration` selects the binary used by both local and hosted `Lint policy` evaluation. `Validation evidence` records results from that evaluation. The nested integration module is a separate existing test boundary and is not implicitly included in root lint scope.
