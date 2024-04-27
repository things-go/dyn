package ginp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/things-go/encoding"

	"github.com/things-go/dyn/transport"
	transportHttp "github.com/things-go/dyn/transport/http"
)

var _ transportHttp.Carrier = (*CarryGin)(nil)
var _ Applier = (*CarryGin)(nil)

type CarryGin struct {
	validation      *validator.Validate
	errorTranslate  transport.ErrorTranslator
	translatorError transport.TranslatorError
	translatorData  transport.TranslatorData
}

func NewCarryGin() *CarryGin {
	return &CarryGin{
		validation: func() *validator.Validate {
			v := validator.New()
			v.SetTagName("binding")
			return v
		}(),
	}
}

func (cy *CarryGin) setEncoding(e *encoding.Encoding) {}
func (cy *CarryGin) setValidation(v *validator.Validate) {
	cy.validation = v
}
func (cy *CarryGin) setErrorTranslator(e transport.ErrorTranslator) {
	cy.errorTranslate = e
}

func (cy *CarryGin) setTranslatorError(e transport.TranslatorError) {
	cy.translatorError = e
}

func (cy *CarryGin) setTranslatorData(e transport.TranslatorData) {
	cy.translatorData = e
}

func (*CarryGin) Bind(c *gin.Context, v any) error {
	return c.ShouldBind(v)
}
func (*CarryGin) BindQuery(c *gin.Context, v any) error {
	return c.ShouldBindQuery(v)
}
func (*CarryGin) BindUri(c *gin.Context, v any) error {
	return c.ShouldBindUri(v)
}
func (cy *CarryGin) Error(c *gin.Context, err error) {
	var obj any
	var statusCode = http.StatusInternalServerError

	if cy.errorTranslate != nil {
		err = cy.errorTranslate.Translate(err)
	}
	if cy.translatorError != nil {
		statusCode, obj = cy.translatorError.TranslateError(err)
	} else {
		obj = err.Error()
	}
	c.AbortWithStatusJSON(statusCode, obj)
}
func (cy *CarryGin) Render(c *gin.Context, v any) {
	if cy.translatorData != nil {
		v = cy.translatorData.TranslateData(v)
	}
	c.JSON(http.StatusOK, v)
}
func (cy *CarryGin) Validator() *validator.Validate {
	return cy.validation
}
func (cy *CarryGin) Validate(ctx context.Context, v any) error {
	return cy.validation.StructCtx(ctx, v)
}
func (cy *CarryGin) StructCtx(ctx context.Context, v any) error {
	return cy.validation.StructCtx(ctx, v)
}
func (cy *CarryGin) Struct(v any) error {
	return cy.validation.Struct(v)
}
func (cy *CarryGin) VarCtx(ctx context.Context, v any, tag string) error {
	return cy.validation.VarCtx(ctx, v, tag)
}
func (cy *CarryGin) Var(v any, tag string) error {
	return cy.validation.Var(v, tag)
}
