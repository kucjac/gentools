---
description: "Implementation tasks for Generic Parser Compatibility"
---

# Tasks: Generic Parser Compatibility

**Input**: Design documents from `specs/001-generic-parser-compatibility/`

**Prerequisites**: `spec.md`, `plan.md`, `research.md`, `data-model.md`, `contracts/public-type-model.md`, and `quickstart.md`

**Tests**: Required by the feature specification and constitution. Add each focused regression test before its implementation, and confirm it fails for the defect it covers.

**Organization**: Tasks are grouped by user story so each delivers an independently testable increment.

## Path Conventions

- Public model: `types/`
- Package loading and `go/types` conversion: `parser/`
- Parser fixtures: `parser/testcases/`
- Nested-module integration coverage: `internal/integration/`
- Compatibility automation: `.github/workflows/ci.yml`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish current-toolchain fixtures and a reproducible serial validation baseline without changing public behavior.

- [X] T001 Record the selected Go 1.27.1 dependency set and single-toolchain rationale in `specs/001-generic-parser-compatibility/research.md`.
- [X] T002 Create version-gated generic fixture directories and package documentation in `parser/testcases/generic/doc.go`, `parser/testcases/genericalias/doc.go`, and `parser/testcases/genericmethods/doc.go`.
- [X] T003 Add a CI-ready root and nested-module validation command matrix to `specs/001-generic-parser-compatibility/quickstart.md`.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Provide additive public generic metadata and a Go-1.27 conversion boundary before loading generic declarations.

**⚠️ CRITICAL**: Complete this phase before implementing either parser behavior or CI coverage.

- [X] T004 Add additive `TypeParameter`, `GenericInfo`, `Constraint`, and `Instantiation` public representations, including `Type` behavior and documented accessors, in `types/generic.go`.
- [X] T005 Add generic-metadata accessors to existing named declarations without adding exported fields or changing their existing APIs in `types/alias.go`, `types/function.go`, `types/interface.go`, and `types/struct.go`.
- [X] T006 Extend the public kind and string/identity behavior for type parameters and instantiated types in `types/kind.go` and `types/generic.go`.
- [X] T007 Add current-toolchain generic conversion that returns actionable errors instead of panicking or producing partial values in `parser/generic_common.go`.
- [X] T008 Use the Go-1.27 `go/types` alias APIs for generic aliases in `parser/generic_go122.go`.
- [X] T009 Request the syntax and type information required to map generic declarations and instantiation sites, while preserving existing comment behavior, in `parser/load-packages.go`.

**Checkpoint**: The additive model and current-toolchain conversion boundary compile on Go 1.27.1; story work can now proceed.

---

## Phase 3: User Story 1 - Analyze Generic Declarations Reliably (Priority: P1) 🎯 MVP

**Goal**: Load generic declarations and uses into the additive type model without panics, false rejections, or lost generic metadata.

**Independent Test**: Load the generic fixture packages on Go 1.27.1 and assert declaration parameters, constraints, instantiated use origins/arguments, cross-package relationships, and generic-receiver method metadata.

### Tests for User Story 1

- [X] T010 [P] [US1] Add generic type, function, constraint, nested-composite, and invalid-input fixtures in `parser/testcases/generic/models.go` and `parser/generic_error_test.go`.
- [X] T011 [P] [US1] Add a separate imported generic declaration and instantiated-use fixture package in `parser/testcases/genericimport/models.go` and `parser/testcases/genericimport/consumer/models.go`.
- [X] T012 [P] [US1] Add generic-alias fixtures, including constraints and aliases used across a package boundary, in `parser/testcases/genericalias/models_go124.go`.
- [X] T013 [P] [US1] Add generic-receiver method fixtures covering receiver parameters, variadics, inputs, and results in `parser/testcases/genericmethods/models_go127.go`.
- [X] T014 [US1] Add focused loader assertions for generic declarations, constraints, type-parameter scope, nested compositions, instantiations, and controlled load failures in `parser/generic_test.go`.
- [X] T015 [US1] Add focused cross-package generic-alias and generic-method assertions in `parser/generic_cross_package_test.go`, `parser/generic_alias_go124_test.go`, and `parser/generic_methods_go127_test.go`.

### Implementation for User Story 1

- [X] T016 [US1] Convert `go/types.TypeParam` values and their constraints into the additive public model, retaining owner and ordinal information, in `parser/generic_common.go`.
- [X] T017 [US1] Attach declaration and receiver generic metadata to existing `Struct`, `Interface`, `Alias`, and `Function` values during scaffold and signature completion in `parser/load-packages.go`.
- [X] T018 [US1] Recognize instantiated named types and generic signatures before existing named/composite fallbacks, preserving origin and ordered type arguments, in `parser/generic_common.go` and `parser/load-packages.go`.
- [X] T019 [US1] Translate generic aliases through the Go 1.27.1 `go/types` APIs while retaining alias identity, parameters, constraints, right-hand side, and arguments in `parser/generic_go122.go`.
- [X] T020 [US1] Preserve generic receiver, ordinary parameter/result, and variadic metadata when converting methods in `parser/load-packages.go` and `parser/generic_common.go`.
- [X] T021 [US1] Replace generic conversion panics and silent mapping omissions with contextual, actionable errors that flow through package loading in `parser/load-packages.go` and `parser/generic_common.go`.

**Checkpoint**: Generic fixtures load correctly on Go 1.27.1 and each parser regression test identifies its affected construct.

---

## Phase 4: User Story 2 - Upgrade Existing Consumers Safely (Priority: P1)

**Goal**: Preserve existing non-generic public contracts, type identities, rendering, errors, and integration behavior while adding generic metadata.

**Independent Test**: Run the unchanged baseline fixture and nested protobuf integration tests on Go 1.27.1; public named declarations retain their existing concrete Gentools types and observable results.

### Tests for User Story 2

- [X] T022 [P] [US2] Add compatibility assertions that existing non-generic declarations retain concrete type, name, full name, kind, equality, method, field, and function-signature behavior in `parser/non_generic_compatibility_test.go`.
- [X] T023 [P] [US2] Add public-model compatibility tests proving new generic accessors are additive and pre-generic values remain empty/unchanged in `types/generic_test.go`.
- [X] T024 [P] [US2] Extend the nested-module regression test to verify existing generated protobuf loading behavior is unchanged in `internal/integration/protobuf/proto_test.go`.
- [X] T025 [US2] Add regression assertions for malformed, unsupported, and toolchain-ineligible source diagnostics in `parser/generic_error_test.go`.

### Implementation for User Story 2

- [X] T026 [US2] Preserve old concrete declaration values, unkeyed-composite compatibility, type identity, and rendering while exposing generic data solely through new accessors in `types/generic.go`, `types/alias.go`, `types/function.go`, `types/interface.go`, and `types/struct.go`.
- [X] T027 [US2] Keep non-generic alias, named-type, field, method, and declaration conversion on their established paths while delegating only recognized generic constructs in `parser/load-packages.go` and `parser/expr.go`.
- [X] T028 [US2] Ensure package-loader diagnostics retain underlying compiler/tooling errors and never publish successful partially initialized generic values in `parser/load-packages.go`.
- [X] T029 [US2] Document the additive API, supported-toolchain behavior, and controlled failure semantics for consumers in `README.md` and `specs/001-generic-parser-compatibility/contracts/public-type-model.md`.

**Checkpoint**: Existing callers and fixtures work without source changes, including the nested protobuf integration module.

---

## Phase 5: User Story 3 - Trust Compatibility Evidence (Priority: P2)

**Goal**: Make CI deliver attributable, repeatable compatibility evidence for every required Go version and existing operating system.

**Independent Test**: Inspect rendered workflow job names and execute the documented commands to confirm Go 1.27.1 runs root build/vet/test across Linux, macOS, and Windows, while lint runs as its separately documented job.

### Tests for User Story 3

- [X] T030 [P] [US3] Add a workflow-structure regression script that asserts the explicit Go-version/OS matrix, checkout/setup action versions, and required validation steps in `scripts/verify-ci-matrix.sh`.
- [X] T031 [P] [US3] Add fixture-selection assertions documenting the language-version requirements of generic aliases and generic-receiver methods in `parser/generic_version_gates_test.go`.

### Implementation for User Story 3

- [X] T032 [US3] Configure Go 1.27.1 across the existing Linux, macOS, and Windows build-and-test matrix, using maintained checkout and Go-setup actions, in `.github/workflows/ci.yml`.
- [X] T033 [US3] Run root module dependency consistency, formatting, vet, and tests with version/OS/category-attributable step names in `.github/workflows/ci.yml`.
- [X] T034 [US3] Add a separately named nested-module test step for `internal/integration` so it is validated on every applicable matrix entry in `.github/workflows/ci.yml`.
- [X] T035 [US3] Move lint to a documented compatible Go/OS subset with a maintained linter action and explicit job name, keeping it separate from build-and-test compatibility failures, in `.github/workflows/ci.yml`.
- [X] T036 [US3] Document the coordinated selected-toolchain update procedure, syntax-specific fixture gates, and expected job evidence in `.github/workflows/ci.yml` and `specs/001-generic-parser-compatibility/quickstart.md`.

**Checkpoint**: CI failures identify Go 1.27.1, the OS, and validation category; a future toolchain change updates the documented selected version and any new syntax-specific fixtures.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Validate the completed implementation, preserve project hygiene, and record compatibility evidence.

- [X] T037 [P] Run formatting and static validation for both modules, recording commands and outcomes in `specs/001-generic-parser-compatibility/quickstart.md`.
- [X] T038 Run the full root and nested-module test suite serially on Go 1.27.1, recording the result in `specs/001-generic-parser-compatibility/quickstart.md`.
- [X] T039 [P] Review exported generic symbols and consumer examples for Go documentation and contract consistency in `types/generic.go`, `README.md`, and `specs/001-generic-parser-compatibility/contracts/public-type-model.md`.
- [X] T040 Audit `go.mod` and `go.sum` so the chosen dependency set builds on Go 1.27.1 and is justified by the compatibility evidence in `specs/001-generic-parser-compatibility/research.md`.
- [X] T041 Verify all repository and feature-document links resolve, the workflow YAML is valid, and task evidence meets the constitution in `specs/001-generic-parser-compatibility/tasks.md`, `.github/workflows/ci.yml`, and `.specify/memory/constitution.md`.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependency; establishes fixtures and validation documentation.
- **Foundational (Phase 2)**: Depends on setup; blocks every user story because it supplies the additive model and Go-version conversion boundary.
- **US1 (Phase 3)**: Depends on Phase 2; delivers the generic parsing MVP.
- **US2 (Phase 4)**: Depends on Phase 2 and validates US1's additive implementation against existing public behavior.
- **US3 (Phase 5)**: Depends on the implementation and syntax-specific fixtures from US1/US2 so CI exercises real current-toolchain behavior.
- **Polish (Phase 6)**: Depends on the desired stories and records final evidence.

### User Story Dependencies

- **US1 (P1)**: Starts after Phase 2 and is the MVP.
- **US2 (P1)**: Starts after Phase 2, but its complete compatibility proof follows the US1 conversion changes.
- **US3 (P2)**: Requires the generic fixtures and validation commands from US1 and the compatibility cases from US2.

### Parallel Opportunities

- T002 and T003 can proceed together after the research decision in T001.
- T010–T013 fixture files can be authored in parallel; T014–T015 then aggregate their assertions.
- T022–T024 are independent regression test files and can run in parallel.
- T030 and T031 can proceed in parallel before workflow implementation.
- T037 and T039 can run in parallel after the code and documentation changes are complete.

## Parallel Example: User Story 1

```text
Task: "Add generic declaration fixtures in parser/testcases/generic/models.go and malformed-source coverage in parser/generic_error_test.go"
Task: "Add cross-package generic fixtures in parser/testcases/genericimport/models.go and parser/testcases/genericimport/consumer/models.go"
Task: "Add Go-1.24 generic-alias fixtures in parser/testcases/genericalias/models_go124.go"
Task: "Add generic-receiver method fixtures in parser/testcases/genericmethods/models_go127.go"
```

## Implementation Strategy

### MVP First (US1)

1. Complete Phases 1 and 2, including the Go 1.27.1 conversion boundary.
2. Write the US1 fixtures and focused failing assertions.
3. Implement generic metadata and conversion paths (T016–T021).
4. Validate generic parsing on Go 1.27.1.

### Incremental Delivery

1. Deliver US1 as correct generic parsing.
2. Add US2's unchanged-consumer and error-contract proof before treating the feature as compatible.
3. Add US3's Go 1.27.1/OS CI proof so future toolchain changes are detected.
4. Complete final current-toolchain evidence and documentation review.

## Notes

- Every task uses the required checklist format: checkbox, sequential ID, optional parallel marker, user-story label where required, and exact file paths.
- Syntax-specific fixture files document the minimum language version of the construct they exercise; Go 1.27.1 includes all selected fixtures.
- Do not edit generated protobuf files; test them through `internal/integration/protobuf/proto_test.go`.

---

## Phase 7: Convergence

- [X] T042 CRITICAL Document the Go 1.27.1-only support-boundary migration guidance and the required release or major-version treatment per Constitution II (partial).
- [X] T043 Reconcile the selected-toolchain source of truth with the CI workflow, matrix verifier, and update procedure so US3/AC2's "one documented place" policy is either implemented or accurately specified (partial).
- [X] T044 Make `UpdatePackages` atomic on conversion failure and add a regression test proving a failed generic conversion does not leave caller-visible partially initialized packages per FR-006 (partial).

---

## Phase 8: Convergence

- [X] T045 Reuse already loaded dependency packages during staged `UpdatePackages` parsing and add cross-package identity coverage so successful updates neither replace nor reference unpublished duplicate packages per FR-003, FR-004, and FR-006 (partial).

---

## Phase 9: Convergence

- [X] T046 Align the stale multi-version CI wording in the task implementation strategy with the Go 1.27.1-only current-toolchain policy per US3 and FR-008 (contradicts).
