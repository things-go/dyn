package ginp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/things-go/dyn/errors"
	transportHttp "github.com/things-go/dyn/transport/http"
	"github.com/things-go/encoding"
)

var _ transportHttp.Convey = (*Conveyor)(nil)

type Conveyor struct {
	encoding   *encoding.Encoding
	validation *validator.Validate
}

type Option func(*Conveyor)

func WithEncoding(e *encoding.Encoding) Option {
	return func(i *Conveyor) {
		i.encoding = e
	}
}

func WithValidation(v *validator.Validate) Option {
	return func(i *Conveyor) {
		i.validation = v
	}
}

func NewImplemented(opts ...Option) *Conveyor {
	i := &Conveyor{
		encoding: encoding.New(),
		validation: func() *validator.Validate {
			v := validator.New()
			v.SetTagName("binding")
			return v
		}(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}
func (i *Conveyor) WithValueUri(req *http.Request, params gin.Params) *http.Request {
	return transportHttp.WithValueUri(req, params)
}
func (i *Conveyor) Bind(c *gin.Context, v any) error {
	return i.encoding.Bind(c.Request, v)
}
func (i *Conveyor) BindQuery(c *gin.Context, v any) error {
	return i.encoding.BindQuery(c.Request, v)
}
func (i *Conveyor) BindUri(c *gin.Context, v any) error {
	return i.encoding.BindUri(c.Request, v)
}
func (*Conveyor) ErrorBadRequest(c *gin.Context, err error) {
	Abort(c, errors.ErrBadRequest(err.Error()))
}
func (*Conveyor) Error(c *gin.Context, err error) {
	Abort(c, err)
}
func (i *Conveyor) Render(c *gin.Context, v any) {
	c.Writer.WriteHeader(http.StatusOK)
	err := i.encoding.Render(c.Writer, c.Request, v)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render failed cause by %v", err)
	}
}
func (i *Conveyor) Validator() *validator.Validate {
	return i.validation
}
func (i *Conveyor) Validate(ctx context.Context, v any) error {
	return i.validation.StructCtx(ctx, v)
}
func (i *Conveyor) StructCtx(ctx context.Context, v any) error {
	return i.validation.StructCtx(ctx, v)
}
func (i *Conveyor) Struct(v any) error {
	return i.validation.Struct(v)
}
func (i *Conveyor) VarCtx(ctx context.Context, v any, tag string) error {
	return i.validation.VarCtx(ctx, v, tag)
}
func (i *Conveyor) Var(v any, tag string) error {
	return i.validation.Var(v, tag)
}
