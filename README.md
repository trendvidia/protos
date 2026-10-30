# protos

The shared schemas of the [protowire](https://github.com/trendvidia/protowire) stack, with their generated Go code.

If your schema imports one of these files, the code you generate from it refers to the types here. In Go, that means `protoc-gen-go` needs to know which Go package each file lives in, and the package has to exist. This repository is that package. It depends only on `google.golang.org/protobuf`.

| File | Go package | Java package | C# namespace |
|---|---|---|---|
| `pxf/annotations.proto`, `pxf/bignum.proto`, `pxf/secret.proto` | `github.com/trendvidia/protos/pxf` | `org.protowire.proto.pxf` | `Protowire.Pxf` |
| `sbe/annotations.proto` | `github.com/trendvidia/protos/sbe` | `org.protowire.proto.sbe` | `Protowire.Sbe` |
| `envelope/v1/envelope.proto` | `github.com/trendvidia/protos/envelope/v1` | `org.protowire.envelope.v1` | `Protowire.Envelopes.V1` |
| `protowire/schema/v1/*.proto` | `github.com/trendvidia/protos/protowire/schema/v1` | `org.protowire.proto.schema.v1` | `Protowire.Schema.V1` |

The Java and C# names are the ones protowire-java and protowire-csharp already publish, so a port that copies these files keeps its class names. `java_outer_classname` is left unset, so the outer class keeps protoc's default name (`Annotations`, `Bignum`).

Swift, TypeScript, C++, Rust, Python and Dart have no per-file package option. Their names come from the proto package or the file path.

## `protowire/schema/v1/annotations.proto` and gRPC reflection

`annotations.proto` uses v1.2 schema syntax (`@validate`, `@required`), so only `pxf build` compiles it. Once compiled, it declares no types, so your generated code builds and runs without registering it.

gRPC server reflection does need it registered. Without it, a reflection client cannot rebuild a file that imports it: `could not resolve import "protowire/schema/v1/annotations.proto"`. A server that offers reflection imports the opt-in package:

```go
import _ "github.com/trendvidia/protos/protowire/schema/v1/annotationsreg"
```

## Regenerating

```bash
PXF=/path/to/pxf ./scripts/generate.sh
```

It needs `protoc` and `pxf`. `protoc-gen-go` is built at the version `go.mod` pins. The `pxf` you pass must bundle this repository's schemas, because `pxf build` serves its bundled copy of `annotations.proto` ahead of any file on disk.

## Status

This repository is the canonical home of these files. protowire bundles exact copies under `proto/` and tests that they equal `protos.Files` of the version it pins, so an edit lands here first and reaches protowire when it bumps this module. The migration is tracked in [#1](https://github.com/trendvidia/protos/issues/1).

protocompile v0.35.1 and earlier register some of these files from their own `gen/pxf` and `gen/protowire/schema/v1`, so a program must not link those together with this module. It panics at startup: `proto: file "pxf/bignum.proto" is already registered`. Later protocompile releases make both packages aliases of this module (trendvidia/protocompile#271).
