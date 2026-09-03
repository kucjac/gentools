# Implementation Plan: CI and Lint Modernization

**Branch**: `004-ci-lint-modernization` | **Date**: 2026-09-03 | **Spec**: [specs/004-ci-lint-modernization/spec.md](spec.md)

**Input**: Feature specification from [specs/004-ci-lint-modernization/spec.md](spec.md)

## Summary

Pin golangci-lint to the latest verified Go 1.27-compatible release, migrate the repository configuration to the linter's v2 schema, and make the local Make target and GitHub Actions lint job invoke the same deterministic command. The existing three-OS build-and-test matrix remains unchanged and lint remains a separate Ubuntu job. Contributor documentation will describe the required Go and linter versions, PATH-based installation, and evidence expected from local and hosted validation.

## Technical Context

**Language/Version**: Go 1.27.1 (`.go-version`); module language version `go 1.27.0`.

**Primary Dependencies**: Standard Go tooling; golangci-lint v2.13.2; GitHub Actions `golangci/golangci-lint-action@v9`; existing `actions/setup-go@v6`.

**Storage**: N/A. Configuration and policy are repository files; lint results are CI annotations/logs and command exit status.

**Testing**: `golangci-lint config verify`, `golangci-lint run ./...`, `make golangci-lint`, existing `make build-and-test`, and the existing nested-module, integration, documentation, and matrix checks.

**Target Platform**: Root module on Linux, macOS, and Windows; hosted lint on Ubuntu; nested integration module remains a separate test scope.

**Project Type**: Public Go source-analysis library with repository quality tooling and GitHub Actions CI.

**Performance Goals**: No library runtime change. Lint should remain bounded by the existing three-minute configuration timeout and use CI caching where the action supports it.

**Constraints**: Preserve public APIs and the existing three-OS build/test coverage. Use PATH-visible tools and repository-pinned versions; do not add host-specific SDK paths, credentials, or a second lint policy. Do not broaden lint scope to nested integration modules unless explicitly required later.

**Scale/Scope**: `.golangci.yml`, `Makefile`, `README.md`, `.github/workflows/ci.yml`, and feature validation artifacts. No production code, public API, dependency, or generated-source changes are planned.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

**Pre-research gate: PASS**

- Library-first design: PASS. This is quality tooling around the existing library and introduces no runtime abstraction.
- Public API compatibility: PASS. No exported identifier, signature, behavior, or representation changes are planned.
- Test-first correctness: PASS. The v2 configuration is verified, a clean lint run is required, and a deliberate violation must demonstrate non-zero failure.
- Integration and compatibility testing: PASS. Existing root, nested-module, consumer, documentation, and three-OS checks remain in scope; lint is independently diagnosable.
- Idiomatic Go and simplicity: PASS. The design uses the supported upstream action and one shared version/policy rather than adding a custom linter wrapper.

**Post-design gate: PASS**

The design changes only repository automation, configuration, and guidance. It adds no project, service, persistence, public API, or dependency complexity. The only intentional compatibility boundary is the already-declared Go 1.27.1 toolchain, which this feature preserves.

## Project Structure

### Documentation (this feature)

```text
specs/004-ci-lint-modernization/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── validation-contract.md
└── tasks.md                         # created by $speckit-tasks
```

### Source and automation

```text
.golangci.yml                         # v2 policy and explicit non-zero gate
Makefile                              # canonical local lint command/version check
.github/workflows/ci.yml              # pinned hosted lint job
README.md                             # contributor tool/version guidance
parser/, types/, examples/, cmd/      # root-module lint scope
internal/integration/                  # separate existing test module, not lint scope
```

**Structure Decision**: Keep the existing repository layout. The feature is implemented in configuration, automation, and documentation files; no source package or module is added.

## Complexity Tracking

No constitution violations or complexity exceptions require justification.
