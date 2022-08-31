package main

import (
	"flag"
	"fmt"
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

type A struct {
	B
}
type B struct {
	xxx map[string]string
}

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
		services := make([]*protogen.File, 0)
		sdkRender := client.NewRender(gen)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			// 1. generate components code
			// struct
			path, name := utils.SplitDirectoryAndFilename(f.GeneratedFilenamePrefix)
			filename := path + "/components/" + name + ".struct.go"
			genstruct.GenerateStructFile(gen, f, filename)

			// interface
			// check there is any service
			if len(f.Services) == 0 {
				continue
			}
			filename = path + "/components/" + name + ".interface.go"
			gentemplate.GenerateComponentInterface(gen, f, filename)
			// types
			filename = path + "/components/" + name + ".types.go"
			gentemplate.GenerateComponentTypes(gen, f, filename)

			// 2. generate gRPC API plugin code
			gentemplate.GenerateAPI(gen, f)

			services = append(services, f)
		}
		// 3. ApplicationContext
		api.GenerateApplicationContext(gen, services)

		// 4. runtime related code
		runtime.GenerateConfig(gen, services)
		runtime.GenerateOptions(gen, services)
		runtime.GenerateNewApplicationContext(gen, services)
		runtime.GenerateComponentRelatedCode(gen, services)

		// 5. generate sdk
		return sdkRender.Render(services)
	})
}
