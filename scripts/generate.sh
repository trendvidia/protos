#!/usr/bin/env bash
# SPDX-License-Identifier: MIT
# Copyright (c) 2026 TrendVidia, LLC.
#
# Regenerates every .pb.go in this repository.
#
# Needs protoc and pxf. protoc-gen-go is built from the version go.mod
# pins. PXF names the pxf binary (default: pxf on PATH). Its bundled
# schemas must be this repository's: pxf serves its bundled copy of
# protowire/schema/v1/annotations.proto ahead of any file on disk, and
# the opt-in package registers what pxf lowers.
set -euo pipefail
cd "$(dirname "$0")/.."

PXF="${PXF:-pxf}"
INCLUDE="${PROTOC_INCLUDE:-$(dirname "$(dirname "$(command -v protoc)")")/include}"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT

go build -o "$work/protoc-gen-go" google.golang.org/protobuf/cmd/protoc-gen-go
gen=(--plugin=protoc-gen-go="$work/protoc-gen-go")

# The stock-syntax schemas: protoc compiles them from source.
protoc -I . -I "$INCLUDE" "${gen[@]}" --go_out=. --go_opt=paths=source_relative \
  envelope/v1/envelope.proto \
  protowire/schema/v1/descriptor.proto \
  protowire/schema/v1/report.proto \
  pxf/annotations.proto \
  pxf/bignum.proto \
  pxf/secret.proto \
  sbe/annotations.proto

# protowire/schema/v1/annotations.proto is v1.2 syntax, so only pxf build
# compiles it. The test schema imports it, so one image carries both.
"$PXF" build -o "$work/image.binpb" internal/testpb

# The opt-in reflection package registers the lowered annotations.proto.
# Its go_package names protowire/schema/v1, where user code imports it
# from and where it stays unregistered unless a program opts in; the M
# flag puts the generated file in annotationsreg instead.
protoc --descriptor_set_in="$work/image.binpb" "${gen[@]}" \
  --go_out=. --go_opt=module=github.com/trendvidia/protos \
  --go_opt=Mprotowire/schema/v1/annotations.proto="github.com/trendvidia/protos/protowire/schema/v1/annotationsreg;annotationsreg" \
  protowire/schema/v1/annotations.proto

# The test schema, generated the way a user's would be: without the M
# flag, so it imports annotations.proto from its go_package.
protoc --descriptor_set_in="$work/image.binpb" "${gen[@]}" \
  --go_out=. --go_opt=module=github.com/trendvidia/protos \
  testpb.proto
