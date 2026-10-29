// SPDX-License-Identifier: MIT
// Copyright (c) 2026 TrendVidia, LLC.

//go:build tools

// Package tools pins the code generator scripts/generate.sh builds, so its
// version is the one go.mod requires.
package tools

import _ "google.golang.org/protobuf/cmd/protoc-gen-go"
