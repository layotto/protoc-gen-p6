import (
    rawGRPC "google.golang.org/grpc"
    grpc_api "mosn.io/layotto/pkg/grpc"
)
{{$pb_name := .PackageName}}
{{$component_name := .ComponentPackageName }}

func NewAPI(ac *grpc_api.ApplicationContext) grpc_api.GrpcAPI {
    result:= &server{
		appId: ac.AppId,
		components: make(map[string]{{ .ComponentPackageName }}.{{.Name}}) ,
	}

    for k,v :=range ac.{{ .Extend.FieldNameInContext }} {
        comp, ok := v.({{ .ComponentPackageName }}.{{.Name}})
        if !ok {
            continue
        }
        // put it in the components map
        result.components[k]=comp
    }
	return result
}

type server struct {
	appId       string
	components  map[string]{{ .ComponentPackageName }}.{{.Name}}
}

{{range .MethodSet}}
func (s *server) {{.Name}}(ctx context.Context, in *{{$pb_name}}.{{.Request}}) (*{{$pb_name}}.{{.Reply}}, error){
	// find the component
	comp := s.components[in.ComponentName]
	if comp == nil {
		return nil, invalidArgumentError("{{.Name}}", grpc_api.ErrComponentNotFound, "{{$pb_name}}", in.ComponentName)
	}

	// convert request
	req := &{{$component_name}}.{{.Request}}{}
	err := copier.CopyWithOption(req, in, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Error when converting the request: %s", err.Error())
    }

	// delegate to the component
	resp, err := comp.{{.Name}}(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// convert response
	out := &{{$pb_name}}.{{.Reply}}{}
	err = copier.CopyWithOption(out, resp, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Error when converting the response: %s", err.Error())
    }
	return out, nil
}
{{end}}

func invalidArgumentError(method string, format string, a ...interface{}) error {
	err := status.Errorf(codes.InvalidArgument, format, a...)
	log.DefaultLogger.Errorf(fmt.Sprintf("%s fail: %+v", method, err))
	return err
}

func (s *server) Init(conn *rawGRPC.ClientConn) error {
	return nil
}

func (s *server) Register(rawGrpcServer *rawGRPC.Server) error {
	{{ .PackageName }}.Register{{.Name}}Server(rawGrpcServer, s)
	return nil
}
