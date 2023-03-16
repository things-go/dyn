package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ctxConveyKey struct{}

// Convey is an HTTP Convey.
type Convey interface {
	// WithValueUri sets the URL params for the given request.
	WithValueUri(*http.Request, gin.Params) *http.Request
	// Bind checks the Method and Content-Type to select codec.Marshaler automatically,
	// Depending on the "Content-Type" header different bind are used.
	Bind(*gin.Context, any) error
	// BindQuery binds the passed struct pointer using the query codec.Marshaler.
	BindQuery(*gin.Context, any) error
	// BindUri binds the passed struct pointer using the uri codec.Marshaler.
	// NOTE: before use this, you should set uri params in the request context with RequestWithUri.
	BindUri(*gin.Context, any) error
	// ErrorBadRequest encode error response.
	ErrorBadRequest(*gin.Context, error)
	// Error encode error response.
	Error(c *gin.Context, err error)
	// Render encode response.
	Render(*gin.Context, any)

	// validator shortcut

	// Validator return default Validator
	Validator() *validator.Validate
	// Validate the request.
	Validate(context.Context, any) error
	// StructCtx validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified
	// and also allows passing of context.Context for contextual validation information.
	StructCtx(context.Context, any) error
	// Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified.
	Struct(any) error
	// VarCtx validates a single variable using tag style validation and allows passing of contextual
	// validation information via context.Context.
	VarCtx(context.Context, any, string) error
	// Var validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified
	// and also allows passing of context.Context for contextual validation information.
	Var(any, string) error
}

func WithValueConvey(ctx context.Context, c Convey) context.Context {
	return context.WithValue(ctx, ctxConveyKey{}, c)
}

// FromTransporter returns the Transporter value stored in ctx, if any.
func FromConvey(ctx context.Context) Convey {
	c, ok := ctx.Value(ctxConveyKey{}).(Convey)
	if !ok {
		panic("convey: must be set Convey into context but it is not!!!")
	}
	return c
}

// ConveyInterceptor convey middleware
func ConveyInterceptor(convey Convey) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(WithValueConvey(c.Request.Context(), convey))
		c.Next()
	}
}
