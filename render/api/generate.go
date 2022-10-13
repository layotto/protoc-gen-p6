package api

import (
	"fmt"
	"github.com/layotto/protoc-gen-p6/utils"
	"google.golang.org/protobuf/compiler/protogen"
)

func GenerateGeneratedComponent(gen *protogen.Plugin, files []*protogen.File) *protogen.GeneratedFile {
	filename := "grpc/generated.go"
	g := gen.NewGeneratedFile(filename, files[0].GoImportPath)

	utils.AddHeader(g, false)
	g.P("package ", "grpc")
	g.P()

	g.P(`type GeneratedComponents struct {`)
	for _, file := range files {
		for _, service := range file.Services {
			// add import
			importPath := "mosn.io/layotto/components/" + string(file.GoPackageName)
			pkgName := g.QualifiedGoIdent(protogen.GoImportPath(importPath).Ident(""))
			// e.g.
			// BlogService map[string]blog.BlogService
			g.P(fmt.Sprintf("\t%s\tmap[string]%s%s", service.GoName, pkgName, service.GoName))
			g.P()
		}
	}
	g.P("}")

	return g
}

func GenerateApplicationContext(gen *protogen.Plugin, files []*protogen.File) *protogen.GeneratedFile {
	filename := "grpc/context_generated.go"
	g := gen.NewGeneratedFile(filename, files[0].GoImportPath)

	utils.AddHeader(g, false)

	g.P("package ", "grpc")
	g.P()
	g.QualifiedGoIdent(protogen.GoImportPath("github.com/dapr/components-contrib/bindings").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("github.com/dapr/components-contrib/pubsub").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("github.com/dapr/components-contrib/secretstores").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("github.com/dapr/components-contrib/state").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/components/configstores").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/components/custom").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/components/file").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/components/hello").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/components/lock").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/components/oss").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/components/rpc").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/components/sequencer").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/pkg/runtime/lifecycle").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/components/pkg/common").Ident(""))
	g.P()

	g.P(`// ApplicationContext contains all you need to construct your GrpcAPI, such as all the components.
// For example, your "SuperState" GrpcAPI can hold the "StateStores" components and use them to implement your own "Super State API" logic.
type ApplicationContext struct {
	AppId                 string
	Hellos                map[string]hello.HelloService
	ConfigStores          map[string]configstores.Store
	Rpcs                  map[string]rpc.Invoker
	PubSubs               map[string]pubsub.PubSub
	StateStores           map[string]state.Store
	Files                 map[string]file.File
	Oss                   map[string]oss.Oss
	LockStores            map[string]lock.LockStore
	Sequencers            map[string]sequencer.Store
	SendToOutputBindingFn func(name string, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error)
	SecretStores          map[string]secretstores.SecretStore
	DynamicComponents     map[lifecycle.ComponentKey]common.DynamicComponent
	CustomComponent       map[string]map[string]custom.Component`)

	for _, file := range files {
		for _, service := range file.Services {
			// add import
			importPath := "mosn.io/layotto/components/" + string(file.GoPackageName)
			pkgName := g.QualifiedGoIdent(protogen.GoImportPath(importPath).Ident(""))

			// e.g.
			// BlogService map[string]blog.BlogService
			g.P(fmt.Sprintf("\t%s\tmap[string]%s%s", service.GoName, pkgName, service.GoName))
			g.P()
		}
	}
	g.P("}")

	return g
}
