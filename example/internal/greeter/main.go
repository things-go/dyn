package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/things-go/dyn/carry"
	"github.com/things-go/dyn/example/gen/hello"
	transportHttp "github.com/things-go/dyn/transport/http"
)

func main() {
	g := gin.Default()
	carrier := carry.NewCarry(carry.WithTranslatorData(translatorData{}))
	g.Use(transportHttp.CarrierInterceptor(carrier))
	g.Use(transportHttp.TransportInterceptor())
	group := g.Group("/api")
	hello.RegisterGreeterHTTPServer(group, new(Greeter))
	g.Run(":9090")
}

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type translatorData struct{}

func (t translatorData) TranslateData(v any) any {
	return &Result{
		Code:    200,
		Message: "ok",
		Data:    v,
	}
}

var _ hello.GreeterHTTPServer = (*Greeter)(nil)

type Greeter struct{}

// GetHello implements hello.GreeterHTTPServer.
func (g *Greeter) GetHello(_ context.Context, _ *hello.GetHelloRequest) (*hello.GetHelloReply, error) {
	return &hello.GetHelloReply{
		Message: "hello world",
	}, nil
}

// SayHello implements hello.GreeterHTTPServer.
func (g *Greeter) SayHello(_ context.Context, req *hello.HelloRequest) (*hello.HelloReply, error) {
	return &hello.HelloReply{
		Message: fmt.Sprintf("hello %s", req.Name),
	}, nil
}
