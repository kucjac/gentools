#!/usr/bin/env bash
set -euo pipefail

repo=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
go_bin=${GO:-go}
scenario=${1:-all}
if [[ ${GENTOOLS_SNAPSHOT:-} != 1 ]]; then
  tmp=$(mktemp -d "${TMPDIR:-/tmp}/gentools-integration.XXXXXX")
  trap 'rm -rf -- "$tmp"' EXIT
  "$repo/scripts/snapshot-tracked-source.sh" "$tmp/source"
  echo "candidate-source evidence: $tmp/source"
  GENTOOLS_SNAPSHOT=1 GO="$go_bin" "$tmp/source/scripts/test-integration.sh" "$scenario"
  exit 0
fi

(cd "$repo/internal/integration" && env -u GOROOT -u GOTOOLCHAIN "$go_bin" run ./runner \
  -manifest scenarios/scenarios.json -scenario "$scenario" -go "$go_bin")
