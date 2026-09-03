#!/usr/bin/env bash
set -euo pipefail

repo=$(git rev-parse --show-toplevel)
tmp=$(mktemp -d)
untracked=$repo/.gentools-snapshot-test-untracked
trap 'rm -rf -- "$tmp" "$untracked"' EXIT
mkdir -p "$untracked"
printf 'not tracked\n' > "$untracked/file"
"$repo/scripts/snapshot-tracked-source.sh" "$tmp/snapshot"
test -f "$tmp/snapshot/go.mod"
test ! -e "$tmp/snapshot/.gentools-snapshot-test-untracked"
git -C "$repo" diff --quiet -- go.mod || { echo "working tree fixture unexpectedly changed" >&2; exit 1; }
test "$(cmp -s "$repo/go.mod" "$tmp/snapshot/go.mod"; echo $?)" = 0
if "$repo/scripts/snapshot-tracked-source.sh" "$tmp/snapshot"; then
  echo "existing destination was accepted" >&2
  exit 1
fi
