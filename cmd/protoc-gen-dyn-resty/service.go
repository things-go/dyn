package main

import (
	"strconv"

	"google.golang.org/protobuf/compiler/protogen"
)

type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/v1/helloworld.proto
	Methods     []*methodDesc
}

type methodDesc struct {
	// method
	Name    string // 方法名
	Num     int    // 方法号
	Request string // 请求结构
	Reply   string // 回复结构
	Comment string // 方法注释
	// http_rule
	Path         string // 路径
	Method       string // 方法
	HasVars      bool   // 是否有url参数
	HasBody      bool   // 是否有消息体
	Body         string // 请求消息体
	ResponseBody string // 回复消息体
}

func executeServiceDesc(g *protogen.GeneratedFile, s *serviceDesc) error {
	// http interface defined
	g.P("type ", s.ServiceType, "HTTPClient", " interface {")
	for _, m := range s.Methods {
		g.P(m.Comment)
		g.P(clientSignature(g, m))
	}
	g.P("}")
	g.P()

	// http client implement.
	g.P("type ", s.ServiceType, "HTTPClientImpl struct {")
	g.P("cc *", g.QualifiedGoIdent(transportHttpPackage.Ident("Client")))
	g.P("}")
	g.P()
	// http client factory method.
	g.P("func New", s.ServiceType, "HTTPClient(c *", g.QualifiedGoIdent(transportHttpPackage.Ident("Client")), ") ", s.ServiceType, "HTTPClient {")
	g.P("return &", s.ServiceType, "HTTPClientImpl{")
	g.P("cc: c,")
	g.P("}")
	g.P("}")
	g.P()

	// http client implement methods.
	for _, m := range s.Methods {
		g.P("func (c *", s.ServiceType, "HTTPClientImpl)", clientSignature(g, m), " {")
		g.P("var err error")
		g.P("var resp ", m.Reply)
		g.P()
		g.P("settings := ", g.QualifiedGoIdent(transportHttpPackage.Ident("DefaultCallOption")), `("`, m.Path, `")`)
		g.P("for _, opt := range opts {")
		g.P("opt(&settings)")
		g.P("}")

		if m.HasVars {
			g.P("path := c.cc.EncodeURL(settings.Path, req, {{not .HasBody}})")
		} else {
			if m.HasBody {
				g.P("path := settings.Path")
			} else {
				g.P("var query string")
				g.P()
				g.P("query, err = c.cc.EncodeQuery(req)")
				g.P("if err != nil {")
				g.P("return nil, err")
				g.P("}")
				g.P("path := settings.Path")
				g.P(`if query != "" {`)
				g.P(`path += "?" + query`)
				g.P("}")
			}
		}
		g.P("ctx = ", g.QualifiedGoIdent(transportHttpPackage.Ident("WithValueCallOption")), "(ctx, settings)")

		reqValue := "nil"
		if m.HasBody {
			reqValue = "req" + m.Body
		}
		g.P(`err = c.cc.Invoke(ctx, "`, m.Method, `", path, `, reqValue, ", &resp", m.ResponseBody, ")")
		g.P("if err != nil {")
		g.P("return nil, err")
		g.P("}")
		g.P("return &resp, nil")
		g.P("}")
	}

	return nil
}

func clientSignature(g *protogen.GeneratedFile, m *methodDesc) string {
	num := ""
	if m.Num != 0 { // unique because additional_bindings
		num = "_" + strconv.Itoa(m.Num)
	}

	return m.Name + num + "(ctx " + g.QualifiedGoIdent(contextPackage.Ident("Context")) +
		", req *" + m.Request + ", opts ..." + g.QualifiedGoIdent(transportHttpPackage.Ident("CallOption")) +
		") (*" + m.Reply + ", error)"
}
