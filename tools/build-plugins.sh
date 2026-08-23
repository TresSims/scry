#!/usr/bin/env bash
# Build every plugin in examples/ as a Go plugin (.so) into plugins/
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
out_dir="$repo_root/plugins"

mkdir -p "$out_dir"

for dir in "$repo_root"/examples/*/; do
	name="$(basename "$dir")"
	echo "Building $name..."
	go build -C "$repo_root" -buildmode=plugin -o "$out_dir/$name.so" "./examples/$name"
done

echo "Done. Plugins written to $out_dir"
