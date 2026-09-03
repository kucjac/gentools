#!/usr/bin/env bash
set -euo pipefail

mode=candidate
if [[ ${1:-} == --committed ]]; then
  mode=committed
  shift
fi
[[ $# -eq 1 ]] || { echo "usage: $0 [--committed] DESTINATION" >&2; exit 2; }
destination=$1
repo=$(git rev-parse --show-toplevel)
if [[ -e $destination ]]; then
  echo "snapshot destination already exists: $destination" >&2
  exit 1
fi
mkdir -p "$destination"
cleanup() { rm -rf -- "$destination"; }
trap cleanup ERR
if [[ $mode == committed ]]; then
  git -C "$repo" archive --format=tar HEAD | tar -x -C "$destination"
else
  while IFS= read -r -d '' path; do
    source_path=$repo/$path
    [[ -f $source_path ]] || { echo "tracked source is missing: $path" >&2; exit 1; }
    mkdir -p "$destination/$(dirname "$path")"
    cp -- "$source_path" "$destination/$path"
  done < <(git -C "$repo" ls-files -z)
fi
[[ -f $destination/go.mod ]] || { echo "snapshot has no root go.mod" >&2; exit 1; }
trap - ERR
