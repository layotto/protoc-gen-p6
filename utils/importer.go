package utils

import (
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func ImportComponent(g *protogen.GeneratedFile, goPackageName protogen.GoPackageName) string {
	importPath := "mosn.io/layotto/components/" + string(goPackageName)
	pkgName := g.QualifiedGoIdent(protogen.GoImportPath(importPath).Ident(""))
	return pkgName
}

func ImportAPI(g *protogen.GeneratedFile, goPackageName protogen.GoPackageName) string {
	importPath := "mosn.io/layotto/pkg/grpc/" + string(goPackageName)
	pkgName := g.QualifiedGoIdent(protogen.GoImportPath(importPath).Ident(""))
	return pkgName
}

func LowerCammel(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}
