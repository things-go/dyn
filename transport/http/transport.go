package http

import (
	"net/http"
	"net/textproto"

	"github.com/gin-gonic/gin"

	"github.com/things-go/dyn/transport"
)

var _ Transporter = (*Transport)(nil)

// Transporter is http Transporter
type Transporter interface {
	transport.Transporter
	Method() string
	Route() string
	GinContext() *gin.Context
}

// Transport is an HTTP transport.
type Transport struct {
	fullPath       string
	method         string
	route          string
	clientIp       string
	requestHeader  header
	responseHeader header
	ginContext     *gin.Context
}

// Kind returns the transport kind.
func (tr *Transport) Kind() transport.Kind { return transport.HTTP }

// FullPath Service full method or path
func (tr *Transport) FullPath() string { return tr.fullPath }

// ClientIp Service full method or path
func (tr *Transport) ClientIp() string { return tr.clientIp }

// RequestHeader return transport request header
// http: http.Header
// grpc: metadata.MD
func (tr *Transport) RequestHeader() transport.Header { return tr.requestHeader }

// ResponseHeader return transport response header
// http: http.Header
// grpc: metadata.MD
func (tr *Transport) ResponseHeader() transport.Header { return tr.responseHeader }

// Method Service http method
func (tr *Transport) Method() string { return tr.method }

// Route Service full route
func (tr *Transport) Route() string { return tr.route }

// GinContext Service gin context
func (tr *Transport) GinContext() *gin.Context { return tr.ginContext }

type header http.Header

// Len returns the number of items in header.
func (h header) Len() int { return len(h) }

// Get returns the value associated with the passed key.
func (h header) Get(key string) string { return http.Header(h).Get(key) }

// Add adds the key, value pair to the header.
func (h header) Add(key, value string) { http.Header(h).Add(key, value) }

// Set stores the key-value pair.
func (h header) Set(key string, value string) { http.Header(h).Set(key, value) }

// Append adds the values to key k, not overwriting what was already stored at
// that key.
//
// k is converted to lowercase before storing in header.
func (h header) Append(key string, vals ...string) {
	if len(vals) == 0 {
		return
	}
	key = textproto.CanonicalMIMEHeaderKey(key)
	h[key] = append(h[key], vals...)
}

// Delete removes the values for a given key k which is converted to lowercase
// before removing it from header.
func (h header) Delete(key string) { textproto.MIMEHeader(h).Del(key) }

// Keys lists the keys stored in this carrier.
func (h header) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range http.Header(h) {
		keys = append(keys, k)
	}
	return keys
}

// Clone returns a copy of h or nil if h is nil.
func (h header) Clone() transport.Header { return transport.Header(header(http.Header(h).Clone())) }

// TransportInterceptor transport middleware
func TransportInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		tr := &Transport{
			c.Request.URL.Path,
			c.Request.Method,
			c.FullPath(),
			c.ClientIP(),
			header(c.Request.Header),
			header(c.Writer.Header()),
			c,
		}
		c.Request = c.Request.WithContext(transport.WithValueTransporter(c.Request.Context(), tr))
		c.Next()
	}
}
