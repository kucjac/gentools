# Feature Specification: OpenAPI Contract Generator

**Feature Branch**: `003-openapi-contract-generator`

**Created**: 2026-09-03

**Status**: Draft

**Input**: User description: "OpenAPI generator request"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Produce an API contract (Priority: P1)

As an API maintainer, I want to describe public endpoints alongside their
handlers and data models, then produce a single API contract document, so that
API consumers can discover the available operations and message shapes without
manually maintaining a second source of truth.

**Why this priority**: A reliable contract generated from the maintained source
removes the most costly documentation drift and is useful on its own.

**Independent Test**: A maintainer can run the documented example against a
small service fixture and compare the resulting contract document with its
reviewed expected result.

**Acceptance Scenarios**:

1. **Given** endpoint handlers with complete operation annotations and public
   request or response models, **When** a maintainer generates the contract,
   **Then** it includes the declared operations, routes, summaries, success
   responses, and reusable data-model definitions.
2. **Given** a response model with named fields, descriptions, and serialized
   field names, **When** the contract is generated, **Then** its model
   definition exposes those names and descriptions consistently.
3. **Given** two operations that refer to the same public model, **When** the
   contract is generated, **Then** both operations refer to one shared model
   definition rather than conflicting copies.

---

### User Story 2 - Receive actionable contract diagnostics (Priority: P2)

As an API maintainer, I want invalid or incomplete endpoint annotations to stop
generation with clear diagnostics, so that I can correct source contracts
before publishing an inaccurate API description.

**Why this priority**: A plausible but wrong API document is more harmful than
no document; failures must identify the affected source declaration.

**Independent Test**: Run the generator against fixtures containing each
documented invalid condition and verify that it fails with the expected
actionable diagnostic.

**Acceptance Scenarios**:

1. **Given** an operation annotation that lacks a required route, method, or
   response declaration, **When** generation is requested, **Then** no output
   is produced and the diagnostic identifies the operation and missing item.
2. **Given** two operations that claim the same method and route, **When**
   generation is requested, **Then** it fails and identifies both conflicting
   operations.
3. **Given** an operation that names a model that cannot be resolved, **When**
   generation is requested, **Then** it fails and identifies the unresolved
   model name and operation.

---

### User Story 3 - Learn from a publishable example (Priority: P3)

As a library user, I want a concise guide and runnable example for contract
generation, so that I can adapt the annotation convention to my own API and
verify the expected result before adopting it.

**Why this priority**: The feature is most valuable when users can reproduce
the workflow rather than infer it from implementation source.

**Independent Test**: Starting from the guide, a user can run the example from
a clean source snapshot and obtain the documented contract without additional
instructions.

**Acceptance Scenarios**:

1. **Given** a user viewing the published guide, **When** they follow the
   supplied command, **Then** they can locate the annotated input and generate
   the documented contract.
2. **Given** the example input or expected contract changes, **When** project
   documentation validation runs, **Then** an outdated example result is
   detected before publication.

### Edge Cases

- An annotation uses an unsupported HTTP method, malformed route, or an
  unknown directive.
- A model contains nested, optional, collection, or cross-package fields.
- A serialized field name is omitted or explicitly excluded from the public
  representation.
- An endpoint has more than one documented success response or a non-success
  error response.
- A generated contract would overwrite an existing file after validation has
  failed.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The project MUST provide a runnable example that derives an
  OpenAPI 3 contract document from endpoint documentation annotations and
  declared public data models.
- **FR-002**: The example MUST define a documented, minimal annotation
  convention for an operation's HTTP method, route, summary, and response
  status/model pairing.
- **FR-003**: The generated contract MUST contain a document identity,
  declared operations, route parameters, success responses, and reusable
  model definitions for every model referenced by a supported operation.
- **FR-004**: Model definitions MUST use documented serialized field names,
  field descriptions, and supported field shapes from the annotated source.
- **FR-005**: The generator MUST reuse a single model definition when the
  same supported model is referenced by multiple operations.
- **FR-006**: The generator MUST reject malformed annotations, duplicate
  method-and-route pairs, unresolved referenced models, and unsupported model
  shapes with diagnostics that name the relevant source declaration.
- **FR-007**: The generator MUST not leave a new or partially replaced output
  artifact when validation fails.
- **FR-008**: The project MUST include automated checks for a representative
  successful contract, shared-model reuse, and every documented invalid-input
  category.
- **FR-009**: The public guide MUST explain the annotation convention,
  example input, command, expected contract, supported scope, and diagnostic
  path.
- **FR-010**: Project validation MUST run the published example from an
  isolated tracked-source snapshot and fail when its reviewed expected contract
  differs from the actual result.

### Key Entities

- **Operation annotation**: Source-adjacent declaration of an operation's
  method, route, summary, and responses.
- **Operation**: A uniquely identified public API action bound to one method
  and route.
- **Model definition**: A reusable public description of a request or response
  shape, including serialized fields and field descriptions.
- **Contract document**: The generated OpenAPI 3 description containing
  operations and reusable model definitions.
- **Generation diagnostic**: An actionable explanation of an invalid source
  contract that prevents document creation.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A maintainer can generate the representative API contract from a
  clean checkout with one documented command in under one minute.
- **SC-002**: The reviewed example contract contains 100% of its annotated
  operations, response statuses, and referenced public models.
- **SC-003**: Automated checks detect 100% of the documented malformed,
  duplicate, unresolved, and unsupported-input fixtures without producing a
  replacement output artifact.
- **SC-004**: Every published contract-generation example reproduces its
  reviewed expected output from an isolated tracked-source snapshot.
- **SC-005**: A user unfamiliar with the feature can identify its supported
  annotation convention, run the example, and diagnose an invalid annotation
  using only the published guide.

## Assumptions

- This feature is a deliberately small, educational contract generator, not a
  replacement for full framework discovery or broad compatibility with other
  annotation ecosystems.
- The first example covers documented GET operations and public request or
  response models; authentication schemes, servers, security policies, and
  arbitrary vendor extensions are outside its initial scope.
- The existing code-generation example validation and public documentation site
  are the appropriate channels for demonstrating and verifying this feature.
- Existing public library behavior remains compatible; the generator is a
  consumer-owned example unless later planning identifies a justified reusable
  library addition.
