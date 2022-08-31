package gentemplate

import (
	"github.com/seeflood/protoc-gen-p6/utils"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	contextPkg         = protogen.GoImportPath("context")
	deprecationComment = "// Deprecated: Do not use."
)

var methodSets = make(map[string]int)

func GenerateComponentInterface(gen *protogen.Plugin, file *protogen.File, filename string) *protogen.GeneratedFile {
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	utils.AddHeader(g, false)
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.QualifiedGoIdent(contextPkg.Ident(""))
	g.P()

	for _, service := range file.Services {
		r := newRender(gen, file, g, service, string(file.GoPackageName))
		g.P(r.renderComponent())
	}

	return g
}

func GenerateComponentTypes(gen *protogen.Plugin, file *protogen.File, filename string) *protogen.GeneratedFile {
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	utils.AddHeader(g, false)
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.QualifiedGoIdent(protogen.GoImportPath("fmt").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/components/ref").Ident(""))
	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/components/pkg/info").Ident(""))
	g.P()

	for _, service := range file.Services {
		r := newRender(gen, file, g, service, string(file.GoPackageName))
		g.P(r.doRender("types", componentTypesTpl))
	}

	return g
}

func GenerateAPI(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := "grpc/" + string(file.GoPackageName) + "/server.go"
	g := gen.NewGeneratedFile(filename, protogen.GoImportPath("mosn.io/layotto/pkg/grpc/"+string(file.GoPackageName)))
	utils.AddHeader(g, true)

	g.P("package ", file.GoPackageName)
	g.P()

	componentPackageName := g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/components/" + file.GoPackageName).Ident(""))
	pbPackageName := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: file.GoImportPath,
	})
	g.P()

	for _, service := range file.Services {
		r := newRender(gen, file, g, service, pbPackageName[:len(pbPackageName)-1])
		r.ComponentPackageName = componentPackageName[:len(componentPackageName)-1]
		g.P(r.render())
	}

	return g
}
