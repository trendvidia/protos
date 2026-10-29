// SPDX-License-Identifier: MIT
// Copyright (c) 2026 TrendVidia, LLC.

// Package protos holds the shared protowire schemas as source. Their
// generated Go is in the packages below it, one per proto directory.
//
// protowire bundles copies of these files and tests that its copies equal
// Files byte for byte, so edits are made here first.
package protos

import "embed"

// Files holds every .proto in this repository at its import path, e.g.
// "pxf/bignum.proto". The test schema under internal/ is not included.
//
//go:embed envelope/v1/*.proto protowire/schema/v1/*.proto pxf/*.proto sbe/*.proto
var Files embed.FS
