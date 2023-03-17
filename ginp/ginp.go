package ginp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/things-go/dyn/genproto/errors"
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
	return func(cy *Conveyor) {
		cy.encoding = e
	}
}

func WithValidation(v *validator.Validate) Option {
	return func(cy *Conveyor) {
		cy.validation = v
	}
}

func NewConveyor(opts ...Option) *Conveyor {
	cy := &Conveyor{
		encoding: encoding.New(),
		validation: func() *validator.Validate {
			v := validator.New()
			v.SetTagName("binding")
			return v
		}(),
	}
	for _, opt := range opts {
		opt(cy)
	}
	return cy
}
func (cy *Conveyor) WithValueUri(req *http.Request, params gin.Params) *http.Request {
	return transportHttp.WithValueUri(req, params)
}
func (cy *Conveyor) Bind(c *gin.Context, v any) error {
	return cy.encoding.Bind(c.Request, v)
}
func (cy *Conveyor) BindQuery(c *gin.Context, v any) error {
	return cy.encoding.BindQuery(c.Request, v)
}
func (cy *Conveyor) BindUri(c *gin.Context, v any) error {
	return cy.encoding.BindUri(c.Request, v)
}
func (*Conveyor) ErrorBadRequest(c *gin.Context, err error) {
	Abort(c, errors.ErrBadRequest(err.Error()))
}
func (*Conveyor) Error(c *gin.Context, err error) {
	Abort(c, err)
}
func (cy *Conveyor) Render(c *gin.Context, v any) {
	c.Writer.WriteHeader(http.StatusOK)
	err := cy.encoding.Render(c.Writer, c.Request, v)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render failed cause by %v", err)
	}
}
func (cy *Conveyor) Validator() *validator.Validate {
	return cy.validation
}
func (cy *Conveyor) Validate(ctx context.Context, v any) error {
	return cy.validation.StructCtx(ctx, v)
}
func (cy *Conveyor) StructCtx(ctx context.Context, v any) error {
	return cy.validation.StructCtx(ctx, v)
}
func (cy *Conveyor) Struct(v any) error {
	return cy.validation.Struct(v)
}
func (cy *Conveyor) VarCtx(ctx context.Context, v any, tag string) error {
	return cy.validation.VarCtx(ctx, v, tag)
}
func (cy *Conveyor) Var(v any, tag string) error {
	return cy.validation.Var(v, tag)
}
