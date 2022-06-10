package middleware

import (
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
	fullMethod string
	route      string
	clientIp   string
	ginContext *gin.Context
}

// Kind returns the transport kind.
func (tr *Transport) Kind() transport.Kind { return transport.KindHTTP }

// FullPath Service full method or path
func (tr *Transport) FullPath() string { return tr.fullMethod }

// ClientIp Service full method or path
func (tr *Transport) ClientIp() string { return tr.clientIp }

// Route Service full route
func (tr *Transport) Route() string { return tr.route }

// GinContext Service gin context
func (tr *Transport) GinContext() *gin.Context { return tr.ginContext }

func TransportInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		tr := &Transport{
			fullMethod: c.Request.URL.Path,
			route:      c.FullPath(),
			clientIp:   c.ClientIP(),
			ginContext: c,
		}
		c.Request = c.Request.WithContext(transport.WithValueTransporter(c.Request.Context(), tr))
		c.Next()
	}
}
