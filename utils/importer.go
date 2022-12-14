package utils

import (
	"google.golang.org/protobuf/compiler/protogen"
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

func AddHeader(g *protogen.GeneratedFile, canEdit bool) {
	header := "// Code generated by github.com/layotto/protoc-gen-p6 ."
	if canEdit {
		header = "// Code generated by github.com/layotto/protoc-gen-p6. DO NOT EDIT."
	}
	g.P(header)
	g.P()
	AddLicense(g)
}

func AddLicense(g *protogen.GeneratedFile) {
	g.P(`// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.`)
	g.P()
}
