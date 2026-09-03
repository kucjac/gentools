---
description: "Dependency-ordered implementation tasks for OpenAPI contract generator"
---

# Tasks: OpenAPI Contract Generator

**Input**: Design documents from `/specs/003-openapi-contract-generator/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, and `quickstart.md`

**Organization**: Tasks are grouped by user story; tests are included because the specification explicitly requires automated checks.

## Phase 1: Setup

**Purpose**: Establish the example and fixture layout.

- [X] T001 Create the `examples/codegen/openapi-contract/` command and `testdata/` fixture paths from `specs/003-openapi-contract-generator/plan.md`
- [X] T002 [P] Add the annotated model and handler fixture declarations in `examples/codegen/openapi-contract/testdata/models.go` and `examples/codegen/openapi-contract/testdata/handlers.go`
- [X] T003 [P] Add the initial documentation-validation contract and quickstart references in `specs/003-openapi-contract-generator/contracts/validation-contract.md` and `specs/003-openapi-contract-generator/quickstart.md`

## Phase 2: Foundational

**Purpose**: Define the internal contract representation, annotation parser, and transactional output boundary required by all stories.

- [X] T004 Define OpenAPI document, path, operation, response, schema, and diagnostic types in `examples/codegen/openapi-contract/main.go`
- [X] T005 [P] Implement strict `@openapi` directive tokenization and validation for method, route, summary, response, and unknown/duplicate keys in `examples/codegen/openapi-contract/main.go`
- [X] T006 [P] Implement atomic output writing that preserves or omits the destination when validation fails in `examples/codegen/openapi-contract/main.go`

## Phase 3: User Story 1 - Produce an API contract (Priority: P1) 🎯 MVP

**Goal**: Generate a deterministic OpenAPI 3 contract with operations, parameters, success responses, and shared model definitions.

**Independent Test**: Run the fixture command and compare its JSON with `examples/codegen/openapi-contract/testdata/openapi.golden.json`; verify both operations reference one shared model definition and field tags/descriptions are present.

### Tests for User Story 1

- [X] T007 [P] [US1] Add successful contract and shared-model reuse assertions in `examples/codegen/openapi-contract/main_test.go`
- [X] T008 [P] [US1] Add representative annotated models, handlers, and reviewed output in `examples/codegen/openapi-contract/testdata/models.go`, `examples/codegen/openapi-contract/testdata/handlers.go`, and `examples/codegen/openapi-contract/testdata/openapi.golden.json`

### Implementation for User Story 1

- [X] T009 [US1] Load source packages with comments, discover annotated exported functions, and build validated operations in `examples/codegen/openapi-contract/main.go`
- [X] T010 [US1] Map supported Gentools primitive, pointer, slice, and struct metadata to OpenAPI schemas using JSON tags and field comments in `examples/codegen/openapi-contract/main.go`
- [X] T011 [US1] Emit deterministic OpenAPI 3 JSON with path parameters, responses, and deduplicated `components.schemas` in `examples/codegen/openapi-contract/main.go`

## Phase 4: User Story 2 - Receive actionable contract diagnostics (Priority: P2)

**Goal**: Reject invalid annotations and model references without producing a replacement artifact.

**Independent Test**: Each invalid fixture fails with an operation/model-field diagnostic and leaves a pre-existing or absent output unchanged.

### Tests for User Story 2

- [X] T012 [P] [US2] Add malformed, duplicate-route, unresolved-model, unsupported-shape, and output-safety tests in `examples/codegen/openapi-contract/main_test.go`

### Implementation for User Story 2

- [X] T013 [US2] Add duplicate method-route, unresolved model, malformed route/method, and missing declaration diagnostics in `examples/codegen/openapi-contract/main.go`
- [X] T014 [US2] Add unsupported model-shape diagnostics for maps, interfaces, channels, functions, and unnamed fields in `examples/codegen/openapi-contract/main.go`
- [X] T015 [US2] Ensure generation performs all validation before temporary-file creation and atomic rename in `examples/codegen/openapi-contract/main.go`

## Phase 5: User Story 3 - Learn from a publishable example (Priority: P3)

**Goal**: Make the annotation convention, command, expected output, supported scope, and failure path discoverable and reproducible.

**Independent Test**: From the guide, a user can run the command and obtain the reviewed output; docs validation catches a changed result.

- [X] T016 [US3] Add OpenAPI generator command, annotation convention, supported scope, and diagnostics guidance to `README.md`
- [X] T017 [P] [US3] Add the generator walkthrough and expected contract excerpt to `docs/code-generation.html`
- [X] T018 [US3] Extend `scripts/verify-docs.sh` to run the generator from the committed tracked-source snapshot and compare normalized output with `openapi.golden.json`
- [X] T019 [US3] Extend `scripts/verify-docs_test.sh` to assert the OpenAPI example validation path and guide links remain present

## Phase 6: Polish and cross-cutting validation

- [X] T020 [P] Run `gofmt`, `git diff --check`, and the example package tests; record final evidence in `specs/003-openapi-contract-generator/quickstart.md`
- [X] T021 [P] Run `make build-and-test` and `make test-integration` to confirm no library or integration regression
- [X] T022 Run `make verify-docs` and the complete quickstart command from an isolated tracked-source snapshot; record results in `specs/003-openapi-contract-generator/quickstart.md`
- [X] T023 Review the implementation against every FR/SC in `specs/003-openapi-contract-generator/spec.md`, then mark all completed tasks `[X]`

## Dependencies and execution order

- Setup precedes Foundational; Foundational blocks all user stories.
- US1 is the MVP and must precede the docs comparison in US3.
- US2 depends on the operation/model builder from US1 but adds independent failure coverage.
- US3 depends on the stable command/output from US1 and diagnostics from US2.
- Polish depends on the desired user stories being complete.

Parallel opportunities: T002-T003, T005-T006, T007-T008, and T012 can proceed in parallel when their file ownership is coordinated. T017 and the implementation work can be developed separately until the final docs integration.

## Implementation strategy

Implement Setup and Foundational, deliver US1 as the MVP, then add US2 failure safety and US3 publication guidance. Finish with isolated-snapshot validation and a full regression pass.

## Phase 7: Convergence

Remaining work identified by code-aware comparison of the implementation with
`spec.md`, `plan.md`, and the existing tasks.

- [X] T024 [US2] Add invalid fixtures and assertions for unknown directives, malformed routes, each missing required annotation key, and absent-output safety in `examples/codegen/openapi-contract/main_test.go` and `examples/codegen/openapi-contract/testdata/invalid/` per FR-008 and SC-003 (partial)
- [X] T025 [US1] Add nested, optional, collection, and cross-package model fixtures and assert the documented serialized-name, description, reference, and optional-field behavior in `examples/codegen/openapi-contract/main_test.go` and `examples/codegen/openapi-contract/testdata/` per US1/AC2 and the documented edge cases (partial)
- [X] T026 [US1] Resolve model definitions by stable package-qualified identity and add a collision/cross-package shared-model test in `examples/codegen/openapi-contract/main.go` and `examples/codegen/openapi-contract/main_test.go` per FR-003, FR-005, and the cross-package edge case (partial)
- [X] T027 [US2] Verify and, if necessary, implement portable atomic replacement of an existing output on the repository’s Windows CI target in `examples/codegen/openapi-contract/main.go` and `examples/codegen/openapi-contract/main_test.go` per FR-007 and the plan portability constraint (partial)
- [X] T028 [US3] Record measured command duration for SC-001 and run the published example from an isolated tracked-source snapshot for SC-004 in `specs/003-openapi-contract-generator/quickstart.md` after the feature files are tracked (partial)
