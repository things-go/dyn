package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/things-go/dyn/transport"
)

// Transporter is http Transporter
type Transporter interface {
	transport.Transporter
	Route() string
	GinContext() *gin.Context
}

// Transport is an HTTP transport.
type Transport struct {
	fullMethod     string
	route          string
	clientIp       string
	requestHeader  httpHeader
	responseHeader httpHeader
	ginContext     *gin.Context
}

// Kind returns the transport kind.
func (tr *Transport) Kind() transport.Kind { return transport.HTTP }

// FullPath Service full method or path
func (tr *Transport) FullPath() string { return tr.fullMethod }

// ClientIp Service full method or path
func (tr *Transport) ClientIp() string { return tr.clientIp }

// RequestHeader return transport request header
// http: http.Header
// grpc: metadata.MD
func (tr *Transport) RequestHeader() transport.Header { return tr.requestHeader }

// ReplyHeader return transport response header
// http: http.Header
// grpc: metadata.MD
func (tr *Transport) ReplyHeader() transport.Header { return tr.responseHeader }

// Route Service full route
func (tr *Transport) Route() string { return tr.route }

// GinContext Service gin context
func (tr *Transport) GinContext() *gin.Context { return tr.ginContext }

func TransportInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		tr := &Transport{
			fullMethod:     c.Request.URL.Path,
			route:          c.FullPath(),
			clientIp:       c.ClientIP(),
			requestHeader:  httpHeader(c.Request.Header),
			responseHeader: httpHeader(c.Writer.Header()),
			ginContext:     c,
		}
		c.Request = c.Request.WithContext(transport.WithValueTransporter(c.Request.Context(), tr))
		c.Next()
	}
}

type httpHeader http.Header

// Get returns the value associated with the passed key.
func (h httpHeader) Get(key string) string { return http.Header(h).Get(key) }

// Add adds the key, value pair to the header.
func (h httpHeader) Add(key, value string) { http.Header(h).Add(key, value) }

// Set stores the key-value pair.
func (h httpHeader) Set(key string, value string) { http.Header(h).Set(key, value) }

// Keys lists the keys stored in this carrier.
func (h httpHeader) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range http.Header(h) {
		keys = append(keys, k)
	}
	return keys
}

// Clone returns a copy of h or nil if h is nil.
func (h httpHeader) Clone() transport.Header {
	return transport.Header(httpHeader(http.Header(h).Clone()))
}
