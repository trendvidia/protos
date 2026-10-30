// SPDX-License-Identifier: MIT
// Copyright (c) 2026 TrendVidia, LLC.

// Package check holds what the tests in this directory share. The tests
// here link the test schema without the opt-in reflection package; the
// ones in optin link it too. That difference is the thing under test, so
// the two cannot share a test binary.
package check

import (
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Files lists every shared schema, each with the language options it
// must carry. The Java and C# names are the ones protowire-java and
// protowire-csharp already publish, so a port that vendors these files
// keeps its class names; an empty JavaOuterClassname means the option is
// unset, which keeps protoc's default outer class (pxf.Annotations, as
// protowire-java names it).
var Files = []struct {
	Path, GoPackage, JavaPackage, JavaOuterClassname, CsharpNamespace string
}{
	{"envelope/v1/envelope.proto", "github.com/trendvidia/protos/envelope/v1;envelopev1", "org.protowire.envelope.v1", "", "Protowire.Envelopes.V1"},
	{"protowire/schema/v1/annotations.proto", "github.com/trendvidia/protos/protowire/schema/v1;schemav1", "org.protowire.proto.schema.v1", "", "Protowire.Schema.V1"},
	{"protowire/schema/v1/descriptor.proto", "github.com/trendvidia/protos/protowire/schema/v1;schemav1", "org.protowire.proto.schema.v1", "", "Protowire.Schema.V1"},
	{"protowire/schema/v1/report.proto", "github.com/trendvidia/protos/protowire/schema/v1;schemav1", "org.protowire.proto.schema.v1", "", "Protowire.Schema.V1"},
	{"pxf/annotations.proto", "github.com/trendvidia/protos/pxf;pxf", "org.protowire.proto.pxf", "", "Protowire.Pxf"},
	{"pxf/bignum.proto", "github.com/trendvidia/protos/pxf;pxf", "org.protowire.proto.pxf", "", "Protowire.Pxf"},
	{"pxf/secret.proto", "github.com/trendvidia/protos/pxf;pxf", "org.protowire.proto.pxf", "", "Protowire.Pxf"},
	{"sbe/annotations.proto", "github.com/trendvidia/protos/sbe;sbe", "org.protowire.proto.sbe", "", "Protowire.Sbe"},
}

// Closure returns fd and every file it imports, transitively, the way
// grpc-go's reflection service collects them (reflection/internal,
// FileDescWithDependencies): an import that resolved to no registered
// file is a placeholder, and has nothing to send.
func Closure(fd protoreflect.FileDescriptor) *descriptorpb.FileDescriptorSet {
	set := &descriptorpb.FileDescriptorSet{}
	seen := map[string]bool{}
	var walk func(protoreflect.FileDescriptor)
	walk = func(f protoreflect.FileDescriptor) {
		if seen[f.Path()] || f.IsPlaceholder() {
			return
		}
		seen[f.Path()] = true
		for i := range f.Imports().Len() {
			walk(f.Imports().Get(i).FileDescriptor)
		}
		set.File = append(set.File, protodesc.ToFileDescriptorProto(f))
	}
	walk(fd)
	return set
}
