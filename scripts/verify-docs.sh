#!/usr/bin/env bash
set -euo pipefail

repo=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
go_bin=${GO:-go}
if [[ ${GENTOOLS_DOCS_SNAPSHOT:-} != 1 ]]; then
  snapshot=$(mktemp -d "${TMPDIR:-/tmp}/gentools-docs-snapshot.XXXXXX")
  trap 'rm -rf -- "$snapshot"' EXIT
  "$repo/scripts/snapshot-tracked-source.sh" --committed "$snapshot/source"
  artifact_arg=${PAGES_ARTIFACT_DIR:-}
  if [[ -n $artifact_arg && $artifact_arg != /* ]]; then
    artifact_arg=$repo/$artifact_arg
  fi
  GENTOOLS_DOCS_SNAPSHOT=1 GO="$go_bin" PAGES_ARTIFACT_DIR="$artifact_arg" \
    "$snapshot/source/scripts/verify-docs.sh"
  exit 0
fi
artifact=${PAGES_ARTIFACT_DIR:-$(mktemp -d "${TMPDIR:-/tmp}/gentools-pages.XXXXXX")}
cleanup=false
if [[ -z ${PAGES_ARTIFACT_DIR:-} ]]; then cleanup=true; trap 'rm -rf -- "$artifact"' EXIT; fi
mkdir -p "$artifact"
tmp=$(mktemp -d "${TMPDIR:-/tmp}/gentools-example.XXXXXX")
trap 'rm -rf -- "$tmp"; $cleanup && rm -rf -- "$artifact"' EXIT
normalize_newlines() {
  sed 's/\r$//' "$1" > "$2"
}
env -u GOROOT -u GOTOOLCHAIN "$go_bin" run ./examples/codegen/struct-summary -input ./examples/codegen/struct-summary/testdata -output "$tmp/summary.go"
normalize_newlines "$tmp/summary.go" "$tmp/summary.normalized.go"
normalize_newlines "$repo/examples/codegen/struct-summary/testdata/zz_summary.golden" "$tmp/summary.golden.normalized"
cmp "$tmp/summary.normalized.go" "$tmp/summary.golden.normalized"
env -u GOROOT -u GOTOOLCHAIN "$go_bin" run ./examples/codegen/struct-accessors -input ./examples/codegen/struct-accessors/testdata -type Account -fields ID,Email -output "$tmp/account_accessors.go"
normalize_newlines "$tmp/account_accessors.go" "$tmp/account_accessors.normalized.go"
normalize_newlines "$repo/examples/codegen/struct-accessors/testdata/zz_account_accessors.golden.go" "$tmp/account_accessors.golden.normalized"
cmp "$tmp/account_accessors.normalized.go" "$tmp/account_accessors.golden.normalized"
env -u GOROOT -u GOTOOLCHAIN "$go_bin" run ./examples/codegen/openapi-contract -input ./examples/codegen/openapi-contract/testdata -output "$tmp/openapi.json"
normalize_newlines "$tmp/openapi.json" "$tmp/openapi.normalized.json"
normalize_newlines "$repo/examples/codegen/openapi-contract/testdata/openapi.golden.json" "$tmp/openapi.golden.normalized.json"
cmp "$tmp/openapi.normalized.json" "$tmp/openapi.golden.normalized.json"
while IFS= read -r -d '' page; do
  while IFS= read -r link; do
    case $link in http*|\#*|mailto:*) continue;; esac
    test -e "$(dirname "$page")/$link" || { echo "broken link $link in $page" >&2; exit 1; }
  done < <(sed -nE 's/.*href="([^"]+)".*/\1/p' "$page")
done < <(find "$repo/docs" -type f -name '*.html' -print0 | sort -z)
cp -R "$repo/docs/." "$artifact/"
echo "validated Pages artifact: $artifact"
