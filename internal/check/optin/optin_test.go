// SPDX-License-Identifier: MIT
// Copyright (c) 2026 TrendVidia, LLC.

package optin

import (
	"slices"
	"strings"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/trendvidia/protos/internal/check"
	"github.com/trendvidia/protos/internal/testpb"
	schemav1 "github.com/trendvidia/protos/protowire/schema/v1"
	_ "github.com/trendvidia/protos/protowire/schema/v1/annotationsreg"
)

// With the opt-in package linked, the import resolves and a reflection
// client rebuilds the whole file set.
func TestOptInReflectionResolves(t *testing.T) {
	fd := testpb.File_testpb_proto
	for i := range fd.Imports().Len() {
		if imp := fd.Imports().Get(i); imp.IsPlaceholder() {
			t.Errorf("%s is still a placeholder import", imp.Path())
		}
	}
	set := check.Closure(fd)
	if _, err := protodesc.NewFiles(set); err != nil {
		t.Fatalf("rebuilding the closure: %v", err)
	}
	var names []string
	for _, f := range set.File {
		names = append(names, f.GetName())
	}
	for _, f := range check.Files {
		if !slices.Contains(names, f.Path) {
			t.Errorf("the closure of %s lacks %s; it has %q", fd.Path(), f.Path, names)
		}
	}
}

// The registered annotations.proto carries the annotation declarations
// (the 1329 carrier), so a reflection client gets what @validate means.
func TestOptInCarriesDeclarations(t *testing.T) {
	fd, err := protoregistry.GlobalFiles.FindFileByPath("protowire/schema/v1/annotations.proto")
	if err != nil {
		t.Fatal(err)
	}
	decls := proto.GetExtension(fd.Options(), schemav1.E_AnnotationDecls).(*schemav1.FileAnnotationDecls)
	var names []string
	for _, d := range decls.GetDeclarations() {
		names = append(names, d.GetName())
	}
	for _, want := range []string{"protowire.schema.v1.validate", "protowire.schema.v1.required"} {
		if !slices.Contains(names, want) {
			t.Errorf("annotation_decls lacks %s; it has %q", want, names)
		}
	}
}

// Every shared schema carries the language options its generated code is
// named by. A file that loses its go_package fails protoc-gen-go in every
// user's build, so this is the test that would notice.
func TestOptions(t *testing.T) {
	for _, f := range check.Files {
		fd, err := protoregistry.GlobalFiles.FindFileByPath(f.Path)
		if err != nil {
			t.Errorf("%s: %v", f.Path, err)
			continue
		}
		o := fd.Options().(*descriptorpb.FileOptions)
		for _, c := range []struct{ name, got, want string }{
			{"go_package", o.GetGoPackage(), f.GoPackage},
			{"java_package", o.GetJavaPackage(), f.JavaPackage},
			{"java_outer_classname", o.GetJavaOuterClassname(), f.JavaOuterClassname},
			{"csharp_namespace", o.GetCsharpNamespace(), f.CsharpNamespace},
		} {
			if c.got != c.want {
				t.Errorf("%s: %s = %q, want %q", f.Path, c.name, c.got, c.want)
			}
		}
		if !o.GetJavaMultipleFiles() {
			t.Errorf("%s: java_multiple_files is not set", f.Path)
		}
	}
}

// protoc-gen-go finds a Go import path for every file the test schema
// imports, with no M flags: the failure trendvidia/protowire#424 reported.
// The control strips one file's go_package and expects the same error
// users got.
func TestProtocGenGoResolvesImports(t *testing.T) {
	set := check.Closure(testpb.File_testpb_proto)
	request := func(set *descriptorpb.FileDescriptorSet) *pluginpb.CodeGeneratorRequest {
		return &pluginpb.CodeGeneratorRequest{FileToGenerate: []string{"testpb.proto"}, ProtoFile: set.File}
	}
	if _, err := (protogen.Options{}).New(request(set)); err != nil {
		t.Fatalf("protoc-gen-go: %v", err)
	}

	stripped := proto.Clone(set).(*descriptorpb.FileDescriptorSet)
	for _, f := range stripped.File {
		if f.GetName() == "pxf/secret.proto" {
			f.Options.GoPackage = nil
		}
	}
	_, err := (protogen.Options{}).New(request(stripped))
	if err == nil || !strings.Contains(err.Error(), `unable to determine Go import path for "pxf/secret.proto"`) {
		t.Errorf("control: err = %v, want unable to determine Go import path for pxf/secret.proto", err)
	}
}
