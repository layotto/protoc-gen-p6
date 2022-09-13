package main

import (
	"flag"
	"fmt"
	mode "github.com/seeflood/protoc-gen-p6/mode"
	"github.com/seeflood/protoc-gen-p6/render/api"
	"github.com/seeflood/protoc-gen-p6/render/client"
	"github.com/seeflood/protoc-gen-p6/render/genstruct"
	"github.com/seeflood/protoc-gen-p6/render/gentemplate"
	"github.com/seeflood/protoc-gen-p6/render/runtime"
	"github.com/seeflood/protoc-gen-p6/utils"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "0.0.1"

func main() {
	// check version flag
	showVersion := flag.Bool("version", false, "show version")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-p6 %v\n", version)
		return
	}

	// render
	var flags flag.FlagSet

	options := protogen.Options{
		ParamFunc: flags.Set,
	}

	options.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		components := make([]*protogen.File, 0)
		apis := make([]*protogen.File, 0)
		needSDK := make([]*protogen.File, 0)

		sdkRender := client.NewRender(gen)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			// 1. generate components code
			// 1.1. generate struct
			_, name := utils.SplitDirectoryAndFilename(f.GeneratedFilenamePrefix)
			dir := "components/" + name + "/"
			filename := dir + "/struct_generated.go"
			genstruct.GenerateStructFile(gen, f, filename)

			// 1.2. generate interface
			// check there is any service
			if len(f.Services) == 0 {
				continue
			}
			// Note: we allow only 1 service in a file
			if len(f.Services) > 1 {
				panic("There are more than one service in a file.")
			}
			// 1.3. check mode
			comments := make([]string, 0)
			for _, comment := range f.Services[0].Comments.LeadingDetached {
				comments = append(comments, comment.String())
			}
			comments = append(comments, f.Services[0].Comments.Leading.String())
			compileMode, args := mode.CheckMode(comments)

			filename = dir + "/interface_generated.go"
			gentemplate.GenerateComponentInterface(gen, f, filename)

			if compileMode == mode.Independent {
				// generate types
				filename = dir + "/types_generated.go"
				gentemplate.GenerateComponentTypes(gen, f, filename)

				components = append(components, f)

				// generate gRPC API plugin code
				gentemplate.GenerateAPI(gen, f)
			} else if compileMode == mode.Extend {
				// TODO multiple extends
				// generate gRPC API plugin code
				gentemplate.GenerateExtendAPI(gen, f, args[0])
			}

			// 2. collect api plugins
			apis = append(apis, f)

			if mode.NeedGenerateSDK(comments) {
				needSDK = append(needSDK, f)
			}
		}
		if len(components) > 0 {
			// 3. generate ApplicationContext
			api.GenerateApplicationContext(gen, components)

			// 4. runtime related code
			runtime.GenerateExtensionComponentConfig(gen, components)
			runtime.GenerateOptions(gen, components)
			runtime.GenerateNewApplicationContext(gen, components)
			runtime.GenerateComponentRelatedCode(gen, components)
		}

		runtime.GenerateWithExtensionGrpcAPI(gen, apis)

		// 5. generate golang sdk
		return sdkRender.Render(needSDK)
	})
}
