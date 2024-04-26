package ginp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/things-go/dyn/transport"
	transportHttp "github.com/things-go/dyn/transport/http"
)

var _ transportHttp.Carrier = (*GinCarry)(nil)

type GinCarry struct {
	validation *validator.Validate
	// translate error
	translate             transport.ErrorTranslator
	renderErrorTranslator transport.RenderErrorTranslator
	renderTranslator      transport.RenderTranslator
}

func NewCarryForGin() *GinCarry {
	return &GinCarry{
		validation: func() *validator.Validate {
			v := validator.New()
			v.SetTagName("binding")
			return v
		}(),
	}
}

func (cy *GinCarry) SetTranslateError(e transport.ErrorTranslator) *GinCarry {
	cy.translate = e
	return cy
}

func (*GinCarry) Bind(c *gin.Context, v any) error {
	return c.ShouldBind(v)
}
func (*GinCarry) BindQuery(c *gin.Context, v any) error {
	return c.ShouldBindQuery(v)
}
func (*GinCarry) BindUri(c *gin.Context, v any) error {
	return c.ShouldBindUri(v)
}
func (cy *GinCarry) Error(c *gin.Context, err error) {
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
	c.AbortWithStatusJSON(statusCode, obj)
}
func (cy *GinCarry) Render(c *gin.Context, v any) {
	if cy.renderTranslator != nil {
		v = cy.renderTranslator.TranslateData(v)
	}
	c.JSON(http.StatusOK, v)
}
func (cy *GinCarry) Validator() *validator.Validate {
	return cy.validation
}
func (cy *GinCarry) Validate(ctx context.Context, v any) error {
	return cy.validation.StructCtx(ctx, v)
}
func (cy *GinCarry) StructCtx(ctx context.Context, v any) error {
	return cy.validation.StructCtx(ctx, v)
}
func (cy *GinCarry) Struct(v any) error {
	return cy.validation.Struct(v)
}
func (cy *GinCarry) VarCtx(ctx context.Context, v any, tag string) error {
	return cy.validation.VarCtx(ctx, v, tag)
}
func (cy *GinCarry) Var(v any, tag string) error {
	return cy.validation.Var(v, tag)
}
