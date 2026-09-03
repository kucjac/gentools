# Feature Specification: CI and Lint Modernization

**Feature Branch**: `004-ci-lint-modernization`
**Created**: 2026-09-03
**Status**: Draft
**Input**: User description: "I want to update the CI and golangci-lint to satisfy best practices, check the latest version of golangci-lint to match Go v1.27, and use it in the pipeline because it fails remotely."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Enforce a reproducible remote lint gate (Priority: P1)

As a maintainer, I want the hosted CI lint job to use an explicitly selected linter release that supports the project's supported Go version so pull requests receive a dependable pass or fail result.

**Why this priority**: The current remote lint job can fail independently of local validation, preventing reliable contribution and release decisions.

**Independent Test**: Trigger the lint job for a representative pull request and confirm it completes with a result that can be reproduced using the repository's documented local validation path.

**Acceptance Scenarios**:

1. **Given** a change that meets the repository's lint policy, **When** the hosted lint job runs, **Then** it completes successfully using the declared linter release and supported Go version.
2. **Given** a change that violates the lint policy, **When** the hosted lint job runs, **Then** it fails and reports actionable findings.
3. **Given** a supported operating-system test run, **When** CI validates the project, **Then** linting remains isolated from unrelated platform-specific test failures.

---

### User Story 2 - Maintain a current and intentional lint policy (Priority: P2)

As a maintainer, I want the repository's lint policy to be compatible with the selected current linter release so checks reflect intentional code-quality rules rather than obsolete configuration behavior.

**Why this priority**: A valid, reviewable policy avoids silent rule loss or unexpected behavior during linter upgrades.

**Independent Test**: Run the repository's documented lint command against the full root module and confirm the effective policy is accepted by the selected linter release.

**Acceptance Scenarios**:

1. **Given** the repository lint configuration, **When** it is loaded by the selected linter release, **Then** it is accepted without deprecated, ignored, or invalid settings.
2. **Given** lint rules intentionally enabled for the project, **When** linting is run locally and remotely, **Then** both environments apply the same policy to the same source scope.

---

### User Story 3 - Give contributors clear validation guidance (Priority: P3)

As a contributor, I want to know the required Go and linter versions and the canonical validation command so I can resolve failures before submitting changes.

**Why this priority**: Clear, portable instructions reduce avoidable remote-only failures and support contributors on different systems.

**Independent Test**: A contributor following repository documentation can identify the required tools and execute the same lint validation used by CI without relying on a host-specific path or hidden state.

**Acceptance Scenarios**:

1. **Given** a contributor with a compatible Go installation on their PATH, **When** they follow the repository guidance, **Then** they can run the canonical lint validation command.
2. **Given** an incompatible Go or linter version, **When** the contributor runs the validation path, **Then** they receive an actionable mismatch or setup failure.

### Edge Cases

- The selected current linter release changes its configuration format or removes a rule used by the repository.
- The linter release does not support the project's declared Go 1.27.x toolchain.
- A hosted runner obtains a different linter release than local development due to an unpinned reference or stale cache.
- The root module passes linting while a nested integration module requires separate validation scope.
- Lint output contains findings but the job is configured to report success.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The repository MUST declare one exact golangci-lint release for CI and contributor validation.
- **FR-002**: The declared linter release MUST support the Go 1.27.x version declared by the repository.
- **FR-003**: The hosted lint job MUST use the declared linter release rather than an unbounded or implicit latest release.
- **FR-004**: The hosted lint job MUST fail when the configured lint policy finds unresolved violations.
- **FR-005**: The lint policy MUST be accepted by the declared linter release without invalid, deprecated, or silently ignored configuration.
- **FR-006**: Local and hosted lint validation MUST use the same declared linter release, policy, and root-module source scope.
- **FR-007**: The repository MUST provide portable contributor guidance for installing or selecting the required Go and linter releases without committing host-specific SDK locations or credentials.
- **FR-008**: The CI workflow MUST preserve the existing cross-platform build-and-test coverage while making lint failures independently diagnosable.
- **FR-009**: The change MUST document the validation evidence needed to confirm the lint gate works in the hosted pipeline and locally.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of hosted lint runs use the single declared linter release, as shown in the job log.
- **SC-002**: A clean root-module revision completes the hosted lint job successfully on the first run with the supported Go 1.27.x toolchain.
- **SC-003**: A deliberate lint-policy violation causes the hosted lint job to fail and identifies at least one actionable location and rule.
- **SC-004**: A contributor can identify the required tool versions and run the canonical local lint validation using only repository documentation and PATH-visible tools.
- **SC-005**: Existing build-and-test jobs continue to cover all three supported CI operating systems after the lint-gate update.

## Assumptions

- The Go 1.27.1 version in `.go-version` and the current CI matrix remains the supported project baseline for this feature.
- The current upstream golangci-lint release at specification time is v2.12.2; implementation will re-verify the release and its Go 1.27 support before pinning it.
- The feature covers the root-module lint gate and its contributor guidance; adding lint execution for nested integration modules is out of scope unless planning finds it is already a stated project requirement.
- Existing public APIs and their observable behavior are out of scope for this quality-tooling change.
