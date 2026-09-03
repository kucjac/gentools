# Research: CI and Lint Modernization

## Decision: Pin golangci-lint to v2.13.2

**Rationale**: The upstream release API reports v2.13.2, published 2026-08-28. Its v2.13.0 changelog explicitly records Go 1.27 support, and the v2.13.2 binary inspected during planning reports it was built with Go 1.27.0. This is compatible with the repository's Go 1.27.1 toolchain.

**Alternatives considered**: The specification's v2.12.2 assumption is outdated; v2.12.2 predates the upstream Go 1.27 support entry. `latest` is rejected because it violates reproducible CI and can change policy behavior without repository review.

Evidence: [v2.13.2 release](https://github.com/golangci/golangci-lint/releases/tag/v2.13.2), [upstream changelog](https://github.com/golangci/golangci-lint/blob/v2.13.2/CHANGELOG.md).

## Decision: Use the official GitHub Action with an exact version and shared policy

**Rationale**: Upstream CI guidance recommends the GitHub Action for GitHub projects and selecting a specific release. The action provides annotations while its exit status preserves the required gate. CI and the Make target must both run `golangci-lint run ./...` against the root module and `.golangci.yml`.

**Alternatives considered**: An unpinned shell-installed binary is less reproducible. `go install` is rejected by upstream guidance because it compiles with the local Go version and can alter dependency behavior. Linting the nested integration module is out of scope per the feature assumptions and would create a second source scope.

Evidence: [upstream CI installation guidance](https://golangci-lint.run/docs/welcome/install/ci/), [upstream local installation guidance](https://golangci-lint.run/docs/welcome/install/local/).

## Decision: Migrate `.golangci.yml` to v2 and make failures non-zero

**Rationale**: The verified v2.13.2 binary rejects the current file because it has no v2 `version` and contains legacy v1 keys. The current `issue-exit-code: 0` reports findings as success, contradicting FR-004. Use `golangci-lint migrate` as a starting point, then manually review/remove settings rejected by the v2 schema. The final config must pass `config verify` before lint findings are assessed.

**Alternatives considered**: Keeping v1 syntax is impossible with the selected v2 binary. Blindly accepting migration output is rejected because the validator identifies invalid options; each retained rule and exclusion needs intentional review.

## Decision: Keep lint isolated from the cross-platform test matrix

The existing `build-and-test` job covers Ubuntu, macOS, and Windows. The Ubuntu-only lint job is independently reported and does not mask platform-specific test failures. Change only its version and arguments while preserving the matrix.

## Decision: Document PATH-based local setup and evidence

`.go-version` and `Makefile` already establish a portable exact Go check. Add one exact linter version and an official-installer example using a PATH-visible user/bin directory, plus commands for version, config verification, clean lint, and deliberate-failure evidence. Host-specific SDK paths and credentials remain excluded.
