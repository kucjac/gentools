# Implementation Plan: Test Quality Guidance

**Branch**: `002-test-quality-guidance` | **Date**: 2026-09-02 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from [spec.md](spec.md)

## Summary

Replace the narrow nested protobuf check with a four-scenario consumer-validation
suite; make test gaps visible through a versioned behavior inventory plus coverage
evidence; and publish tested, static GitHub Pages documentation that links directly
to runnable examples. The design keeps Gentools a library: code-generation examples
are consumers of its public source-analysis API, not a new Gentools generator.

## Technical Context

**Language/Version**: Go 1.27.1, the repository-selected build, test, CI, and example toolchain.

**Primary Dependencies**: Existing standard-library and `golang.org/x/tools/go/packages` library dependencies; GitHub Pages Actions (`checkout`, `configure-pages`, `upload-pages-artifact`, and `deploy-pages`). No new runtime library dependency.

**Storage**: Version-controlled JSON behavior inventory, generated Markdown assessment report, fixture source, static documentation, and ephemeral coverage/snapshot artifacts.

**Testing**: Standard Go unit tests; four isolated nested-module consumer scenarios; coverage profiles as supporting evidence; fixture mutation checks; static-site link/build validation; CI and a release-only published-module smoke test.

**Target Platform**: Public Go module consumers on Linux, macOS, and Windows; GitHub Pages for public documentation.

**Project Type**: Public Go source-analysis library with development-only validation tooling and static documentation.

**Performance Goals**: The complete local gap assessment and integration validation complete in under 10 minutes from a warm clean checkout; no change to library runtime behavior or package-loading complexity.

**Constraints**: Preserve exported API compatibility; keep generated protobuf output generated and reproducible; distinguish candidate-source tests from released-module tests; do not treat raw line coverage as completeness; publish no credentials or private fixture data; retain the selected-toolchain and OS matrix.

**Scale/Scope**: Root `parser` and `types` packages, the existing nested integration module, one development command, one static documentation site, and at least one runnable code-generation example. New public Gentools APIs and a general-purpose generator are out of scope.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Pre-research gate: PASS**

- Library-first design: PASS. New tooling and examples consume existing public APIs; no new reusable runtime abstraction or public API is required.
- Public API compatibility: PASS. The plan only adds verification, fixtures, documentation, and a development command. Any correction to the protobuf fixture's declared package path is a fixture-reproducibility repair and must be tested.
- Test-first correctness: PASS. The behavior inventory records test level and evidence; every new scenario has both a positive result and an intentionally failing mutation check.
- Integration and compatibility testing: PASS. The design explicitly separates a staged candidate-source consumer contract from a release-only published-module smoke test.
- Idiomatic Go and simplicity: PASS. Standard Go tooling and a small development-only command avoid external coverage services and site-generator dependencies.

**Post-design gate: PASS**

The design adds no service, persistence system, public library API, or long-lived
generated output. GitHub Pages deployment receives only the minimum documented
permissions and a static public artifact. The behavior inventory is reviewable
source of truth; coverage is evidence rather than a policy gate.

## Project Structure

### Documentation (this feature)

```text
specs/002-test-quality-guidance/
├── plan.md              # This file ($speckit-plan command output)
├── research.md          # Phase 0 output ($speckit-plan command)
├── data-model.md        # Phase 1 output ($speckit-plan command)
├── quickstart.md        # Phase 1 output ($speckit-plan command)
├── contracts/           # Phase 1 output ($speckit-plan command)
└── tasks.md             # Phase 2 output ($speckit-tasks command - NOT created by $speckit-plan)
```

### Source Code (repository root)
```text
cmd/testinventory/                 # development-only inventory validation and report generation
docs/
├── index.html                     # static GitHub Pages entry point
├── code-generation.html           # guide that links to exact runnable example sources
└── testing/
    ├── behavior-inventory.json    # reviewed behavior-to-verification mapping
    └── behavior-inventory.md      # generated, reviewable assessment report

examples/codegen/struct-summary/   # runnable consumer that emits a Go source summary

internal/integration/
├── go.mod                          # nested consumer module
├── runner/                         # snapshots tracked candidate source into a temporary fixture
└── scenarios/
    ├── inspect/                    # successful public-contract inspection
    ├── generated/                  # generated protobuf source and imported dependencies
    ├── crosspackage/               # generic type use across consumer packages
    └── invalid/                    # expected diagnostic for invalid source/dependency

parser/ and types/                 # existing unit tests expanded from inventory gaps
scripts/
├── test-integration.sh             # scenario runner and mutation verification
├── assess-test-gaps.sh             # reproducible inventory/coverage entry point
└── verify-docs.sh                  # example, link, and static-site validation

.github/workflows/
├── ci.yml                          # root, inventory, docs, and candidate-source validation
├── pages.yml                       # main-only Pages deployment after build validation
└── release-smoke.yml               # tag/release published-module consumer smoke test
```

**Structure Decision**: Keep library code in `parser/` and `types/`; place
development-only checks outside those packages. Retain `internal/integration`
as the consumer-module boundary, but organize it by observable scenario rather
than one protobuf assertion. Use hand-authored static HTML because the requested
site is small and this avoids adding a site-generator toolchain.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
