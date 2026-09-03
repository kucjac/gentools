# Research: Test Quality Guidance

## Candidate-source consumer validation

**Decision**: Retain a nested consumer module, split it into four scenario
fixtures, and run each against a temporary snapshot of the tracked candidate
source. Label this as a candidate-source contract test. Add a separate
tag/release smoke test with no `replace` directive to prove a published module
can be consumed.

**Rationale**: The existing nested module uses a local replacement and has one
protobuf test. A temporary snapshot prevents accidental reads from unrelated
working-tree state while still allowing pull-request validation of the candidate
commit. Only a no-replacement release test can prove publication and module
resolution.

**Alternatives considered**:

- Keep the single protobuf test: rejected because it does not assert fields,
  imported dependencies, failure behavior, or independent consumer use.
- Require published-version tests in every pull request: rejected because a
  candidate version is not published yet.
- Test a permanent local replacement directly: rejected because it makes the
  source-under-test boundary implicit.

## Unit-test gap assessment

**Decision**: Use a checked-in JSON behavior inventory, validated by a small
standard-library Go command and rendered as Markdown. Pair it with
`go test -coverprofile` and `go tool cover -func` evidence; classify every
exported contract and high-risk internal behavior as verified, integration-only,
excluded with rationale, or a prioritized gap.

**Rationale**: Go coverage measures basic blocks and is not a contract
inventory; it cannot establish useful assertions, error behavior, or public API
compatibility. A structured inventory can verify that all exported declarations
are classified and that cited test paths still exist.

**Alternatives considered**:

- Enforce one package-percentage threshold: rejected because a high percentage
  can still omit critical errors and because generated/unreachable code distorts
  the result.
- Adopt an external coverage service: rejected because it adds an account and
  dependency without solving behavior classification.
- Maintain an unvalidated prose list: rejected because it becomes stale and
  cannot detect unclassified API additions.

## Code-generation examples

**Decision**: Provide a runnable consumer example that loads a package through
Gentools and emits a deterministic Go source summary from selected declarations.
Treat protobuf as generated-source consumption coverage; correct and regenerate
the existing protobuf fixture only after a reproducibility preflight confirms
the generator versions and package path.

**Rationale**: Gentools is a source-analysis library, not a generator. A small
consumer-owned renderer demonstrates the requested code-generation use case
without changing the public library. The current `.proto` `go_package` value
does not match the package path asserted by its test, so regeneration cannot be
claimed safe until that mismatch is resolved.

**Alternatives considered**:

- Add a general Gentools code generator: rejected as a new public product and
  out of scope.
- Hand-edit committed generated protobuf output: rejected by the constitution
  and not reproducible.
- Document parsing only: rejected because it does not demonstrate the requested
  code-generation workflow.

## GitHub Pages publication

**Decision**: Publish a static artifact through a dedicated Pages workflow:
validate/build in one job, then deploy the artifact from `main` in a dependent
job. Use `actions/checkout@v6`, `actions/configure-pages@v5`,
`actions/upload-pages-artifact@v4`, and `actions/deploy-pages@v4`; configure
the repository's Pages source as GitHub Actions.

**Rationale**: GitHub's custom-workflow model supports static content without a
site generator, keeps deployment permissions isolated, and permits pull
requests to validate the same artifact without publishing it. Deployment needs
`contents: read`, `pages: write`, and `id-token: write`; it reports the
published URL as an output.

**Alternatives considered**:

- Branch publishing: rejected because it mixes generated deployment output with
  source and requires deployment commits.
- One-job deployment: rejected because example and link validation must finish
  before publication.
- A documentation-site framework: rejected because a small static site does
  not justify an additional runtime or package manager.

## Sources

- [GitHub Pages custom workflows](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages)
- [Configure a GitHub Pages publishing source](https://docs.github.com/en/pages/getting-started-with-github-pages/configuring-a-publishing-source-for-your-github-pages-site)
- [Go coverage for applications and integration tests](https://go.dev/doc/build-cover)
- [Go cover command](https://pkg.go.dev/cmd/cover)
- [Go testing package](https://pkg.go.dev/testing)
