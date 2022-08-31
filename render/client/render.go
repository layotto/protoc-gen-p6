package client

import (
	"fmt"
	"github.com/seeflood/protoc-gen-p6/utils"
	"google.golang.org/protobuf/compiler/protogen"
)

type renderImpl struct {
	gen *protogen.Plugin
}

func (r *renderImpl) Render(files []*protogen.File) error {
	g := r.initGeneratedFile()
	r.generateSDK(files, g)
	return nil
}

type SdkRender interface {
	Render([]*protogen.File) error
}

func NewRender(gen *protogen.Plugin) SdkRender {
	return &renderImpl{
		gen: gen,
	}
}

func (r *renderImpl) initGeneratedFile() *protogen.GeneratedFile {
	//var fileWithService *protogen.File
	//for _, f := range r.gen.Files {
	//	if len(f.Services) > 0 {
	//		fileWithService = f
	//	}
	//}
	//path, _ := utils.Directory(fileWithService)
	//filename := path + "/client/client_generated.go"
	filename := "client/client_generated.go"

	g := r.gen.NewGeneratedFile(filename, r.gen.Files[0].GoImportPath)

	utils.AddHeader(g, false)

	g.P("package ", "client")
	g.P()

	runtimev1pbPkg = g.QualifiedGoIdent(runtimev1pbImport.Ident(""))
	contextPkg = g.QualifiedGoIdent(contextImport.Ident(""))
	grpcPkg = g.QualifiedGoIdent(grpcImport.Ident(""))
	g.P()

	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var _ ", contextPkg, "Context")
	g.P()

	return g
}

func (r *renderImpl) generateSDK(files []*protogen.File, g *protogen.GeneratedFile) {
	// 1. add `Client` interface
	g.P(`// Client is the interface for runtime client implementation.
type Client interface {
	runtimeAPI

	s3.ObjectStorageServiceClient
`)

	for _, file := range files {
		for _, service := range file.Services {
			// add import
			pkgName := g.QualifiedGoIdent(protogen.GoIdent{
				GoImportPath: file.GoImportPath,
			})
			g.P("// ", file.GoImportPath)
			g.P(pkgName, service.GoName, "Client")
			g.P()
		}
	}

	g.P("}")

	// 2. constructor
	g.P(`
// NewClientWithConnection instantiates runtime client using specific connection.
func NewClientWithConnection(conn *grpc.ClientConn) Client {
	return &GRPCClient{
		connection:                 conn,`)
	// protoClient:                v1.NewRuntimeClient(conn),
	g.P("\t\tprotoClient:                ", runtimev1pbPkg, "NewRuntimeClient(conn),")

	// clients for different services
	// e.g. 		BlogServiceClient: blog.NewBlogServiceClient(conn),
	for _, file := range files {
		for _, service := range file.Services {
			// add import
			pkgName := g.QualifiedGoIdent(protogen.GoIdent{
				GoImportPath: file.GoImportPath,
			})
			// add comment
			g.P("// ", file.GoImportPath)
			// generate client
			clientName := service.GoName + "Client"
			g.P(fmt.Sprintf("\t\t%s: %sNew%s(conn),", clientName, pkgName, clientName))
			g.P()
		}
	}
	// the end of the constructor
	g.P(`	}
}`)
	// 3. generate the GRPCClient struct
	g.P(`
// GRPCClient is the gRPC implementation of runtime client.
type GRPCClient struct {
	connection  *grpc.ClientConn`)
	g.P("\tprotoClient ", runtimev1pbPkg, "RuntimeClient")
	// clients for different services
	// e.g.  blog.BlogServiceClient
	for _, file := range files {
		for _, service := range file.Services {
			// add import
			pkgName := g.QualifiedGoIdent(protogen.GoIdent{
				GoImportPath: file.GoImportPath,
			})
			g.P("// ", file.GoImportPath)
			clientName := service.GoName + "Client"
			g.P(pkgName, clientName)
		}
	}
	g.P("}")
}
