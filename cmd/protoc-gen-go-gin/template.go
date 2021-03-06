package main

import (
	"bytes"
	"strings"
	"text/template"
)

var httpTemplate = `
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
{{$useCustomResp := .UseCustomResponse}}
{{$rpcMode := .RpcMode}}
{{$allowFromAPI := .AllowFromAPI}}
type {{$svrType}}HTTPServer interface {
{{- range .MethodSets}}
{{- if eq $rpcMode "rpcx"}}
	{{.Name}}(context.Context, *{{.Request}}, *{{.Reply}}) error
{{- else}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
{{- end}}
    Validate(context.Context, any) error
	ErrorEncoder(c *gin.Context, err error, isBadRequest bool)
{{- if $useCustomResp}}
	ResponseEncoder(c *gin.Context, v any)
{{- end}}
}

type Unimplemented{{$svrType}}HTTPServer struct {}

{{- range .MethodSets}}
{{- if eq $rpcMode "rpcx"}}
func (*Unimplemented{{$svrType}}HTTPServer) {{.Name}}(context.Context, *{{.Request}}, *{{.Reply}}) error {
	return errors.New("method {{.Name}} not implemented")
}
{{- else}}
func (*Unimplemented{{$svrType}}HTTPServer) {{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error) {
	return nil, errors.New("method {{.Name}} not implemented")
}
{{- end}}
{{- end}}
func (*Unimplemented{{$svrType}}HTTPServer) Validate(context.Context, any) error { return nil }
func (*Unimplemented{{$svrType}}HTTPServer) ErrorEncoder(c *gin.Context, err error, isBadRequest bool) {
	var code = 500
	if isBadRequest {
		code = 400
	}
	c.String(code, err.Error())
}
{{- if $useCustomResp}}
func (*Unimplemented{{$svrType}}HTTPServer) ResponseEncoder(c *gin.Context, v any) {
	c.JSON(200, v)
}
{{- end}}

func Register{{$svrType}}HTTPServer(g *gin.RouterGroup, srv {{$svrType}}HTTPServer) {
	r := g.Group("")
	{{- range .Methods}}
	r.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
	{{- end}}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv {{$svrType}}HTTPServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		shouldBind := func(req any) error {
			{{- if .HasBody}}
			if err := c.ShouldBind(req); err != nil {
				return err
			}
			{{- if not (eq .Body "")}}
			if err := c.ShouldBindQuery(req); err != nil {
				return err
			}
			{{- end}}
			{{- else}}
			{{- if not (eq .Method "PATCH")}}
			if err := c.ShouldBindQuery(req); err != nil {
				return err
			}
			{{- end}}
			{{- end}}
			{{- if .HasVars}}
			if err := c.ShouldBindUri(req); err != nil {
				return err
			}
			{{- end}}
			return srv.Validate(c.Request.Context(), req)
		}

		var err error
		var req {{.Request}}
		{{- if eq $rpcMode "rpcx"}}
		var reply {{.Reply}}
		{{- else}}
		var reply *{{.Reply}}
		{{- end}}
		if err = shouldBind(&req); err != nil {
			srv.ErrorEncoder(c, err, true)
			return
		}
		{{- if eq $rpcMode "rpcx"}}
		err = srv.{{.Name}}(c.Request.Context(), &req, &reply)
		{{- else}}
		reply, err = srv.{{.Name}}(c.Request.Context(), &req)
		{{- end}}
		if err != nil {
			srv.ErrorEncoder(c, err, false)
			return
		}
		{{- if eq $rpcMode "rpcx"}}
		{{- if $useCustomResp}}
		srv.ResponseEncoder(c, &reply)
		{{- else}}
		c.JSON(200, &reply)
		{{- end}}
		{{- else}}
		{{- if $useCustomResp}}
		srv.ResponseEncoder(c, reply)
		{{- else}}
		c.JSON(200, reply)		
		{{- end}}
		{{- end}}
	}
}
{{end}}

{{- if $allowFromAPI}}
type From{{$svrType}}HTTPServer interface {
{{- range .MethodSets}}
{{- if eq $rpcMode "rpcx"}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- else}}
	{{.Name}}(context.Context, *{{.Request}}, *{{.Reply}}) error
{{- end}}
{{- end}}
    Validate(context.Context, any) error
	ErrorEncoder(c *gin.Context, err error, isBadRequest bool)
{{- if $useCustomResp}}
	ResponseEncoder(c *gin.Context, v any)
{{- end}}
}

type From{{$svrType}} struct {
	From{{$svrType}}HTTPServer
}

func NewFrom{{$svrType}}HTTPServer(from From{{$svrType}}HTTPServer) {{$svrType}}HTTPServer {
	return &From{{$svrType}}{from}
}

{{- range .MethodSets}}
{{- if eq $rpcMode "rpcx"}}
func (f *From{{$svrType}}) {{.Name}}(ctx context.Context, req *{{.Request}}, reply *{{.Reply}}) error {
	result, err := f.From{{$svrType}}HTTPServer.{{.Name}}(ctx, req)
	if err != nil {
		return err
	}
	if result == nil {
		*reply = {{.Reply}}{}
	} else {
		*reply = *result
	}
	return nil
}
{{- else}}
func (f *From{{$svrType}}) {{.Name}}(ctx context.Context, req *{{.Request}}) (*{{.Reply}}, error) {
	var err error 
	var reply {{.Reply}}

	err = f.From{{$svrType}}HTTPServer.{{.Name}}(ctx, req, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}
{{- end}}
{{- end}}
{{- end}}
`

type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/v1/helloworld.proto
	Methods     []*methodDesc
	MethodSets  map[string]*methodDesc

	UseCustomResponse bool
	RpcMode           string
	AllowFromAPI      bool
}

type methodDesc struct {
	// method
	Name    string // ?????????
	Num     int    // ?????????
	Request string // ????????????
	Reply   string // ????????????
	// http_rule
	Path         string // ??????
	Method       string // ??????
	HasVars      bool   // ?????????url??????
	HasBody      bool   // ??????????????????
	Body         string // ?????????
	ResponseBody string //
}

func (s *serviceDesc) execute() string {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("gin").Parse(strings.TrimSpace(httpTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}
