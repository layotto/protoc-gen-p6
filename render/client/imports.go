package client

import "google.golang.org/protobuf/compiler/protogen"

const (
	contextImport      = protogen.GoImportPath("context")
	grpcImport         = protogen.GoImportPath("google.golang.org/grpc")
	runtimev1pbImport  = protogen.GoImportPath("mosn.io/layotto/spec/proto/runtime/v1")
	deprecationComment = "// Deprecated: Do not use."
)

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	contextPkg     string
	runtimev1pbPkg string
	grpcPkg        string
)
