# Quickstart: Validate Test Quality Guidance

## Prerequisites

- Checkout at the repository root.
- Go version selected by `.go-version`.
- `GODEBUG=gotypesalias=1` for the alias fixtures, matching CI.
- Network access for first-time module downloads. Pages deployment additionally
  requires repository Pages to use GitHub Actions as its publishing source.

## Assess behavior coverage

Generate a coverage profile, then validate and render the behavior inventory:

```sh
GODEBUG=gotypesalias=1 go test -coverprofile=/tmp/gentools.coverprofile ./...
go run ./cmd/testinventory assess \
  --inventory docs/testing/behavior-inventory.json \
  --report docs/testing/behavior-inventory.md \
  --coverage /tmp/gentools.coverprofile
```

Expected result: all exported contracts are classified, cited evidence exists,
and no unapproved high-risk gap remains. Read the generated report rather than
using its coverage percentage as a completeness decision; see the [artifact
contract](contracts/quality-artifacts.md).

## Run consumer scenarios

```sh
GODEBUG=gotypesalias=1 scripts/test-integration.sh all
```

Expected result: `inspect`, `generated`, and `crosspackage` pass with their
stated output; `invalid` passes only by receiving its expected failure
diagnostic. The runner reports that this is candidate-source evidence. Use the
release smoke workflow with an explicit version to obtain published-module
evidence.

## Validate examples and Pages artifact

```sh
GODEBUG=gotypesalias=1 scripts/verify-docs.sh
```

Expected result: each documented code-generation example runs from a clean
snapshot, generated output matches, and static-site links resolve. The Pages
workflow deploys this same validated artifact only from `main`; see
[research.md](research.md) for the deployment decision and permissions.

## Recorded local evidence

On 2026-09-03, the full `make test-quality` workflow completed successfully
in under 10 minutes with the selected Go 1.27.1 toolchain. An independent
walkthrough confirmed that the Pages landing page routes to the generator,
consumer-scenario, gap-assessment, and publishing guidance, and that the
generator's checked-in golden output matches the executable example.

The committed-snapshot boundary is covered by the snapshot helper test. A
repository administrator must still enable GitHub Actions as the Pages source
and record the first deployment URL in `docs/publishing.md`.
