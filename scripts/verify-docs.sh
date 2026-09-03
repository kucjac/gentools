#!/usr/bin/env bash
set -euo pipefail

repo=$(git rev-parse --show-toplevel)
go_bin=${GO:-go}
artifact=${PAGES_ARTIFACT_DIR:-$(mktemp -d "${TMPDIR:-/tmp}/gentools-pages.XXXXXX")}
cleanup=false
if [[ -z ${PAGES_ARTIFACT_DIR:-} ]]; then cleanup=true; trap 'rm -rf -- "$artifact"' EXIT; fi
mkdir -p "$artifact"
tmp=$(mktemp -d "${TMPDIR:-/tmp}/gentools-example.XXXXXX")
trap 'rm -rf -- "$tmp"; $cleanup && rm -rf -- "$artifact"' EXIT
env -u GOROOT -u GOTOOLCHAIN "$go_bin" run ./examples/codegen/struct-summary -input ./examples/codegen/struct-summary/testdata -output "$tmp/summary.go"
cmp "$tmp/summary.go" "$repo/examples/codegen/struct-summary/testdata/zz_summary.golden"
for page in "$repo"/docs/*.html; do
  test -f "$page"
  while IFS= read -r link; do
    case $link in http*|\#*|mailto:*) continue;; esac
    test -e "$(dirname "$page")/$link" || { echo "broken link $link in $page" >&2; exit 1; }
  done < <(sed -nE 's/.*href="([^"]+)".*/\1/p' "$page")
done
cp -R "$repo/docs/." "$artifact/"
echo "validated Pages artifact: $artifact"
