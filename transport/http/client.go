package http

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/things-go/encoding"
	"golang.org/x/oauth2"
)

var noRequestBodyMethods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodHead:    {},
	http.MethodDelete:  {},
	http.MethodConnect: {},
	http.MethodOptions: {},
	http.MethodTrace:   {},
}

type Client struct {
	cc    *resty.Client
	codec *encoding.Encoding
	// A TokenSource is anything that can return a token.
	tokenSource oauth2.TokenSource
	// validate request
	validate func(any) error
	// call option
	callOptions []CallOption
}

type ClientOption func(*Client)

func WithEncoding(codec *encoding.Encoding) ClientOption {
	return func(c *Client) {
		c.codec = codec
	}
}

func WithTokenSource(t oauth2.TokenSource) ClientOption {
	return func(c *Client) {
		c.tokenSource = t
	}
}

// Deprecated: use WithCallOption(WithCoNoAuth) instead.
func WithNoAuth() ClientOption {
	return func(c *Client) {
		c.callOptions = append(c.callOptions, WithCoNoAuth())
	}
}

func WithValidate(f func(any) error) ClientOption {
	return func(c *Client) {
		c.validate = f
	}
}

func WithCallOption(cs ...CallOption) ClientOption {
	return func(c *Client) {
		c.callOptions = append(c.callOptions, cs...)
	}
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		cc:    resty.New(),
		codec: encoding.New(),
	}
	for _, opt := range opts {
		opt(c)
	}
	c.cc.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		if r.RawResponse != nil {
			body := r.RawResponse.Body
			defer body.Close()
			r.RawResponse.Body = io.NopCloser(bytes.NewBuffer(r.Body()))
		}
		return nil
	})
	return c
}

func (c *Client) Deref() *resty.Client { return c.cc }

func (c *Client) CallSetting(path string, cos ...CallOption) *CallSettings {
	cs := &CallSettings{
		contentType: "application/json",
		accept:      "application/json",
		Path:        path,
		header:      make(http.Header),
		noAuth:      false,
	}
	for _, co := range c.callOptions {
		co(cs)
	}
	for _, co := range cos {
		co(cs)
	}
	return cs
}

// Invoke the request
// NOTE: Do not use this function. use Execute instead.
func (c *Client) Invoke(ctx context.Context, method, path string, in, out any, settings *CallSettings) error {
	if c.validate != nil {
		if err := c.validate(in); err != nil {
			return err
		}
	}
	ctx = WithValueCallOption(ctx, *settings)
	r := c.cc.R().SetContext(ctx)
	if in != nil {
		reqBody, err := c.codec.Encode(settings.contentType, in)
		if err != nil {
			return err
		}
		r = r.SetBody(reqBody)
	}
	if !settings.noAuth {
		if c.tokenSource == nil {
			return errors.New("transport: token source should be not nil")
		}
		tk, err := c.tokenSource.Token()
		if err != nil {
			return err
		}
		r.SetHeader("Authorization", tk.Type()+" "+tk.AccessToken)
	}
	r.SetHeader("Content-Type", settings.contentType)
	r.SetHeader("Accept", settings.accept)
	for k, vs := range settings.header {
		for _, v := range vs {
			r.Header.Add(k, v)
		}
	}
	resp, err := r.Execute(method, c.cc.BaseURL+path)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return &ErrorReply{
			Code:   resp.StatusCode(),
			Body:   resp.Body(),
			Header: resp.Header(),
		}
	}
	defer resp.RawResponse.Body.Close()
	return c.codec.InboundForResponse(resp.RawResponse).NewDecoder(resp.RawResponse.Body).Decode(out)
}

func hasRequestBody(method string) bool {
	_, ok := noRequestBodyMethods[method]
	return !ok
}

func (c *Client) Execute(ctx context.Context, method, path string, req, resp any, opts ...CallOption) error {
	var r any

	settings := c.CallSetting(path, opts...)
	hasBody := hasRequestBody(method)
	if hasBody {
		r = req
	}
	url := c.EncodeURL(settings.Path, req, !hasBody)
	return c.Invoke(ctx, method, url, r, &resp, settings)
}

// Get method does GET HTTP request. It's defined in section 4.3.1 of RFC7231.
func (c *Client) Get(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodGet, path, req, resp, opts...)
}

// Head method does HEAD HTTP request. It's defined in section 4.3.2 of RFC7231.
func (c *Client) Head(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodHead, path, req, resp, opts...)
}

// Post method does POST HTTP request. It's defined in section 4.3.3 of RFC7231.
func (c *Client) Post(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodPost, path, req, resp, opts...)
}

// Put method does PUT HTTP request. It's defined in section 4.3.4 of RFC7231.
func (c *Client) Put(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodPut, path, req, resp, opts...)
}

// Delete method does DELETE HTTP request. It's defined in section 4.3.5 of RFC7231.
func (c *Client) Delete(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodDelete, path, req, resp, opts...)
}

// Options method does OPTIONS HTTP request. It's defined in section 4.3.7 of RFC7231.
func (c *Client) Options(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodOptions, path, req, resp, opts...)
}

// Patch method does PATCH HTTP request. It's defined in section 2 of RFC5789.
func (c *Client) Patch(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodPatch, path, req, resp, opts...)
}

// EncodeURL encode msg to url path.
// pathTemplate is a template of url path like http://helloworld.dev/{name}/sub/{sub.name}.
func (c *Client) EncodeURL(pathTemplate string, msg any, needQuery bool) string {
	return c.codec.EncodeURL(pathTemplate, msg, needQuery)
}

// EncodeQuery encode v into “URL encoded” form
// ("bar=baz&foo=quux") sorted by key.
func (c *Client) EncodeQuery(v any) (string, error) {
	vv, err := c.codec.EncodeQuery(v)
	if err != nil {
		return "", err
	}
	return vv.Encode(), nil
}
