type {{ $.Name }} interface {
    Init(context.Context, *Config) error

{{range .MethodSet}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}