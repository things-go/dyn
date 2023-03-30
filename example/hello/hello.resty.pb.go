// Code generated by protoc-gen-dyn-resty. DO NOT EDIT.
// versions:
// - protoc-gen-dyn-resty v0.1.0
// - protoc                v3.21.2
// source: hello.proto

package examples

import (
	context "context"
	errors "errors"
	http "github.com/things-go/dyn/transport/http"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = errors.New
var _ = context.TODO
var _ = http.NewClient

type GreeterHTTPClient interface {
	// SayHello Sends a hello
	// I am a trailing comment
	SayHello(ctx context.Context, req *HelloRequest, opts ...http.CallOption) (*HelloReply, error)
}

type GreeterHTTPClientImpl struct {
	cc *http.Client
}

func NewGreeterHTTPClient(c *http.Client) GreeterHTTPClient {
	return &GreeterHTTPClientImpl{
		cc: c,
	}
}

func (c *GreeterHTTPClientImpl) SayHello(ctx context.Context, req *HelloRequest, opts ...http.CallOption) (*HelloReply, error) {
	var err error
	var resp HelloReply

	settings := http.DefaultCallOption("/v1/hello")
	for _, opt := range opts {
		opt(&settings)
	}
	path := settings.Path
	ctx = http.WithValueCallOption(ctx, settings)
	err = c.cc.Invoke(ctx, "POST", path, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}