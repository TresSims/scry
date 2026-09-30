#!/usr/bin/env bash

set -euo pipefail

test_opts=""

while getopts "i:o:t" opt; do
  case $opt in
    i) in_dir=$OPTARG;;
    o) out_dir=$OPTARG;;
    t) test_opts="-cover -covermode=set -coverpkg=./..."
  esac
done

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "$out_dir"

for dir in "$in_dir"/*/; do
  name="$(basename "$dir")"
  echo "Building $name..."
  go build -C "$repo_root" \
    -buildmode=plugin -o "$out_dir/$name.so" \
    $test_opts \
    "$in_dir/$name"
done

echo $out_dir
ls $out_dir

echo "Done. Plugins written to $out_dir"
