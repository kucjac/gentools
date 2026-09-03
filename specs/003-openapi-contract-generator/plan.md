# Implementation Plan: OpenAPI Contract Generator

**Branch**: `feat/003-openapi-contract-generator` | **Date**: 2026-09-03 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/003-openapi-contract-generator/spec.md`

## Summary

Add a runnable, consumer-owned example that loads annotated Go source through
Gentools, validates a small operation annotation convention, and writes a
deterministic OpenAPI 3 JSON document.

## Technical Context

**Language/Version**: Go 1.27.1 (`.go-version`), module language Go 1.27.0.

**Primary Dependencies**: Standard library JSON support and existing `parser`
and `types` packages; no new dependency.

**Storage**: Filesystem input and atomically replaced JSON output.

**Testing**: `go test`, example command, tracked-source snapshot validation,
and `make verify-docs`.

**Target Platform**: Portable Go environments covered by repository CI.

**Project Type**: Public Go library with consumer-facing code-generation examples.

**Performance Goals**: Generate the small documented fixture in under one minute.

**Constraints**: Deterministic JSON, actionable diagnostics, no partial output
on validation failure, backward-compatible library behavior, no network or credentials.

**Scale/Scope**: Educational consumer example with documented GET operations
and public structs; path parameters, primitive/pointer/slice/struct fields,
JSON tags, and shared definitions.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Pre-research gate: PASS**

- Library-first design: PASS. This feature is explicitly a consumer-facing educational example, not a new Gentools library capability; it consumes existing library metadata and does not claim a reusable public generator API.
- Public API compatibility: PASS. No existing exported API changes are planned.
- Test-first correctness: PASS. Success, diagnostics, duplicate detection, unresolved models, unsupported shapes, and output safety are tested.
- Integration and compatibility testing: PASS. The example is run from an isolated tracked-source snapshot and through documentation validation.
- Idiomatic Go and simplicity: PASS. OpenAPI JSON is modeled with standard-library structs and no external dependency.

**Post-design gate: PASS**

The design adds only an example generator, fixtures, tests, and documentation
validation. It does not alter runtime library behavior or introduce a dependency.

## Project Structure

```text
specs/003-openapi-contract-generator/{plan,research,data-model,quickstart}.md
specs/003-openapi-contract-generator/contracts/validation-contract.md
examples/codegen/openapi-contract/
├── main.go
├── main_test.go
└── testdata/{models.go,handlers.go,openapi.golden.json}
scripts/verify-docs.sh
scripts/verify-docs_test.sh
docs/code-generation.html
README.md
```

**Structure Decision**: Keep the implementation beside the existing codegen
examples. The command accepts an input directory and output path, uses
`parser.LoadPackages` with comments enabled, and writes only after complete
validation. The docs verifier runs it against the committed tracked-source
snapshot and compares normalized output with the golden contract.

## Complexity Tracking

No constitution violations or complexity exceptions require justification.
