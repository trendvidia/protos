# protos

The shared schemas of the [protowire](https://github.com/trendvidia/protowire) stack, with their generated Go code.

If your schema imports one of these files, the code you generate from it refers to the types here. In Go, that means `protoc-gen-go` needs to know which Go package each file lives in, and the package has to exist. This repository is that package. It depends only on `google.golang.org/protobuf`.

| File | Go package | Java package | C# namespace |
|---|---|---|---|
| `pxf/annotations.proto`, `pxf/bignum.proto`, `pxf/secret.proto` | `github.com/trendvidia/protos/pxf` | `org.protowire.protos.pxf` | `Protowire.Protos.Pxf` |
| `sbe/annotations.proto` | `github.com/trendvidia/protos/sbe` | `org.protowire.protos.sbe` | `Protowire.Protos.Sbe` |
| `envelope/v1/envelope.proto` | `github.com/trendvidia/protos/envelope/v1` | `org.protowire.protos.envelope.v1` | `Protowire.Protos.Envelope.V1` |
| `protowire/schema/v1/*.proto` | `github.com/trendvidia/protos/protowire/schema/v1` | `org.protowire.protos.schema.v1` | `Protowire.Protos.Schema.V1` |

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

Not released yet. protowire still bundles its own copies of these files, and protocompile still publishes Go code for some of them. The migration, and the decisions still open, are tracked in [#1](https://github.com/trendvidia/protos/issues/1).

Until protocompile's `gen/pxf` and `gen/protowire/schema/v1` stop registering these files, a program must not link them together with this module. It panics at startup: `proto: file "pxf/bignum.proto" is already registered`.
