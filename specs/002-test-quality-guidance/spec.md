# Feature Specification: Test Quality Guidance

**Feature Branch**: `002-test-quality-guidance`

**Created**: 2026-09-02

**Status**: Draft

**Input**: User description: "Improve the project's integration testing beyond the current approach; identify missing unit-test coverage; and provide high-quality, runnable code-generation examples, ideally accompanied by GitHub Pages documentation that explains how to use them."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Validate Real Consumer Workflows (Priority: P1)

As a library maintainer, I want representative integration checks that exercise the library as an independently installed consumer would, so that releases catch packaging, dependency, generated-source, and cross-package regressions that isolated tests miss.

**Why this priority**: The current integration coverage is a single narrow generated-model scenario and does not demonstrate that a consumer can reliably use the supported library workflows.

**Independent Test**: Execute the documented integration suite from a clean consumer fixture and verify that every supported end-to-end scenario has a visible result and that deliberately broken fixture expectations cause the relevant check to fail.

**Acceptance Scenarios**:

1. **Given** a clean representative consumer fixture, **When** it loads source using the published library contract, **Then** it can inspect the documented declarations and relationships with the expected results.
2. **Given** generated source that uses the supported code-generation workflow, **When** the fixture is validated, **Then** the generated declarations and their imported dependencies can be inspected without relying on uncommitted local state.
3. **Given** an invalid consumer input or unresolved dependency, **When** the integration suite runs, **Then** it fails with an actionable diagnostic rather than reporting a partial success.

---

### User Story 2 - Identify Unit-Test Gaps (Priority: P1)

As a maintainer, I want a repeatable assessment of untested library behavior, so I can prioritize missing unit tests by public-contract risk instead of relying on an opaque percentage alone.

**Why this priority**: Adding integration checks cannot establish whether individual public behaviors, invalid inputs, and error paths are adequately protected.

**Independent Test**: Run the documented assessment on the repository and compare its results with a reviewed behavior inventory; every reported gap includes the affected behavior, risk category, and recommended test level.

**Acceptance Scenarios**:

1. **Given** the current public and internal behavior inventory, **When** the assessment is run, **Then** it identifies covered behaviors, unverified behaviors, and explicitly excluded generated or unreachable code.
2. **Given** a behavior with no meaningful unit-level assertion, **When** it is reviewed, **Then** the assessment records why it is a gap and whether a unit, integration, compatibility, or documentation example check is appropriate.
3. **Given** a newly added behavior, **When** the assessment is repeated, **Then** its test status is updated without manually rewriting the whole report.

---

### User Story 3 - Learn Code Generation from Published Examples (Priority: P2)

As a library user, I want a public GitHub Pages guide with runnable, versioned code-generation examples, so I can choose the correct starting point and reproduce the workflow in my own project.

**Why this priority**: Current documentation shows library calls but does not provide a discoverable learning path for the requested code-generation use cases.

**Independent Test**: Starting from the published guide, follow each example in a fresh checkout and verify that its stated input, command, generated output, and inspection result agree with the example's expected outcome.

**Acceptance Scenarios**:

1. **Given** a prospective user visiting the published documentation, **When** they choose a code-generation use case, **Then** they can find prerequisites, input source, expected output, and a copyable run path in one place.
2. **Given** an example's source input changes, **When** the project validates its documentation, **Then** stale generated output or expected results are detected before publication.
3. **Given** a user encounters a toolchain, dependency, or generation failure, **When** they consult the guide, **Then** they can distinguish expected setup requirements from an unsupported input and find the next diagnostic step.

### Edge Cases

- A fixture passes only because it reads an unpublished local replacement rather than the intended released-library contract.
- Generated source has an invalid package path, stale output, an unavailable generator, or imported generated dependencies.
- A behavior is intentionally not unit tested because it can only be verified at integration level; the assessment must record the rationale rather than label it missing by default.
- Coverage changes without a corresponding behavior-level test, producing a misleading percentage increase.
- A published example is incompatible with the selected toolchain or cannot run from a clean checkout.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The project MUST define and maintain a behavior inventory that maps supported library contracts to their verification status and appropriate test level. Public contracts MAY be represented by a reviewed package-level entry when that entry names the package, explains the grouping, and cites representative evidence.
- **FR-002**: The project MUST provide a repeatable unit-test-gap assessment that reports unverified behaviors, their risk priority, evidence source, and a recommended verification level.
- **FR-003**: The assessment MUST distinguish intentionally excluded generated code, unreachable code, and integration-only behavior from missing unit coverage.
- **FR-004**: The project MUST provide integration scenarios that run as representative consumers and cover a successful source-inspection workflow, generated-source workflow, cross-package dependency workflow, and an expected-failure workflow.
- **FR-005**: Each integration scenario MUST declare its source inputs, expected observable result, isolation assumptions, and failure diagnostic expectation.
- **FR-006**: The standard project validation workflow MUST run the integration scenarios independently and make their results identifiable to maintainers.
- **FR-007**: The project MUST publish a GitHub Pages documentation site that introduces the supported code-generation use cases and routes users to runnable examples.
- **FR-008**: Each published code-generation example MUST include its purpose, prerequisites, source input, execution steps, expected generated result, expected inspection result, and troubleshooting guidance.
- **FR-009**: The project MUST validate each published example from a clean checkout and detect divergence between its documented expected result and the actual result.
- **FR-010**: Documentation and assessment outputs MUST state their supported-toolchain boundary and link to the authoritative compatibility policy.

### Key Entities *(include if feature involves data)*

- **Behavior inventory entry**: A supported contract or meaningful behavior, its owner, risk level, verification status, evidence, and recommended test level.
- **Integration scenario**: A representative consumer workflow with its inputs, expected observable result, setup assumptions, and diagnostic expectation.
- **Runnable example**: A versioned learning artifact that connects a code-generation use case to source input, commands, expected output, and troubleshooting.
- **Published guide**: The navigable GitHub Pages content that explains how users select, run, and validate runnable examples.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The behavior inventory accounts for 100% of the public `parser` and `types` package surfaces and all identified high-risk internal behaviors, with each entry classified as verified, excluded with rationale, or requiring work. A reviewed package-level entry may account for an entire public package surface.
- **SC-002**: The gap assessment produces a prioritized, reviewable result from a clean checkout in under 10 minutes and identifies the verification level for 100% of reported gaps.
- **SC-003**: At least four independently runnable integration scenarios cover the success, generated-source, cross-package, and expected-failure workflows; each fails when its expected result is deliberately changed.
- **SC-004**: Every runnable published example completes successfully from a clean checkout using the supported toolchain, and its documented expected result matches the actual result.
- **SC-005**: A maintainer unfamiliar with the examples can select and complete the relevant documented code-generation workflow without oral guidance, with every prerequisite and diagnostic path available from the published guide.

## Assumptions

- GitHub Pages is the requested public documentation channel; publication configuration and visual design are implementation decisions for planning.
- The existing selected toolchain remains the authority for examples and validation unless a separately approved compatibility change updates it.
- The feature improves and may replace the current nested generated-model integration coverage rather than preserving its exact structure.
- Code generation means supported workflows that consume or inspect generated Go source; defining a new public code generator is outside this feature unless planning discovers an already-supported contract that lacks documentation.
- The behavior inventory will prioritize public compatibility and failure behavior; raw line coverage is supporting evidence, not the definition of completeness.
