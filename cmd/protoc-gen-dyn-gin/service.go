package main

import (
	"strconv"

	"google.golang.org/protobuf/compiler/protogen"
)

type serviceDesc struct {
	Deprecated      bool   // deprecated or not
	ServiceType     string // Greeter
	ServiceName     string // helloworld.Greeter
	Metadata        string // api/v1/helloworld.proto
	LeadingComment  string // leading comment
	TrailingComment string // trailing comment
	Comment         string // combine leading and trailing comment
	Methods         []*methodDesc

	UseEncoding bool
}

type methodDesc struct {
	Deprecated bool // deprecated or not
	// method
	Name            string // 方法名
	Num             int    // 方法号
	Request         string // 请求结构
	Reply           string // 回复结构
	LeadingComment  string // leading comment
	TrailingComment string // trailing comment
	Comment         string // combine leading and trailing comment

	// http_rule
	Path         string // 路径
	Method       string // 方法
	HasVars      bool   // 是否有url参数
	HasBody      bool   // 是否有消息体
	Body         string // 请求消息体
	ResponseBody string // 回复消息体
}

func executeServiceDesc(g *protogen.GeneratedFile, s *serviceDesc) error {
	if args.EnableMetadata {
		g.P("const ", serviceTypeMetadataKey(s.ServiceType), " = \"", serviceTypeMetadataValue(s.ServiceType, s.LeadingComment)+"\"")
	}
	g.P()
	methodSets := make(map[string]struct{})
	// http interface defined
	if s.Deprecated {
		g.P(deprecationComment)
	}
	g.P("// ", serverInterfaceName(s.ServiceType), " ", s.Comment)
	g.P("type ", serverInterfaceName(s.ServiceType), " interface {")
	for _, m := range s.Methods {
		_, ok := methodSets[m.Name]
		if ok { // unique because additional_bindings
			continue
		}
		methodSets[m.Name] = struct{}{}
		if m.Deprecated {
			g.P(deprecationComment)
		}
		g.P(m.Comment)
		g.P(serverMethodName(g, m))
	}
	g.P("}")
	g.P()
	// register http server handler
	if s.Deprecated {
		g.P(deprecationComment)
	}
	g.P("func Register", s.ServiceType, "HTTPServer(g *", g.QualifiedGoIdent(ginPackage.Ident("RouterGroup")), ", srv ", serverInterfaceName(s.ServiceType), ") {")
	g.P(`r := g.Group("")`)
	g.P("{")
	for _, m := range s.Methods {
		useMdMiddleware := ""
		if args.EnableMetadata {
			useMdMiddleware = "" +
				g.QualifiedGoIdent(transportHttpPackage.Ident("MetadataInterceptor")) +
				"(" +
				g.QualifiedGoIdent(transportHttpPackage.Ident("Metadata")) +
				"{Service: " + serviceTypeMetadataKey(s.ServiceType) + ", Method: \"" + methodMetadataValue(m.Name, m.LeadingComment) + "\"}" +
				"), "
		}
		g.P("r.", m.Method, `("`, m.Path, `", `, useMdMiddleware, serverHandlerMethodName(s.ServiceType, m), "(srv))")

	}
	g.P("}")
	g.P("}")
	g.P()
	// handler
	for _, m := range s.Methods {
		if m.Deprecated {
			g.P(deprecationComment)
		}
		g.P("func ", serverHandlerMethodName(s.ServiceType, m), "(srv ", s.ServiceType, "HTTPServer", ") ", g.QualifiedGoIdent(ginPackage.Ident("HandlerFunc")), " {")
		{ // gin.HandleFunc closure
			g.P("return func(c *", g.QualifiedGoIdent(ginPackage.Ident("Context")), ") {")
			g.P("carrier := ", g.QualifiedGoIdent(transportHttpPackage.Ident("FromCarrier")), "(c.Request.Context())")
			{ // binding
				g.P("shouldBind := func(req *", m.Request, ") error {")
				if s.UseEncoding {
					if m.HasBody {
						g.P("if err := carrier.Bind(c, req", m.Body, "); err != nil {")
						g.P("return err")
						g.P("}")
						if m.Body != "" {
							g.P("if err := carrier.BindQuery(c, req); err != nil {")
							g.P("return err")
							g.P("}")
						}
					} else {
						if m.Method != "PATCH" {
							g.P("if err := carrier.BindQuery(c, req", m.Body, "); err != nil {")
							g.P("return err")
							g.P("}")
						}
					}
					if m.HasVars {
						g.P("if err := carrier.BindUri(c, req); err != nil {")
						g.P("return err")
						g.P("}")
					}
				} else {
					if m.HasBody {
						g.P("if err := c.ShouldBind(req", m.Body, "); err != nil {")
						g.P("return err")
						g.P("}")
						if m.Body != "" {
							g.P("if err := c.ShouldBindQuery(req); err != nil {")
							g.P("return err")
							g.P("}")
						}
					} else {
						if m.Method != "PATCH" {
							g.P("if err := c.ShouldBindQuery(req", m.Body, "); err != nil {")
							g.P("return err")
							g.P("}")
						}
					}
					if m.HasVars {
						g.P("if err := c.ShouldBindUri(req); err != nil {")
						g.P("return err")
						g.P("}")
					}
				}
				g.P("return carrier.Validate(c.Request.Context(), req)")
				g.P("}")
			}
			g.P()
			{ // done
				g.P("var err error")
				g.P("var req ", m.Request)
				g.P("var reply *", m.Reply)
				g.P()
				g.P("if err = shouldBind(&req); err != nil {")
				g.P("carrier.Error(c, err)")
				g.P("return")
				g.P("}")
				g.P("reply, err = srv.", m.Name, "(c.Request.Context(), &req)")
				g.P("if err != nil {")
				g.P("carrier.Error(c, err)")
				g.P("return")
				g.P("}")
				g.P("carrier.Render(c, reply", m.ResponseBody, ")")
			}
			g.P("}")
		}
		g.P("}")
		g.P()
	}

	return nil
}

func serviceTypeMetadataKey(serverType string) string {
	return "__" + serverType + "_Metadata_Service"
}
func serviceTypeMetadataValue(serverType, leadingComment string) string {
	if leadingComment != "" {
		return lineComment(leadingComment)
	} else {
		return serverType
	}
}

func methodMetadataValue(name, leadingComment string) string {
	if leadingComment != "" {
		return lineComment(leadingComment)
	} else {
		return name
	}
}

func serverInterfaceName(serverType string) string {
	return serverType + "HTTPServer"
}

func serverMethodName(g *protogen.GeneratedFile, m *methodDesc) string {
	return m.Name + "(" + g.QualifiedGoIdent(contextPackage.Ident("Context")) + ", *" + m.Request + ") (*" + m.Reply + ", error)"
}

func serverHandlerMethodName(serverType string, m *methodDesc) string {
	return "_" + serverType + "_" + m.Name + strconv.Itoa(m.Num) + "_HTTP_Handler"
}
