#!/usr/bin/env bash
set -euo pipefail
repo=$(git rev-parse --show-toplevel)
test -f "$repo/docs/index.html"
test -f "$repo/docs/code-generation.html"
test -f "$repo/docs/style.css"
grep -q 'href="style.css"' "$repo/docs/index.html"
grep -q 'name="viewport"' "$repo/docs/index.html"
grep -q 'GitHub Pages' "$repo/README.md"
grep -q 'kucjac.github.io/gentools' "$repo/README.md"
if grep -q 'href="missing.html"' "$repo/docs/index.html"; then
  echo "fixture contains intentionally broken link" >&2
  exit 1
fi
