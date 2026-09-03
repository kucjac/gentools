#!/usr/bin/env bash
set -euo pipefail
repo=$(git rev-parse --show-toplevel)
grep -q 'GENTOOLS_DOCS_SNAPSHOT' "$repo/scripts/verify-docs.sh"
grep -q 'snapshot-tracked-source.sh' "$repo/scripts/verify-docs.sh"
grep -q 'snapshot-tracked-source.sh" --committed' "$repo/scripts/verify-docs.sh"
test -f "$repo/docs/index.html"
test -f "$repo/docs/code-generation.html"
test -f "$repo/docs/publishing.html"
test -f "$repo/docs/testing/integration.html"
test -f "$repo/docs/testing/unit-test-gaps.html"
test -f "$repo/docs/testing/behavior-inventory.html"
test -f "$repo/docs/style.css"
grep -q 'href="style.css"' "$repo/docs/index.html"
grep -q 'name="viewport"' "$repo/docs/index.html"
grep -q 'testing/integration.html' "$repo/docs/index.html"
grep -q 'testing/unit-test-gaps.html' "$repo/docs/index.html"
grep -q 'struct-accessors' "$repo/docs/code-generation.html"
grep -q '<title>Code generation' "$repo/docs/code-generation.html"
grep -q 'Generate typed Go accessors' "$repo/docs/code-generation.html"
grep -q 'func (v Account) GetID() string {' "$repo/docs/code-generation.html"
grep -q 'Struct summary' "$repo/docs/code-generation.html"
grep -q 'examples/codegen/struct-summary/testdata' "$repo/docs/code-generation.html"
grep -q 'var StructNames = \[\]string{"Account", "Session"}' "$repo/docs/code-generation.html"
grep -q 'If the summary is empty' "$repo/docs/code-generation.html"
grep -q 'openapi-contract' "$repo/docs/code-generation.html"
grep -q 'openapi method=GET' "$repo/README.md"
grep -q 'openapi-contract' "$repo/scripts/verify-docs.sh"
grep -q 'struct-accessors' "$repo/README.md"
grep -q 'GitHub Pages' "$repo/README.md"
grep -q 'kucjac.github.io/gentools' "$repo/README.md"
if grep -q 'href="missing.html"' "$repo/docs/index.html"; then
  echo "fixture contains intentionally broken link" >&2
  exit 1
fi
