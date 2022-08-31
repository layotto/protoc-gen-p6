const (
	serviceName = "{{ .PackageName }}"
)

// Config is the component's configuration
type Config struct {
	ref.Config
	Type     string            `json:"type"`
	Metadata map[string]string `json:"metadata"`
}

type Registry interface {
	Register(fs ...*Factory)
	Create(compType string) ( {{ $.Name }} , error)
}

type Factory struct {
	CompType      string
	FactoryMethod func()  {{ $.Name }}
}

func NewFactory(compType string, f func()  {{ $.Name }} ) *Factory {
	return &Factory{
		CompType:      compType,
		FactoryMethod: f,
	}
}

type registry struct {
	stores  map[string]func()  {{ $.Name }}
	info *info.RuntimeInfo
}

func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(serviceName)
	return &registry{
		stores:  make(map[string]func()  {{ $.Name }} ),
		info: info,
	}
}

func (r *registry) Register(fs ...*Factory) {
	for _, f := range fs {
		r.stores[f.CompType] = f.FactoryMethod
		r.info.RegisterComponent(serviceName, f.CompType)
	}
}

func (r *registry) Create(compType string) ( {{ $.Name }} , error) {
	if f, ok := r.stores[compType]; ok {
		r.info.LoadComponent(serviceName, compType)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not registered", compType)
}