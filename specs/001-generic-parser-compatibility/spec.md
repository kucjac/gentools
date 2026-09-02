# Feature Specification: Generic Parser Compatibility

**Feature Branch**: `001-generic-parser-compatibility`

**Created**: 2026-09-02

**Status**: Revised

**Input**: Support generic Go source using only Go 1.27.1, the repository's selected current toolchain. The repository is maintained by one user, so compatibility with older Go releases is not required.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Analyze Generic Declarations Reliably (Priority: P1)

As the library consumer, I want package loading and source analysis to describe generic declarations and their uses so generation and inspection work with my current Go source.

**Why this priority**: Correct generic handling is the central defect and the library must remain useful for the selected current toolchain.

**Independent Test**: Load packages containing generic types, functions, interfaces, aliases, constraints, instantiated uses, and methods; verify the public model exposes their names, relationships, kinds, and signatures without an error or panic.

**Acceptance Scenarios**:

1. **Given** a package containing generic declarations and constraints, **When** it is loaded with the selected toolchain, **Then** its type parameters, constraints, owner scope, and declaration metadata are queryable.
2. **Given** a package that instantiates local or imported generic types, **When** it is loaded, **Then** the origin and ordered type arguments remain distinguishable from the uninstantiated declaration.
3. **Given** a package containing generic aliases or methods on generic receiver types, **When** it is loaded, **Then** their identity and metadata are retained without changing the established concrete Gentools declaration type.

### User Story 2 - Upgrade Existing Consumers Safely (Priority: P1)

As the sole library consumer, I want existing non-generic source and public API usage to keep working after generic support is added.

**Why this priority**: The toolchain support boundary is changing, but existing library contracts should remain stable for current source.

**Independent Test**: Run the unchanged baseline fixtures and nested protobuf integration test with the selected toolchain; existing named declarations retain their concrete types and observable behavior.

**Acceptance Scenarios**:

1. **Given** existing non-generic packages and callers, **When** they are analyzed, **Then** their names, kinds, equality, fields, methods, signatures, and rendering remain unchanged.
2. **Given** malformed current-toolchain source, **When** it is loaded, **Then** the caller receives the original compiler or tooling diagnostic rather than a panic or a successful partial result.

### User Story 3 - Trust Current-Toolchain Evidence (Priority: P2)

As the maintainer, I want one repeatable current-toolchain validation workflow so regressions are caught without resource-heavy multi-version testing.

**Why this priority**: A single selected toolchain reflects the repository's maintenance model and keeps validation practical on the available machine.

**Independent Test**: Run the root and nested-module validation commands once with the selected toolchain; inspect the workflow to confirm it runs those same checks on every retained operating system.

**Acceptance Scenarios**:

1. **Given** the selected toolchain, **When** CI runs, **Then** dependency consistency, formatting, vetting, root tests, nested-module tests, and lint have separately identifiable results.
2. **Given** a future toolchain update, **When** the maintainer follows the documented coordinated update procedure, **Then** the workflow and validation instructions remain understandable without a compatibility matrix.

### Edge Cases

- A declaration has multiple type parameters, repeated names in separate scopes, or constraints that use unions and approximation terms.
- A generic type appears inside pointers, slices, arrays, maps, channels, function signatures, structs, interfaces, or aliases.
- A generic declaration is imported from another package with a different package identifier.
- Source is malformed or cannot be type-checked by the selected toolchain.
- Existing non-generic declarations resemble generic syntax but must retain their public behavior.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The library MUST load and analyze generic declarations supported by the selected current toolchain without false rejection.
- **FR-002**: The library MUST expose type parameters, constraints, declaration scope, and instantiated type arguments through additive public behavior.
- **FR-003**: The library MUST preserve the origin and ordered arguments of generic type instantiations across package boundaries.
- **FR-004**: The library MUST preserve existing non-generic exported APIs and observable results for source supported by the selected toolchain.
- **FR-005**: Generic aliases and methods on generic receiver types supported by Go 1.27.1 MUST retain their identity, parameters, constraints, receiver information, inputs, results, and type arguments. Methods do not declare their own type parameters in Go.
- **FR-006**: Compiler and tooling diagnostics MUST remain actionable; generic conversion MUST not panic or return a successful partially initialized value.
- **FR-007**: Regression coverage MUST include successful generic declarations, constraints, aliases, instantiations, cross-package use, nested composition, malformed source, and existing non-generic behavior.
- **FR-008**: CI MUST run the applicable checks using only the selected current toolchain while retaining the repository's operating-system coverage.
- **FR-009**: Documentation MUST state the single-toolchain policy and the procedure for updating that selected toolchain.

### Success Criteria *(mandatory)*

- **SC-001**: All root and nested-module tests pass with the selected toolchain.
- **SC-002**: Focused regression tests cover each generic category listed in FR-007 and fail when their covered behavior is disabled.
- **SC-003**: CI identifies the validation category and operating system for every failure without multiplying work across legacy toolchains.
- **SC-004**: A representative package with generic declarations and imported instantiations loads in one run without a panic or loss of declared type arguments.
- **SC-005**: A representative non-generic package produces the same public results before and after the feature, except for explicitly additive generic information.

## Assumptions

- The `.go-version` file is the repository authority for the selected current toolchain and names Go 1.27.1.
- Go 1.19 and other historical Go releases are outside this feature's support boundary; compatibility adapters and multi-version CI are not required.
- Existing operating-system coverage remains in scope.
- Generated protobuf code remains unchanged and is exercised only through its integration test.
