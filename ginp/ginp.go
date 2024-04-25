package ginp

import (
	"context"
	stdErrors "errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/things-go/dyn/errors"
	"github.com/things-go/dyn/transport"
	transportHttp "github.com/things-go/dyn/transport/http"
	"github.com/things-go/encoding"
)

var _ transportHttp.Carrier = (*Carry)(nil)

type Carry struct {
	encoding   *encoding.Encoding
	validation *validator.Validate
	// translate error
	translate transport.ErrorTranslator
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

// Deprecated: Use BindURI not need this.
func (cy *Carry) WithValueUri(req *http.Request, params gin.Params) *http.Request {
	return transportHttp.WithValueUri(req, params)
}

// Deprecated: Use BindURI not need this.
func (cy *Carry) BindUri(c *gin.Context, v any) error {
	return cy.encoding.BindUri(c.Request, v)
}

func (cy *Carry) Bind(c *gin.Context, v any) error {
	return cy.encoding.Bind(c.Request, v)
}
func (cy *Carry) BindQuery(c *gin.Context, v any) error {
	return cy.encoding.BindQuery(c.Request, v)
}
func (cy *Carry) BindURI(c *gin.Context, raws url.Values, v any) error {
	return cy.encoding.BindURI(raws, v)
}

func (cy *Carry) Error(c *gin.Context, err error) {
	if cy.translate != nil {
		err = cy.translate.Translate(err)
	}
	if e := new(validator.ValidationErrors); stdErrors.As(err, e) {
		err = errors.ErrBadRequest(err.Error())
	}
	Abort(c, err)
}
func (cy *Carry) Render(c *gin.Context, v any) {
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
