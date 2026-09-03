---
description: "Dependency-ordered implementation tasks for CI and lint modernization"
---

# Tasks: CI and Lint Modernization

**Input**: Design documents from `/specs/004-ci-lint-modernization/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, and `quickstart.md`

**Organization**: Tasks are grouped by user story. The feature changes repository automation and documentation only; no runtime Go model or public API task is required.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish the exact tool and baseline evidence needed by all stories.

- [X] T001 Confirm the active feature directory and record the existing Go baseline from `.go-version`, `go.mod`, and `.github/workflows/ci.yml`
- [X] T002 [P] Install or select golangci-lint `v2.13.2` from the official binary installer using a PATH-visible location, and record `golangci-lint version` output in `specs/004-ci-lint-modernization/research.md`
- [X] T003 [P] Capture the current root lint configuration and existing local lint target from `.golangci.yml` and `Makefile` for migration comparison

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Define the one version, source scope, and validation contract that every story must preserve.

**Checkpoint**: The version/scope decisions are fixed before editing CI, policy, or contributor guidance.

- [X] T004 Define the single shared linter version constant/value as `v2.13.2` in the implementation approach documented by `specs/004-ci-lint-modernization/research.md`
- [X] T005 [P] Verify the selected binary accepts Go 1.27 source/toolchain behavior and document the result in `specs/004-ci-lint-modernization/research.md`
- [X] T006 [P] Review `specs/004-ci-lint-modernization/contracts/validation-contract.md` and `specs/004-ci-lint-modernization/quickstart.md` for exact command, scope, exit-status, and evidence requirements

---

## Phase 3: User Story 1 - Enforce a reproducible remote lint gate (Priority: P1) 🎯 MVP

**Goal**: Make hosted lint deterministic, independently diagnosable, and failing on unresolved findings while preserving the three-OS build/test matrix.

**Independent Test**: On a clean revision, the hosted lint job uses v2.13.2 and root scope `./...` and passes; a deliberate enabled-rule violation produces actionable output and a non-zero job result; the three matrix operating systems remain present.

### Implementation for User Story 1

- [X] T007 [US1] Update the `golangci-lint` job in `.github/workflows/ci.yml` to select exactly `v2.13.2` through `golangci/golangci-lint-action@v9`
- [X] T008 [US1] Configure the hosted lint action in `.github/workflows/ci.yml` to use `.golangci.yml`, execute root scope `./...`, and preserve actionable annotations and non-zero failure behavior
- [X] T009 [US1] Verify `.github/workflows/ci.yml` still contains Ubuntu, macOS, and Windows entries for the independent `build-and-test` matrix and keeps lint as a separate job
- [X] T010 [US1] Run the hosted-workflow static checks and a local equivalent of the clean root lint command, recording tool version, scope, result, and any configuration diagnostics in `specs/004-ci-lint-modernization/quickstart.md`
- [X] T011 [US1] Exercise the deliberate lint-violation scenario described in `specs/004-ci-lint-modernization/quickstart.md` and record the reported rule/location and non-zero result

**Checkpoint**: User Story 1 is independently deliverable once the workflow diff is reviewed and clean/failing lint evidence is available.

---

## Phase 4: User Story 2 - Maintain a current and intentional lint policy (Priority: P2)

**Goal**: Convert the legacy configuration to the selected v2 schema while retaining only intentional rules and ensuring invalid or deprecated settings cannot silently pass.

**Independent Test**: `golangci-lint config verify --config .golangci.yml` succeeds with v2.13.2, and local lint applies the same enabled rules and exclusions as CI without an exit-code success override.

### Implementation for User Story 2

- [X] T012 [P] [US2] Run `golangci-lint migrate --config .golangci.yml` with v2.13.2 as a migration reference and compare the result against the v2 configuration schema
- [X] T013 [US2] Rewrite `.golangci.yml` with an explicit v2 configuration version and v2-compatible run/output/linter/issues keys while preserving intentional enabled linters and exclusions
- [X] T014 [US2] Remove or replace legacy settings in `.golangci.yml` that v2.13.2 rejects, including obsolete `linters-settings` fields and the `issue-exit-code: 0` success override
- [X] T015 [US2] Update the `golangci-lint` target in `Makefile` to require the exact linter version, use the same `.golangci.yml` and `./...` scope as CI, and propagate non-zero findings
- [X] T016 [US2] Validate `.golangci.yml` with `golangci-lint config verify --config .golangci.yml` and run the clean root lint command with v2.13.2, resolving all policy/configuration errors
- [X] T017 [US2] Compare local and CI lint commands in `.github/workflows/ci.yml`, `Makefile`, and `specs/004-ci-lint-modernization/contracts/validation-contract.md` to confirm one policy and source scope

**Checkpoint**: User Story 2 is independently deliverable when config verification and clean local lint pass with v2.13.2 and the policy has no obsolete settings.

---

## Phase 5: User Story 3 - Give contributors clear validation guidance (Priority: P3)

**Goal**: Let contributors identify and select the exact Go/linter versions and reproduce the CI lint gate using PATH-visible tools.

**Independent Test**: A contributor following repository documentation can discover Go 1.27.1, golangci-lint v2.13.2, the installation path, the canonical lint command, and actionable mismatch/failure behavior without host-specific paths or secrets.

### Implementation for User Story 3

- [X] T018 [US3] Update the development and validation sections of `README.md` to name Go 1.27.1, golangci-lint v2.13.2, PATH-based setup, and `make golangci-lint` as the canonical local gate
- [X] T019 [P] [US3] Update `specs/004-ci-lint-modernization/quickstart.md` with exact version checks, config verification, clean lint, deliberate-failure, and preserved-matrix validation commands
- [X] T020 [P] [US3] Update `specs/004-ci-lint-modernization/contracts/validation-contract.md` if needed so local and hosted command/version/exit-status requirements match the implemented files
- [X] T021 [US3] Verify documentation contains no host-specific SDK locations, credentials, unbounded `latest` linter references, or contradictory lint scope/version guidance

**Checkpoint**: User Story 3 is independently deliverable when a fresh contributor can follow README/quickstart guidance and reproduce the local gate.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Run the complete validation matrix, review evidence, and prepare the feature for task analysis/implementation review.

- [X] T022 [P] Run `git diff --check` and repository documentation checks for `README.md` and `specs/004-ci-lint-modernization/`
- [X] T023 [P] Run `make build-and-test` with Go 1.27.1 and confirm no unrelated root-module regression
- [X] T024 [P] Run `make test-integration` and `make verify-docs` to confirm nested integration and Pages validation remain intact
- [X] T025 Run the complete `specs/004-ci-lint-modernization/quickstart.md` validation sequence and record final evidence for each success criterion
- [X] T026 Review the final diff against `specs/004-ci-lint-modernization/spec.md`, `plan.md`, and `contracts/validation-contract.md`, then resolve any unmet FR/SC before implementation review

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No implementation dependency; establishes the checkout baseline.
- **Foundational (Phase 2)**: Depends on Setup and blocks story work until version, scope, and contracts are confirmed.
- **User Story 1 (Phase 3)**: Depends on T004-T006; MVP for the hosted gate.
- **User Story 2 (Phase 4)**: Depends on T004-T006 and should complete before final CI evidence because CI consumes the policy.
- **User Story 3 (Phase 5)**: Depends on the final command/version decisions from US1/US2, but documentation edits can begin in parallel after T006.
- **Polish (Phase 6)**: Depends on all desired story tasks and includes the final cross-cutting validation.

### User Story Dependencies

- **US1 (P1)**: Depends only on foundational decisions; independently testable.
- **US2 (P2)**: Depends only on foundational decisions; its final local/CI comparison integrates with US1's workflow files.
- **US3 (P3)**: Depends on the implemented command/version details from US1/US2 for final accuracy; independently testable as documentation and command discovery.

### Parallel Opportunities

- T002 and T003 can run in parallel after setup begins.
- T005 and T006 can run in parallel after the version decision is recorded.
- T012 can run in parallel with initial workflow review; T019 and T020 can run in parallel with README work after command decisions are stable.
- T022-T024 can run in parallel after implementation, provided they use separate validation processes/files.
- US1 and the initial US2 migration review can be assigned to separate workers after Phase 2, but edits to `.github/workflows/ci.yml` and `.golangci.yml` should be coordinated before final comparison.

## Parallel Example: User Story 1

```text
Task T007: update the pinned action version in .github/workflows/ci.yml
Task T009: verify the existing three-OS matrix and separate job boundary in .github/workflows/ci.yml
```

These tasks inspect/edit the same workflow and should be merged carefully; T008 and T010 follow the resulting workflow shape.

## Parallel Example: User Story 2

```text
Task T012: produce migration reference output for .golangci.yml
Task T015: design the Makefile version/scope invocation against the validation contract
```

T013-T014 must consume the migration review before the final configuration is written; T016-T017 then validate parity.

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Setup and Foundational phases.
2. Pin the hosted action and preserve the independent build/test matrix.
3. Run clean and deliberate-failure lint checks.
4. Stop for review once the hosted lint gate is reproducible and independently diagnosable.

### Incremental Delivery

1. Add US1 for the deterministic hosted gate.
2. Add US2 for schema migration and local/CI policy parity.
3. Add US3 for contributor guidance.
4. Run Polish validation and review all requirements and success criteria.

### Notes

- Every task uses the required checkbox, sequential ID, optional `[P]` marker, story label where applicable, and an exact repository file path or command target.
- No task changes production Go code or public APIs.
- The nested integration module remains a separate test scope, as documented in the plan and validation contract.
