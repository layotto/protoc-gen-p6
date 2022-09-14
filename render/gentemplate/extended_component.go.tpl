type {{ $.Name }} interface {

{{range .MethodSet}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}