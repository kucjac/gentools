# OpenAPI Contract Generator Quickstart

With Go `1.27.1` selected from `.go-version`, run:

```sh
go test ./examples/codegen/openapi-contract
go run ./examples/codegen/openapi-contract \
  -input ./examples/codegen/openapi-contract/testdata \
  -output /tmp/gentools-openapi.json
```

Inspect `/tmp/gentools-openapi.json` and compare it with
`examples/codegen/openapi-contract/testdata/openapi.golden.json`.

Operation comments use:

```text
@openapi method=GET route=/pets/{id} summary="Get a pet" response=200:Pet
```

The guide and fixture cover shared model references, JSON field names and
descriptions. Invalid fixtures in the Go tests cover malformed directives,
duplicate routes, unresolved models, unsupported shapes, and output safety.
The full documentation check runs the same command from an isolated committed
source snapshot with `make verify-docs`.

## Implementation evidence

- Go 1.27.1 focused example tests: passed.
- Root `go test ./...`: passed.
- Current-source Pages/docs validation and `scripts/verify-docs_test.sh`: passed.
- `git diff --check`: passed.
- Representative generation duration: 62ms on the development host.
- Windows/amd64 cross-compilation of the example tests: passed.
- `make build-and-test`: passed.
- GitHub Actions pipeline integration coverage, including the separate
  consumer scenarios, passed on Ubuntu, macOS, and Windows.
- Committed-snapshot `make verify-docs`: passed after the feature commit; the
  snapshot intentionally reads `HEAD`.
