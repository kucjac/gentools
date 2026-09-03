# Unit-test gap assessment

Run `make test-gaps` with the toolchain named in `.go-version`. The assessment
checks that every exported `parser` and `types` declaration has an inventory
classification, every evidence path exists, and high-risk gaps are rejected.

`go tool cover -func` and the coverage profile are evidence only: Go coverage
counts basic blocks, not public contracts or quality of assertions. Add new
symbols to `behavior-inventory.json` in the same change, classify integration-
only and generated behavior with a rationale, and add a focused regression test
for every high-risk behavior. The selected toolchain and compatibility policy
are documented in the repository [README](../../README.md).
