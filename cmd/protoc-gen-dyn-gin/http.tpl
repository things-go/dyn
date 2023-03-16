{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
type {{$svrType}}HTTPServer interface {
{{- range .MethodSets}}
	{{.Comment}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}

func Register{{$svrType}}HTTPServer(g *gin.RouterGroup, srv {{$svrType}}HTTPServer) {
	r := g.Group("")
	{
	{{- range .Methods}}
		r.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
	{{- end}}
	}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv {{$svrType}}HTTPServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		carrier := http.FromConvey(c.Request.Context())
		{{- if .HasVars}}
		c.Request = carrier.RequestWithUri(c.Request, c.Params)
		{{- end}}
		shouldBind := func(req *{{.Request}}) error {
		    {{- if .HasBody}}
			if err := carrier.Bind(c, req{{.Body}}); err != nil {
				return err
			}
			{{- if not (eq .Body "")}}
			if err := carrier.BindQuery(c, req); err != nil {
				return err
			}
			{{- end}}
			{{- else}}
			{{- if not (eq .Method "PATCH")}}
			if err := carrier.BindQuery(c, req{{.Body}}); err != nil {
				return err
			}
			{{- end}}
			{{- end}}
			{{- if .HasVars}}
			if err := carrier.BindUri(c, req); err != nil {
				return err
			}
			{{- end}}
			return carrier.Validate(c.Request.Context(), req)
		}

		var err error
		var req {{.Request}}
        var reply *{{.Reply}}

		if err = shouldBind(&req); err != nil {
			carrier.ErrorBadRequest(c, err)
			return
		}
		reply, err = srv.{{.Name}}(c.Request.Context(), &req)
		if err != nil {
			carrier.Error(c, err)
			return
		}
		carrier.Render(c, reply{{.ResponseBody}})
	}
}
{{end}}