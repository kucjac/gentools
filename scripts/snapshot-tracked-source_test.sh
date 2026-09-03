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

committed_repo=$tmp/committed-repo
mkdir -p "$committed_repo/scripts"
cp "$repo/scripts/snapshot-tracked-source.sh" "$committed_repo/scripts/snapshot-tracked-source.sh"
chmod +x "$committed_repo/scripts/snapshot-tracked-source.sh"
printf 'module example.invalid/committed-snapshot\n\ngo 1.27.1\n' > "$committed_repo/go.mod"
git -C "$committed_repo" init -q
git -C "$committed_repo" config user.email test@example.invalid
git -C "$committed_repo" config user.name snapshot-test
git -C "$committed_repo" add .
git -C "$committed_repo" commit -qm 'initial snapshot fixture'
printf '// working tree only\n' >> "$committed_repo/go.mod"
(cd "$committed_repo" && ./scripts/snapshot-tracked-source.sh --committed "$tmp/committed-snapshot")
git -C "$committed_repo" show HEAD:go.mod > "$tmp/committed-go.mod"
cmp "$tmp/committed-go.mod" "$tmp/committed-snapshot/go.mod"
if cmp -s "$committed_repo/go.mod" "$tmp/committed-snapshot/go.mod"; then
  echo "committed snapshot included dirty tracked content" >&2
  exit 1
fi
