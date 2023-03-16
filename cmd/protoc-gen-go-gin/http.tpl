{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
{{$useCustomResp := .UseCustomResponse}}
{{$useEncoding := .UseEncoding}}
type {{$svrType}}HTTPServer interface {
{{- range .MethodSets}}
	{{.Comment}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
	// Validate the request.
    Validate(context.Context, any) error
	// ErrorEncoder encode error response.
	ErrorEncoder(c *gin.Context, err error, isBadRequest bool)
{{- if $useEncoding}}
    // Bind checks the Method and Content-Type to select codec.Marshaler automatically,
    // Depending on the "Content-Type" header different bind are used.
    Bind(c *gin.Context, v any) error
    // BindQuery binds the passed struct pointer using the query codec.Marshaler.
    BindQuery(c *gin.Context, v any) error
    // BindUri binds the passed struct pointer using the uri codec.Marshaler.
    // NOTE: before use this, you should set uri params in the request context with RequestWithUri.
    BindUri(c *gin.Context, v any) error
    // RequestWithUri sets the URL params for the given request.
    RequestWithUri(req *http.Request, params gin.Params) *http.Request
    // Render encode response.
    Render(c *gin.Context, v any)
{{- else}}
{{- if $useCustomResp}}
	// Render encode response.
	Render(c *gin.Context, v any)
{{- end}}
{{- end}}
}

func Register{{$svrType}}HTTPServer(g *gin.RouterGroup, srv {{$svrType}}HTTPServer) {
	r := g.Group("")
	{{- range .Methods}}
	r.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
	{{- end}}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv {{$svrType}}HTTPServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		{{- if and $useEncoding .HasVars}}
		c.Request = srv.RequestWithUri(c.Request, c.Params)
		{{- end}}
		shouldBind := func(req *{{.Request}}) error {
		    {{- if $useEncoding}}
		    {{- if .HasBody}}
			if err := srv.Bind(c, req{{.Body}}); err != nil {
				return err
			}
			{{- if not (eq .Body "")}}
			if err := srv.BindQuery(c, req); err != nil {
				return err
			}
			{{- end}}
			{{- else}}
			{{- if not (eq .Method "PATCH")}}
			if err := srv.BindQuery(c, req{{.Body}}); err != nil {
				return err
			}
			{{- end}}
			{{- end}}
			{{- if .HasVars}}
			if err := srv.BindUri(c, req); err != nil {
				return err
			}
			{{- end}}
		    {{- else}}
			{{- if .HasBody}}
			if err := c.ShouldBind(req{{.Body}}); err != nil {
				return err
			}
			{{- if not (eq .Body "")}}
			if err := c.ShouldBindQuery(req); err != nil {
				return err
			}
			{{- end}}
			{{- else}}
			{{- if not (eq .Method "PATCH")}}
			if err := c.ShouldBindQuery(req{{.Body}}); err != nil {
				return err
			}
			{{- end}}
			{{- end}}
			{{- if .HasVars}}
			if err := c.ShouldBindUri(req); err != nil {
				return err
			}
			{{- end}}
			{{- end}}
			return srv.Validate(c.Request.Context(), req)
		}

		var err error
		var req {{.Request}}
        var reply *{{.Reply}}

		if err = shouldBind(&req); err != nil {
			srv.ErrorEncoder(c, err, true)
			return
		}
		reply, err = srv.{{.Name}}(c.Request.Context(), &req)
		if err != nil {
			srv.ErrorEncoder(c, err, false)
			return
		}
		{{- if or $useEncoding $useCustomResp}}
		srv.Render(c, reply{{.ResponseBody}})
        {{- else}}
        c.JSON(200, reply{{.ResponseBody}})
        {{- end}}
	}
}
{{end}}