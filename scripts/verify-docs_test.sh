#!/usr/bin/env bash
set -euo pipefail
repo=$(git rev-parse --show-toplevel)
test -f "$repo/docs/index.html"
test -f "$repo/docs/code-generation.html"
if grep -q 'href="missing.html"' "$repo/docs/index.html"; then
  echo "fixture contains intentionally broken link" >&2
  exit 1
fi
