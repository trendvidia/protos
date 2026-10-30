#!/usr/bin/env bash
# SPDX-License-Identifier: MIT
# Copyright (c) 2026 TrendVidia, LLC.
#
# Pushes the schemas to buf.build/trendvidia/protos, labelled with a release
# tag. With --check it stages and builds the module and pushes nothing.
#
#   scripts/push-bsr.sh v0.1.3
#   scripts/push-bsr.sh --check
#
# The module is staged in a temporary directory because
# protowire/schema/v1/annotations.proto declares `annotation`, RFC-001 syntax
# that stock buf cannot compile, and buf.yaml can exclude only directories.
# Nothing else here imports that file, and pxf bundles its own copy.
set -euo pipefail
cd "$(dirname "$0")/.."

label="${1:?usage: push-bsr.sh <tag> | --check}"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT

cp bsr/buf.yaml LICENSE README.md "$work/"
git ls-files '*.proto' \
  | grep -v '^internal/' \
  | grep -vx 'protowire/schema/v1/annotations.proto' \
  | while read -r f; do
      mkdir -p "$work/$(dirname "$f")"
      cp "$f" "$work/$f"
    done

buf build "$work" -o /dev/null
(cd "$work" && buf ls-files)

if [ "$label" != "--check" ]; then
  (cd "$work" && buf push --label "$label" \
    --source-control-url "https://github.com/trendvidia/protos/tree/$label")
fi
