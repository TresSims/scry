#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
out_dir="$repo_root/testdata/bundle/generated/"

mkdir -p "$out_dir"

for dir in "$repo_root"/testdata/bundle/bad_plugins/*/; do
  name="$(basename "$dir")"
  echo "Building $name..."
  go build -C "$repo_root" -buildmode=plugin -o "$out_dir"/$name.so "./testdata/bundle/bad_plugins/$name"
done

echo "Done. Plugins written to $out_dir"
