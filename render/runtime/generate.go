package runtime

import (
	"fmt"
	"github.com/seeflood/protoc-gen-p6/utils"
	"google.golang.org/protobuf/compiler/protogen"
)

func GenerateExtensionComponentConfig(gen *protogen.Plugin, files []*protogen.File) *protogen.GeneratedFile {
	filename := "runtime/config_generated.go"
	g := gen.NewGeneratedFile(filename, files[0].GoImportPath)

	// 1. headers
	utils.AddHeader(g, false)
	g.P()
	g.P("package ", "runtime")
	g.P()

	// 2. ExtensionComponentConfig
	g.P(`type ExtensionComponentConfig struct {`)

	for _, file := range files {
		for _, service := range file.Services {
			// add import
			importPath := "mosn.io/layotto/components/" + string(file.GoPackageName)
			pkgName := g.QualifiedGoIdent(protogen.GoImportPath(importPath).Ident(""))
			// add comment
			g.P("// ", file.GoImportPath)
			g.P("// ", pkgName)

			// e.g. 	BlogService map[string]blog.BlogService `json:"blog"`
			g.P(fmt.Sprintf("\t%s\tmap[string]%sConfig\t`json:\"%s\"`", service.GoName, pkgName, pkgName[:len(pkgName)-1]))
			g.P()
		}
	}
	// the end
	g.P("}")

	return g
}

func GenerateOptions(gen *protogen.Plugin, files []*protogen.File) *protogen.GeneratedFile {
	filename := "runtime/options_generated.go"
	g := gen.NewGeneratedFile(filename, files[0].GoImportPath)

	// 1. headers
	utils.AddHeader(g, false)
	g.P()
	g.P("package ", "runtime")
	g.P()

	// 2. extensionComponentFactorys
	g.P(`type extensionComponentFactorys struct {`)

	for _, file := range files {
		for _ = range file.Services {
			// add import
			importPath := "mosn.io/layotto/components/" + string(file.GoPackageName)
			pkgName := g.QualifiedGoIdent(protogen.GoImportPath(importPath).Ident(""))
			// add comment
			g.P("// ", file.GoImportPath)
			g.P("// ", pkgName)

			nickName := pkgName[:len(pkgName)-1]
			// e.g. 	oss           []*oss.Factory
			g.P(fmt.Sprintf("\t%s\t[]*%sFactory", nickName, pkgName))
			g.P()
		}
	}
	// the end of extensionComponentFactorys
	g.P("}")

	// 3. with function
	// for example:
	// func WithOssFactory(oss ...*oss.Factory) Option {
	//	return func(o *runtimeOptions) {
	//		o.services.oss = append(o.services.oss, oss...)
	//	}
	//}
	for _, file := range files {
		for _, service := range file.Services {
			// add import
			pkgName := utils.ImportComponent(g, file.GoPackageName)

			nickName := pkgName[:len(pkgName)-1]
			//  func WithOssFactory(oss ...*oss.Factory) Option {
			g.P(fmt.Sprintf("func With%sFactory(%s ...*%s.Factory) Option {", service.GoName, nickName, nickName))
			// 	return func(o *runtimeOptions) {
			g.P("\treturn func(o *runtimeOptions) {")
			//		o.services.oss = append(o.services.oss, oss...)
			g.P(fmt.Sprintf("\t\to.services.%s = append(o.services.%s, %s...)", nickName, nickName, nickName))
			// the end
			g.P(`	}
}`)
			g.P()
		}
	}

	// 4. WithExtensionGrpcAPI
	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/pkg/grpc/extension/s3").Ident(""))
	g.P(`func WithExtensionGrpcAPI() Option {
	return WithGrpcAPI(
		s3.NewS3Server,`)
	for _, file := range files {
		for _ = range file.Services {
			// add import
			pkgName := utils.ImportAPI(g, file.GoPackageName)

			g.P(pkgName, "NewAPI,")
		}
	}
	// the end of WithExtensionGrpcAPI
	g.P(`	)
}`)
	return g
}

func GenerateNewApplicationContext(gen *protogen.Plugin, files []*protogen.File) *protogen.GeneratedFile {
	filename := "runtime/context_generated.go"
	g := gen.NewGeneratedFile(filename, files[0].GoImportPath)

	// 1. headers
	utils.AddHeader(g, false)
	g.P()
	g.P("package ", "runtime")
	g.P()

	g.QualifiedGoIdent(protogen.GoImportPath("mosn.io/layotto/pkg/grpc").Ident("grpc"))

	// 2. newApplicationContext
	g.P(`func newApplicationContext(m *MosnRuntime) *grpc.ApplicationContext {
	return &grpc.ApplicationContext{
		AppId:                 m.runtimeConfig.AppManagement.AppId,
		Hellos:                m.hellos,
		ConfigStores:          m.configStores,
		Rpcs:                  m.rpcs,
		PubSubs:               m.pubSubs,
		StateStores:           m.states,
		Files:                 m.files,
		Oss:                   m.oss,
		LockStores:            m.locks,
		Sequencers:            m.sequencers,
		SendToOutputBindingFn: m.sendToOutputBinding,
		SecretStores:          m.secretStores,
		DynamicComponents:     m.dynamicComponents,
		CustomComponent:       m.customComponent,`)

	// e.g.
	// 			LockStores: m.locks,
	for _, file := range files {
		for _, service := range file.Services {
			// e.g.
			// BlogService : m.blogService
			g.P(fmt.Sprintf("\t%s\t:\tm.%s,", service.GoName, utils.LowerCammel(service.GoName)))
			g.P()
		}
	}
	// the end of newApplicationContext
	g.P(`	}
}
`)

	return g
}

func GenerateComponentRelatedCode(gen *protogen.Plugin, files []*protogen.File) *protogen.GeneratedFile {
	filename := "runtime/component_generated.go"
	g := gen.NewGeneratedFile(filename, files[0].GoImportPath)

	// 1. headers
	utils.AddHeader(g, false)
	g.P()
	g.P("package ", "runtime")
	g.P()

	// 2. extensionComponents
	g.P(`type extensionComponents struct {`)

	// e.g.
	// 				blogService map[string]blog.BlogService
	for _, file := range files {
		for _, service := range file.Services {
			pkgName := utils.ImportComponent(g, file.GoPackageName)

			g.P(fmt.Sprintf("\t%s\tmap[string]%s%s", utils.LowerCammel(service.GoName), pkgName, service.GoName))
			g.P()
		}
	}
	// the end of newApplicationContext
	g.P(`}`)
	g.P()

	// 3. newExtensionComponents
	g.P(`func newExtensionComponents() *extensionComponents {
	return &extensionComponents{`)
	// e.g.
	// blogService: make(map[string]blog.BlogService),
	for _, file := range files {
		for _, service := range file.Services {
			pkgName := utils.ImportComponent(g, file.GoPackageName)

			g.P(fmt.Sprintf("\t%s:\tmake(map[string]%s%s),", utils.LowerCammel(service.GoName), pkgName, service.GoName))
			g.P()
		}
	}
	g.P(`	}
}`)
	g.P()

	// 4. initXXX function
	for _, file := range files {
		for _, service := range file.Services {
			pkgName := utils.ImportComponent(g, file.GoPackageName)
			pkgName = pkgName[:len(pkgName)-1]
			g.P(fmt.Sprintf(`func (m *MosnRuntime) init%s(factorys ...*%s.Factory) error {
	log.DefaultLogger.Infof("[runtime] init %s")

	// 1. register all implementation
	reg := %s.NewRegistry(m.info)
	reg.Register(factorys...)
	// 2. loop initializing`, service.GoName, pkgName, service.GoName, pkgName))

			g.P(fmt.Sprintf(`	for name, config := range m.runtimeConfig.%s {
		// 2.1. create the component
		c, err := reg.Create(config.Type)
		if err != nil {`, service.GoName))
			g.P("m.errInt(err, \"create the component %s failed\", name)")
			g.P(`			return err
		}
		//inject secret to component
		if config.Metadata, err = m.Injector.InjectSecretRef(config.SecretRef, config.Metadata); err != nil {
			return err
		}`)

			g.P(`// 2.2. init
		if err := c.Init(context.TODO(), &config); err != nil {`)
			g.P("m.errInt(err, \"init the component %s failed\", name)")
			g.P(`return err
		}`)
			g.P(fmt.Sprintf(`		m.%s[name] = c
	}
	return nil
}`, utils.LowerCammel(service.GoName)))
			g.P()
		}
	}

	// 5. initExtensionComponent
	g.P("func (m *MosnRuntime) initExtensionComponent(s services) error {")

	for _, file := range files {
		for _, service := range file.Services {
			pkgName := utils.ImportComponent(g, file.GoPackageName)
			pkgName = pkgName[:len(pkgName)-1]

			g.P(fmt.Sprintf(`	if err := m.init%s(s.%s...); err != nil {
		return err
	}`, service.GoName, pkgName))
			g.P()
		}
	}

	g.P(`	return nil
}`)
	return g
}
