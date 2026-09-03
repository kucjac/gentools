# Consumer integration scenarios

`scripts/test-integration.sh all` snapshots tracked source into a temporary
directory before it runs `inspect`, `generated`, `crosspackage`, and `invalid`.
Its output is **candidate-source evidence**: it proves the checked-out library
contract, not a published release. Each scenario has a manifest entry in
`internal/integration/scenarios/scenarios.json` that declares its input,
expected result, isolation boundary, and diagnostic expectation.

Use the tag-triggered release smoke workflow for **published-module evidence**.
It intentionally has no local `replace` directive.

## Generated protobuf preflight

The protobuf fixture is regenerated, never edited by hand. From the repository
root, use the checked generator and include paths:

```sh
protoc \
  -I internal/integration/protobuf -I /usr/include \
  --go_out=paths=source_relative:internal/integration/protobuf \
  internal/integration/protobuf/models.proto
git diff --check -- internal/integration/protobuf/models.pb.go
```

Record the `protoc` and `protoc-gen-go` versions with any generated-output
change, then run the generated consumer scenario.
