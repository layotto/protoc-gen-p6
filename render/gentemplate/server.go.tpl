import (
    rawGRPC "google.golang.org/grpc"
    grpc_api "mosn.io/layotto/pkg/grpc"
)

func NewAPI(ac *grpc_api.ApplicationContext) grpc_api.GrpcAPI {
	return &server{
		appId: ac.AppId,
		components: ac.{{.Name}},
	}
}

type server struct {
	appId       string
	components  map[string]{{ .ComponentPackageName }}.{{.Name}}
}


func (s *server) Init(conn *rawGRPC.ClientConn) error {
	return nil
}

func (s *server) Register(rawGrpcServer *rawGRPC.Server) error {
	{{ .PackageName }}.Register{{.Name}}Server(rawGrpcServer, s)
	return nil
}
