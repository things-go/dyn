package http

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type ctxCarrierKey struct{}

// Carrier is an HTTP Carrier.
type Carrier interface {
	// WithValueUri sets the URL params for the given request.
	//
	// Deprecated: Use BindURI not need this.
	WithValueUri(*http.Request, gin.Params) *http.Request
	// BindUri binds the passed struct pointer using the uri codec.Marshaler.
	// NOTE: before use this, you should set uri params in the request context with RequestWithUri.
	//
	// Deprecated: Use BindURI not need this.
	BindUri(*gin.Context, any) error

	// Bind checks the Method and Content-Type to select codec.Marshaler automatically,
	// Depending on the "Content-Type" header different bind are used.
	Bind(*gin.Context, any) error
	// BindQuery binds the passed struct pointer using the query codec.Marshaler.
	BindQuery(*gin.Context, any) error
	// BindUri binds the passed struct pointer using the uri codec.Marshaler.
	BindURI(*gin.Context, url.Values, any) error
	// Error encode error response.
	Error(*gin.Context, error)
	// Render encode response.
	Render(*gin.Context, any)
	// Validate the request.
	Validate(context.Context, any) error
}

// WithValueCarrier returns the value associated with ctxCarrierKey is
// Carrier.
func WithValueCarrier(ctx context.Context, c Carrier) context.Context {
	return context.WithValue(ctx, ctxCarrierKey{}, c)
}

// FromCarrier returns the Carrier value stored in ctx, if not exist cause panic.
func FromCarrier(ctx context.Context) Carrier {
	c, ok := ctx.Value(ctxCarrierKey{}).(Carrier)
	if !ok {
		panic("carrier: must be set Carrier into context but it is not!!!")
	}
	return c
}

// CarrierInterceptor carrier middleware.
func CarrierInterceptor(carrier Carrier) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(WithValueCarrier(c.Request.Context(), carrier))
		c.Next()
	}
}
