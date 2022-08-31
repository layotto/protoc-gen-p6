package gentemplate

import (
	"bytes"
	_ "embed"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
	"text/template"
)

//go:embed component.go.tpl
var componentTpl string

//go:embed component_types.go.tpl
var componentTypesTpl string

//go:embed server.go.tpl
var apiTpl string

func newRender(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, s *protogen.Service, packageName string) *render {
	// check deprecated
	if s.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}
	r := &render{
		Name:         s.GoName,
		FullName:     string(s.Desc.FullName()),
		FilePath:     file.Desc.Path(),
		PackageName:  packageName,
		GoImportPath: string(file.GoImportPath),
	}

	for _, method := range s.Methods {
		r.Methods = append(r.Methods, toMethod(method)...)
	}
	// add message for debug
	for _, method := range s.Methods {
		r.Messages = append(r.Messages, toMessage(method.Input))
	}
	return r
}

type render struct {
	Name     string // Greeter
	FullName string // helloworld.Greeter
	FilePath string // api/helloworld/helloworld.proto

	Methods              []*method
	Messages             []*message
	MessageSet           map[string]*protogen.Message
	MethodSet            map[string]*method
	DebugInfo            string
	PackageName          string
	ComponentPackageName string
	GoImportPath         string
}

func (s *render) render() string {
	return s.doRender("l8", apiTpl)
}

func (s *render) renderComponent() string {
	return s.doRender("component", componentTpl)
}

func (s *render) doRender(name string, tplName string) string {
	if s.MethodSet == nil {
		s.MethodSet = map[string]*method{}
		for _, m := range s.Methods {
			m := m
			s.MethodSet[m.Name] = m
		}
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New(name).Parse(strings.TrimSpace(tplName))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return buf.String()
}

// InterfaceName render interface name
//func (s *render) InterfaceName() string {
//	return s.Name + "XXX"
//}

type message struct {
	Name         string
	GoImportPath string
	Fields       []*field
}

type field struct {
	Name string
	Kind string
}

type method struct {
	Name    string // SayHello
	Num     int    // 一个 rpc 方法可以对应多个 http 请求
	Request string // SayHelloReq
	Reply   string // SayHelloResp
}

func toMessage(input *protogen.Message) *message {
	return &message{
		Name:         input.GoIdent.GoName,
		GoImportPath: string(input.GoIdent.GoImportPath),
		Fields:       toFields(input.Fields),
	}
}

func toFields(fields []*protogen.Field) []*field {
	r := make([]*field, len(fields))
	for _, protoField := range fields {
		r = append(r, &field{
			Name: protoField.GoName,
			Kind: string(protoField.Desc.FullName()),
		})
	}
	return r
}

func toMethod(m *protogen.Method) []*method {
	var methods []*method

	methods = append(methods, buildMethodDesc(m))
	return methods
}

func buildMethodDesc(m *protogen.Method) *method {
	defer func() { methodSets[m.GoName]++ }()
	md := &method{
		Name:    m.GoName,
		Num:     methodSets[m.GoName],
		Request: m.Input.GoIdent.GoName,
		Reply:   m.Output.GoIdent.GoName,
	}
	return md
}
