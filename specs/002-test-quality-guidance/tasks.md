---

description: "Actionable implementation tasks for Test Quality Guidance"
---

# Tasks: Test Quality Guidance

**Input**: Design documents from `specs/002-test-quality-guidance/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [quality-artifacts.md](contracts/quality-artifacts.md), and [quickstart.md](quickstart.md)

**Tests**: Required. The feature itself establishes test quality; write each
new automated check before its implementation and verify that it fails for the
documented broken condition.

**Organization**: Tasks are grouped by user story so each delivers a separately
testable increment.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can proceed in parallel because it changes a separate file and has no incomplete prerequisite.
- **[Story]**: Maps a task to its user story. Setup, foundational, and polish tasks do not have a story label.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create safe common primitives used by consumer scenarios and
example validation.

- [X] T001 Add generated coverage and temporary validation-artifact exclusions to `.gitignore`
- [X] T002 Create tracked-source snapshot helper with cleanup and fail-closed checks in `scripts/snapshot-tracked-source.sh`
- [X] T003 Add snapshot helper tests for tracked inputs, untracked-file exclusion, and cleanup in `scripts/snapshot-tracked-source_test.sh`
- [X] T004 Add `test-quality`, `test-integration`, `test-gaps`, and `verify-docs` validation entry points in `Makefile`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish the explicit contracts and CI scaffolding required by all
three stories.

**⚠️ CRITICAL**: Finish this phase before claiming any scenario, inventory, or
published example is authoritative.

- [X] T005 Create an inventory fixture and command-contract tests in `cmd/testinventory/main_test.go`
- [X] T006 Implement deterministic JSON inventory loading, exported-symbol discovery, evidence-path validation, and Markdown rendering in `cmd/testinventory/main.go`
- [X] T007 Add the initial machine-readable behavior inventory schema and all current exported-contract entries in `docs/testing/behavior-inventory.json`
- [X] T008 Add coverage collection and `testinventory assess` orchestration in `scripts/assess-test-gaps.sh`
- [X] T009 Add named CI checks for the inventory, candidate-source scenarios, and documentation validation in `.github/workflows/ci.yml`

**Checkpoint**: The repository can classify public contracts, create isolated
candidate snapshots, and expose the three validation categories in CI.

---

## Phase 3: User Story 1 - Validate Real Consumer Workflows (Priority: P1) 🎯 MVP

**Goal**: Replace the single narrow protobuf check with independently runnable
consumer scenarios that make candidate-source versus published-module evidence
unambiguous.

**Independent Test**: Run `scripts/test-integration.sh all` from a clean
checkout. The inspect, generated-source, and cross-package scenarios must pass;
the invalid scenario must pass only by observing its expected non-zero diagnostic;
each scenario must fail when its assertion is deliberately mutated.

### Tests for User Story 1

- [X] T010 [P] [US1] Add scenario-runner tests for selection, snapshot boundary labeling, missing inputs, cleanup, and expected-failure handling in `internal/integration/runner/runner_test.go`
- [X] T011 [P] [US1] Add a public-contract consumer fixture and assertions for declaration relationships in `internal/integration/scenarios/inspect/main_test.go`
- [X] T012 [P] [US1] Add generated-protobuf fixture assertions for all fields and imported dependency types in `internal/integration/scenarios/generated/main_test.go`
- [X] T013 [P] [US1] Add a cross-package generic consumer fixture that asserts instantiation origin and ordered arguments in `internal/integration/scenarios/crosspackage/main_test.go`
- [X] T014 [P] [US1] Add an invalid-source and unresolved-dependency fixture that asserts an actionable non-zero diagnostic in `internal/integration/scenarios/invalid/main_test.go`

### Implementation for User Story 1

- [X] T015 [US1] Implement scenario selection, temporary candidate snapshotting, module replacement injection, and labeled output in `internal/integration/runner/runner.go`
- [X] T016 [US1] Create the successful source-inspection consumer module in `internal/integration/scenarios/inspect/go.mod` and `internal/integration/scenarios/inspect/main.go`
- [X] T017 [US1] Correct the protobuf `go_package` mismatch, record reproducible generator prerequisites, and implement the generated-source scenario in `internal/integration/scenarios/generated/models.proto`, `internal/integration/scenarios/generated/go.mod`, and `internal/integration/scenarios/generated/main.go`
- [X] T018 [US1] Create the consumer-owned cross-package generic source and executable scenario in `internal/integration/scenarios/crosspackage/go.mod` and `internal/integration/scenarios/crosspackage/main.go`
- [X] T019 [US1] Create the invalid source and unresolved-import scenario inputs in `internal/integration/scenarios/invalid/go.mod`, `internal/integration/scenarios/invalid/main.go`, and `internal/integration/scenarios/invalid/testdata/invalid.go`
- [X] T020 [US1] Implement all-scenario execution and deliberate assertion mutation verification in `scripts/test-integration.sh`
- [X] T021 [US1] Replace the former one-assertion protobuf test with scenario-routing coverage in `internal/integration/protobuf/proto_test.go`
- [X] T022 [US1] Document candidate-source command output, scenario meanings, and release-evidence distinction in `docs/testing/integration.md`

**Checkpoint**: A maintainer can run four isolated scenarios and identify both
the source boundary and the exact failed expectation.

---

## Phase 4: User Story 2 - Identify Unit-Test Gaps (Priority: P1)

**Goal**: Turn public-contract test coverage into a repeatable, prioritized
assessment and close the highest-risk current unit-test gaps.

**Independent Test**: Generate a coverage profile and run `testinventory
assess`; every exported contract must be classified, every evidence path must
exist, and a fixture with an unclassified symbol or stale test path must fail.

### Tests for User Story 2

- [X] T023 [P] [US2] Add built-in type lookup, conversion, equality, and panic-contract tests in `types/builtin_test.go`
- [X] T024 [P] [US2] Add array, map, pointer, and channel rendering, equality, zero-value, and element tests in `types/value_types_test.go`
- [X] T025 [P] [US2] Add package declaration, collision, lookup, must-method, and package-path tests in `types/package_test.go`
- [X] T026 [P] [US2] Add struct-tag split, lookup, escaping, and absent-key tests in `types/struct_tag_test.go`
- [X] T027 [P] [US2] Add interface compatibility and negative receiver, variadic, and signature tests in `types/interface_test.go`
- [X] T028 [P] [US2] Add alias/function rendering, equality, `Kind`, and `Dereference` tests in `types/contracts_test.go`
- [X] T029 [P] [US2] Add `PackageNameOfDir`, empty configuration, invalid path, no-op update, and partial-update error tests in `parser/load-packages_contract_test.go`

### Implementation for User Story 2

- [X] T030 [US2] Classify all new unit tests, parser fixture checks, consumer scenarios, examples, generated protobuf exclusions, and toolchain-only behavior in `docs/testing/behavior-inventory.json`
- [X] T031 [US2] Generate and commit the first reviewed risk-prioritized report in `docs/testing/behavior-inventory.md`
- [X] T032 [US2] Add assessment usage, coverage caveats, gap-remediation rules, and the supported-toolchain link in `docs/testing/unit-test-gaps.md`
- [X] T033 [US2] Add a deterministic inventory-report freshness check to `scripts/assess-test-gaps.sh`

**Checkpoint**: The assessment identifies missing work by contract and risk,
not by percentage alone, and its first high-risk gaps have focused unit tests.

---

## Phase 5: User Story 3 - Learn Code Generation from Published Examples (Priority: P2)

**Goal**: Give users a public, tested GitHub Pages learning path and a runnable
consumer example that generates deterministic Go source from Gentools analysis.

**Independent Test**: From a clean snapshot, run the example and documentation
validator. The generated file must equal its expected output, every site link
must resolve, and a changed expected output or broken link must fail validation.

### Tests for User Story 3

- [X] T034 [P] [US3] Add golden-output and invalid-input tests for the code-generation example in `examples/codegen/struct-summary/main_test.go`
- [X] T035 [P] [US3] Add static-site local-link, source-link, and example-command validation tests in `scripts/verify-docs_test.sh`

### Implementation for User Story 3

- [X] T036 [US3] Create a representative analyzed-source fixture in `examples/codegen/struct-summary/testdata/models.go`
- [X] T037 [US3] Implement the deterministic source-summary generator using the public Gentools API in `examples/codegen/struct-summary/main.go`
- [X] T038 [US3] Commit the expected generated Go output in `examples/codegen/struct-summary/testdata/zz_summary.golden.go`
- [X] T039 [US3] Create the GitHub Pages landing page with prerequisites, toolchain boundary, testing-guide links, and example routing in `docs/index.html`
- [X] T040 [US3] Create the code-generation guide with input, command, expected output, inspection result, and troubleshooting in `docs/code-generation.html`
- [X] T041 [US3] Implement clean-snapshot example execution, golden comparison, static artifact preparation, and link checks in `scripts/verify-docs.sh`
- [X] T042 [US3] Add a Pages build-and-deploy workflow using the validated static artifact and least-privilege permissions in `.github/workflows/pages.yml`
- [X] T043 [US3] Add GitHub Actions publishing-source setup and deployment URL verification guidance in `docs/publishing.md`

**Checkpoint**: The site routes a new user to a verified generation example,
and deployment cannot publish an artifact with a stale example or broken link.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Prove release behavior, validate all acceptance criteria, and make
the project entry points discoverable.

- [X] T044 Add a tag-triggered no-replacement published-module consumer smoke workflow in `.github/workflows/release-smoke.yml`
- [X] T045 Update development, validation, integration, test-gap, example, and GitHub Pages entry points in `README.md`
- [X] T046 Run every command in `specs/002-test-quality-guidance/quickstart.md` and record sanitized results in `specs/002-test-quality-guidance/quickstart.md`
- [X] T047 Run formatting, vetting, root tests, nested consumer scenarios, inventory assessment, documentation validation, and workflow syntax checks from `Makefile`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts immediately. T003 depends on T002; T004 depends on T002 and T003.
- **Foundational (Phase 2)**: T005 through T009 depend on T004. T006 depends on T005; T007 depends on T006; T008 depends on T006 and T007; T009 depends on T008.
- **US1 (Phase 3)**: T010 through T014 can begin after T009. T015 depends on T010; T016 through T019 depend on their corresponding test task and T015; T020 depends on T016 through T019; T021 and T022 depend on T020.
- **US2 (Phase 4)**: T023 through T029 can begin after T009. T030 depends on T023 through T029 and US1's evidence tasks T016 through T022; T031 through T033 depend on T030.
- **US3 (Phase 5)**: T034 and T035 can begin after T009. T036 depends on T034; T037 depends on T034 and T036; T038 depends on T037; T039 and T040 depend on T035; T041 depends on T037 through T040; T042 and T043 depend on T041.
- **Polish (Phase 6)**: T044 depends on T020; T045 depends on T022, T032, and T043; T046 and T047 depend on all desired story tasks.

### User Story Dependencies

- **US1 (P1)**: Independent once the snapshot helper exists; this is the MVP.
- **US2 (P1)**: Its command scaffolding is foundational, but its behavior inventory can be completed after the US1 evidence exists.
- **US3 (P2)**: Independent of US1 and US2 implementation, apart from the shared snapshot helper; it links to their published guidance once available.

### Parallel Opportunities

- T010 through T014 are separate scenario tests and may proceed in parallel.
- T023 through T029 are separate test files and may proceed in parallel.
- T034 and T035 may proceed in parallel; T039 and T040 then modify different Pages files and may proceed in parallel.
- After T009, a developer may start US1 scenario tests while another starts US3 example tests; a third may add US2 unit tests.

## Parallel Example: User Story 1

```text
Task: "Add public-contract consumer assertions in internal/integration/scenarios/inspect/main_test.go"
Task: "Add generated protobuf assertions in internal/integration/scenarios/generated/main_test.go"
Task: "Add cross-package generic assertions in internal/integration/scenarios/crosspackage/main_test.go"
Task: "Add invalid-input diagnostic assertions in internal/integration/scenarios/invalid/main_test.go"
```

## Implementation Strategy

### MVP First (User Story 1)

1. Complete T001 through T009.
2. Complete T010 through T022.
3. Run the US1 independent test and mutation checks.
4. Review the source-boundary labels before treating scenario results as release evidence.

### Incremental Delivery

1. Deliver the consumer scenario suite and CI visibility.
2. Add the inventory command, baseline classification, and focused high-risk unit tests.
3. Add the runnable generator example and validate its static Pages artifact.
4. Enable the release-only smoke test and publish the documentation after all checks pass.

### Format Validation

All 47 tasks use the required checkbox, sequential task ID, optional parallel
marker only where independent, story label only for story work, and at least one
exact target path.

---

## Phase 7: Convergence

- [X] T048 Record and validate the approved reviewed package-level public API groups (`parser.*` and `types.*`) in `docs/testing/behavior-inventory.json` per FR-001 and SC-001
- [X] T049 Make `internal/integration/runner/runner.go` execute manifest-selected scenarios with candidate snapshots and add assertion-mutation coverage in `internal/integration/runner/runner_test.go` per FR-004 and US1/AC1 (partial)
- [X] T050 Run a documented protobuf generator preflight, regenerate `internal/integration/protobuf/models.pb.go` from `internal/integration/protobuf/models.proto`, and verify reproducibility without manual generated-code edits per Constitution additional constraints (partial)
- [X] T051 Complete the focused unit-test matrix in `types/builtin_test.go`, `types/value_types_test.go`, `types/package_test.go`, `types/struct_tag_test.go`, and `types/interface_test.go` per FR-003 and SC-001 (missing)
- [X] T052 Add a committed-snapshot integration test and document the working-tree versus committed-candidate boundary in `scripts/snapshot-tracked-source_test.sh` and `docs/testing/integration.md` per FR-009 and SC-004 (partial)
- [X] T053 Record the gap-assessment elapsed time and an independent documentation walkthrough in `specs/002-test-quality-guidance/quickstart.md` per SC-002 and SC-005 (missing)
- [X] T054 Have a repository administrator enable GitHub Actions as the Pages source and record the first verified deployment URL in `docs/publishing.md` per FR-007 (partial)
