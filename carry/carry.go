package ginp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/things-go/dyn/transport"
	transportHttp "github.com/things-go/dyn/transport/http"
	"github.com/things-go/encoding"
)

var _ transportHttp.Carrier = (*Carry)(nil)

type Carry struct {
	encoding   *encoding.Encoding
	validation *validator.Validate
	// translate error
	translate             transport.ErrorTranslator
	renderErrorTranslator transport.RenderErrorTranslator
	renderTranslator      transport.RenderTranslator
}

type Option func(*Carry)

func WithEncoding(e *encoding.Encoding) Option {
	return func(cy *Carry) {
		cy.encoding = e
	}
}
func WithValidation(v *validator.Validate) Option {
	return func(cy *Carry) {
		cy.validation = v
	}
}
func WithTranslateError(t transport.ErrorTranslator) Option {
	return func(cy *Carry) {
		cy.translate = t
	}
}
func NewCarry(opts ...Option) *Carry {
	cy := &Carry{
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

func (cy *Carry) Bind(c *gin.Context, v any) error {
	return cy.encoding.Bind(c.Request, v)
}
func (cy *Carry) BindQuery(c *gin.Context, v any) error {
	return cy.encoding.BindQuery(c.Request, v)
}
func (cy *Carry) BindUri(c *gin.Context, v any) error {
	return cy.encoding.BindUri(transportHttp.UrlValues(c.Params), v)
}
func (cy *Carry) Error(c *gin.Context, err error) {
	var obj any
	var statusCode = http.StatusInternalServerError

	if cy.translate != nil {
		err = cy.translate.Translate(err)
	}
	if cy.renderErrorTranslator != nil {
		statusCode, obj = cy.renderErrorTranslator.TranslateError(err)
	} else {
		obj = err.Error()
	}
	c.Writer.WriteHeader(statusCode)
	if err := cy.encoding.Render(c.Writer, c.Request, obj); err != nil {
		c.String(http.StatusInternalServerError, "Render failed cause by %v", err)
	}
}
func (cy *Carry) Render(c *gin.Context, v any) {
	if cy.renderTranslator != nil {
		v = cy.renderTranslator.TranslateData(v)
	}
	c.Writer.WriteHeader(http.StatusOK)
	err := cy.encoding.Render(c.Writer, c.Request, v)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render failed cause by %v", err)
	}
}
func (cy *Carry) Validator() *validator.Validate {
	return cy.validation
}
func (cy *Carry) Validate(ctx context.Context, v any) error {
	return cy.validation.StructCtx(ctx, v)
}
func (cy *Carry) StructCtx(ctx context.Context, v any) error {
	return cy.validation.StructCtx(ctx, v)
}
func (cy *Carry) Struct(v any) error {
	return cy.validation.Struct(v)
}
func (cy *Carry) VarCtx(ctx context.Context, v any, tag string) error {
	return cy.validation.VarCtx(ctx, v, tag)
}
func (cy *Carry) Var(v any, tag string) error {
	return cy.validation.Var(v, tag)
}
