# Implementation Plan: Generic Parser Compatibility

**Branch**: `001-generic-parser-compatibility` | **Date**: 2026-09-02 | **Spec**: [specs/001-generic-parser-compatibility/spec.md](spec.md)

**Input**: Feature specification from [specs/001-generic-parser-compatibility/spec.md](spec.md)

## Summary

Make package loading and the public type model correctly describe Go generic
declarations, instantiations, generic aliases, and methods on generic receiver types
without changing established non-generic contracts. The implementation will
centralize conversion from `go/types` into additive Gentools generic metadata,
use the current `go/types` APIs directly, and prove behavior with fixture-driven
unit tests, nested-module integration tests, and a single-toolchain CI workflow.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go library; Go 1.27.1 is the only supported build, test, and CI toolchain.

**Primary Dependencies**: Standard-library `go/ast` and `go/types`; `golang.org/x/tools/go/packages` and `golang.org/x/mod/modfile` selected for Go 1.27.1.

**Storage**: N/A; package and type metadata exist in memory for the duration of loading.

**Testing**: Go unit tests in `parser/` and `types/`; fixture packages under `parser/testcases/`; the nested integration module under `internal/integration/`; CI runs formatting, module consistency, vet, tests, and lint with Go 1.27.1.

**Target Platform**: Public Go modules and source packages on Linux, macOS, and Windows; CI retains its existing operating-system matrix.

**Project Type**: Public source-analysis library.

**Performance Goals**: No regression in package-loading complexity: generic metadata is derived during the existing load/convert pass and must not require a second whole-package load.

**Constraints**: Preserve all existing exported API signatures and non-generic observable results. Go 1.27.1-only `go/types` APIs may be used directly. Keep generated protobuf code unchanged. Local validation is deliberately serial to limit memory use.

**Scale/Scope**: One root module plus one nested integration module; cover generic type parameters, constraints, named/alias types, function signatures, methods, and declared type uses. Expression-level source indexing that is unrelated to the existing public declaration model is out of scope.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Pre-research gate: PASS**

- Library-first design: PASS. The change remains in reusable parser and type-model packages; CI changes only prove the library contract.
- Public API compatibility: PASS with an additive-only design. Existing exported types and methods are not reshaped; new generic metadata is exposed through new types and accessors. Existing concrete types remain the representation of non-generic declarations.
- Test-first correctness: PASS. Every converter path and generic language category receives a regression fixture and focused assertions, including controlled unsupported-input outcomes.
- Integration and compatibility testing: PASS. Cross-package fixtures and the existing nested integration module are part of the validation matrix.
- Idiomatic Go and simplicity: PASS. Conversion stays centered on `go/types` with no legacy compatibility adapter required.

**Post-design gate: PASS**

The design adds no runtime service, persistence layer, or extra project. The single Go 1.27.1 support boundary removes the need for version-adapter complexity.

## Project Structure

### Documentation (this feature)

```text
specs/001-generic-parser-compatibility/
├── plan.md              # This file ($speckit-plan command output)
├── research.md          # Phase 0 output ($speckit-plan command)
├── data-model.md        # Phase 1 output ($speckit-plan command)
├── quickstart.md        # Phase 1 output ($speckit-plan command)
├── contracts/           # Phase 1 output ($speckit-plan command)
└── tasks.md             # Phase 2 output ($speckit-tasks command - NOT created by $speckit-plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
parser/
├── load-packages.go             # package loading and go/types conversion
├── expr.go                      # AST helpers used for declaration metadata
├── generic_*.go                 # Go 1.27.1 generic conversion
├── *_test.go                    # loader and regression tests
└── testcases/
    ├── models.go                # existing baseline fixtures
    ├── generic/                 # generic declaration fixtures
    └── genericimport/           # cross-package generic fixtures

types/
├── type.go                      # shared Type contract and generic accessors
├── generic.go                   # additive generic metadata and instantiated types
├── {alias,function,interface,struct}.go
└── *_test.go

internal/integration/
├── go.mod                       # nested integration module
└── protobuf/                    # generated-model integration coverage

.github/workflows/ci.yml         # selected-toolchain and platform validation
```

**Structure Decision**: Keep the existing package layout. Parser-specific
translation belongs in `parser/`; reusable public representations belong in
`types/`; fixtures stay alongside current parser test fixtures. No new module,
application, or generated source is introduced.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
