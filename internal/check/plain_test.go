// SPDX-License-Identifier: MIT
// Copyright (c) 2026 TrendVidia, LLC.

package check

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/trendvidia/protos/internal/testpb"
)

const annotationsPath = "protowire/schema/v1/annotations.proto"

// Every shared schema but annotations.proto registers when its package
// is linked. annotations.proto registers only through the opt-in
// package, which this test binary does not link.
func TestFilesRegister(t *testing.T) {
	_ = testpb.File_testpb_proto
	for _, f := range Files {
		_, err := protoregistry.GlobalFiles.FindFileByPath(f.Path)
		if f.Path == annotationsPath {
			if err == nil {
				t.Errorf("%s is registered without the opt-in package", f.Path)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: %v", f.Path, err)
		}
	}
}

// Without the opt-in package, generated code that imports annotations.proto
// builds and runs, with the import as a placeholder. A reflection client
// then cannot rebuild the file set: what grpc-go v1.84.0 did when measured
// against a real server (trendvidia/protos#1).
func TestWithoutOptInReflectionCannotResolve(t *testing.T) {
	fd := testpb.File_testpb_proto
	var placeholder bool
	for i := range fd.Imports().Len() {
		if imp := fd.Imports().Get(i); imp.Path() == annotationsPath {
			placeholder = imp.IsPlaceholder()
		}
	}
	if !placeholder {
		t.Fatalf("%s is not a placeholder import of %s", annotationsPath, fd.Path())
	}
	_, err := protodesc.NewFiles(Closure(fd))
	if err == nil || !strings.Contains(err.Error(), "could not resolve import \""+annotationsPath+"\"") {
		t.Errorf("rebuilding the closure: err = %v, want could not resolve import %q", err, annotationsPath)
	}
}
