#!/usr/bin/env bash
set -euo pipefail

repo=$(git rev-parse --show-toplevel)
go_bin=${GO:-go}
coverage=${COVERAGE_FILE:-$(mktemp "${TMPDIR:-/tmp}/gentools.cover.XXXXXX")}
if [[ -z ${COVERAGE_FILE:-} ]]; then trap 'rm -f -- "$coverage"' EXIT; fi
env -u GOROOT -u GOTOOLCHAIN "$go_bin" test -coverprofile="$coverage" ./...
env -u GOROOT -u GOTOOLCHAIN "$go_bin" run ./cmd/testinventory assess --inventory "$repo/docs/testing/behavior-inventory.json" --report "$repo/docs/testing/behavior-inventory.md" --coverage "$coverage"
git -C "$repo" diff --exit-code -- docs/testing/behavior-inventory.md
