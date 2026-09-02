<!--
Sync Impact Report
- Version change: unratified scaffold -> 1.0.0
- Modified principles: none; all five principles established from the scaffold
- Added sections: Additional Constraints; Development Workflow
- Removed sections: none
- Follow-up TODOs: none
-->

# Gentools Constitution

## Core Principles

### I. Library-First Design

Gentools features MUST be implemented as focused, reusable Go library
capabilities. Public behavior MUST be exposed through clear package APIs, with
implementation details kept private unless they are part of the documented
contract. Each feature MUST have a stated purpose and an independently
testable behavior.

Rationale: A small, composable library API is the project’s primary product and
keeps code generation and source analysis reusable by downstream applications.

### II. Public API Compatibility

Exported identifiers, interfaces, method signatures, documented behavior,
error contracts, and observable representations MUST be treated as compatibility
contracts. Changes MUST preserve existing callers and valid inputs unless a
breaking change is explicitly approved, documented with migration guidance, and
released under the appropriate major-version policy. New behavior MUST NOT
silently change established behavior for existing inputs.

Rationale: Gentools is consumed as a public Go module; downstream users need
stable upgrades and predictable source compatibility.

### III. Test-First Correctness

Every behavior change MUST include tests that express the intended contract.
Bug fixes MUST include a regression test. Tests MUST cover successful behavior,
invalid inputs, and error or panic behavior where those outcomes are part of
the contract. Implementation work MUST pass the relevant test suite before
review.

Rationale: AST and type-reflection behavior has many edge cases, and executable
contracts prevent regressions in both public APIs and parsing semantics.

### IV. Integration and Compatibility Testing

Changes affecting package loading, dependency resolution, generated models,
cross-package type relationships, or public contracts MUST include integration
coverage in addition to unit tests. Compatibility-sensitive changes MUST test
representative existing usage patterns and supported Go environments.

Rationale: Correctness depends on interactions among parser, types, Go tooling,
and imported packages that isolated unit tests cannot fully represent.

### V. Idiomatic Go and Simplicity

Code MUST follow established Go conventions, provide documentation for
exported declarations, and keep dependencies and abstractions minimal. A
design MUST justify added complexity, and a simpler implementation MUST be
preferred when it satisfies the documented contract and performance needs.

Rationale: Idiomatic, documented code lowers adoption and maintenance costs for
library users while keeping the implementation understandable.

## Additional Constraints

The module MUST remain buildable with the Go versions declared by the project’s
supported CI and module configuration. Changes MUST use standard Go tooling and
MUST preserve portability across the operating systems covered by CI unless an
explicit support change is approved.

Dependencies MUST have a documented need, compatible licensing, and a stable
version selection. Generated code MUST be reproducible from its source inputs
and MUST NOT be edited manually when regeneration is available.

Public APIs MUST include Go documentation. New exported symbols MUST avoid
unnecessary exposure of internal representation, and API changes MUST update
examples or documentation when usage changes.

## Development Workflow

Work MUST begin with a clear requirement or defect statement and an identified
public or internal contract. Before implementation, the change MUST identify
compatibility impact, affected packages, and required unit and integration
tests.

A change MUST be reviewed against the constitution, formatted and statically
validated with the project’s Go tooling, and tested with the applicable test
commands. Pull requests MUST describe compatibility impact, test evidence, and
any approved exceptions. Releases MUST communicate breaking changes and
migration steps when applicable.

## Governance

This constitution governs project changes and takes precedence over informal
practice. Amendments MUST be proposed in writing, include a Sync Impact Report,
state their compatibility and migration impact, and receive maintainer review
before being merged. An amendment MUST update the version and last-amended
date in the document.

Constitution versions follow semantic versioning: MAJOR denotes removed or
redefined governance obligations; MINOR denotes a new principle or materially
expanded obligation; PATCH denotes clarification, correction, or other
non-semantic refinement. The project MUST use the smallest applicable bump.

Each change review MUST verify compliance with applicable principles and record
any approved exception and its expiry or follow-up. Maintainers MUST review
this constitution when release policy, supported Go versions, or public API
compatibility practices materially change.

**Version**: 1.0.0 | **Ratified**: 2026-09-02 | **Last Amended**: 2026-09-02
