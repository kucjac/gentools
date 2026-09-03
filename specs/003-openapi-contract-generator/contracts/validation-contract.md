# Validation Contract

The canonical command is:

```sh
go run ./examples/codegen/openapi-contract \
  -input ./examples/codegen/openapi-contract/testdata \
  -output /tmp/gentools-openapi.json
```

The command exits zero and produces the reviewed JSON contract for valid input.
It exits non-zero for malformed directives, duplicate method/route pairs,
unresolved models, or unsupported fields. Diagnostics identify the source
operation or model field. Validation errors must not create a new output or
replace an existing output.

The documentation gate runs this command from a committed tracked-source
snapshot and compares the result with `testdata/openapi.golden.json` after
normalizing line endings.
