#!/usr/bin/env bash
# SPDX-License-Identifier: MIT
# Copyright (c) 2026 TrendVidia, LLC.
#
# Fails when a schema here differs from protowire's copy in anything but
# its option lines.
#
# This repository is the canonical home of these schemas (#1), but until
# protowire bundles them from here (#1, step 5) people still edit
# protowire's copies: protowire#429 changed two of them a day after this
# repository was cut. This check is what notices. Delete it at step 5.
#
# PROTOWIRE_DIR compares against a local checkout instead of fetching
# protowire's main branch.
set -euo pipefail
cd "$(dirname "$0")/.."

files=(
  envelope/v1/envelope.proto
  protowire/schema/v1/annotations.proto
  protowire/schema/v1/descriptor.proto
  protowire/schema/v1/report.proto
  pxf/annotations.proto
  pxf/bignum.proto
  pxf/secret.proto
  sbe/annotations.proto
)

body() { grep -v '^option ' | cat -s; }

status=0
for f in "${files[@]}"; do
  upstream="proto/${f#protowire/}"
  if [[ -n "${PROTOWIRE_DIR:-}" ]]; then
    theirs="$(body < "$PROTOWIRE_DIR/$upstream")"
  else
    theirs="$(curl -fsSL "https://raw.githubusercontent.com/trendvidia/protowire/main/$upstream" | body)"
  fi
  if ! diff -u <(echo "$theirs") <(body < "$f") > /dev/null; then
    echo "checkdrift: $f differs from protowire's $upstream"
    diff -u <(echo "$theirs") <(body < "$f") || true
    status=1
  fi
done
[[ $status == 0 ]] && echo "checkdrift: ${#files[@]} files match protowire main"
exit $status
