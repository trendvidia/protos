// SPDX-License-Identifier: MIT
// Copyright (c) 2026 TrendVidia, LLC.

package protos_test

import (
	"io/fs"
	"slices"
	"testing"

	"github.com/trendvidia/protos"
	"github.com/trendvidia/protos/internal/check"
)

// Files holds exactly the shared schemas, so protowire's copy test sees
// every one of them and nothing else.
func TestFilesAreTheSharedSchemas(t *testing.T) {
	var got []string
	err := fs.WalkDir(protos.Files, ".", func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			got = append(got, p)
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	var want []string
	for _, f := range check.Files {
		want = append(want, f.Path)
	}
	slices.Sort(want)
	slices.Sort(got)
	if !slices.Equal(got, want) {
		t.Errorf("Files holds %q, want %q", got, want)
	}
}
