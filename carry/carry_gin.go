package carry

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
	validation     *validator.Validate
	transformError transport.TransformError
	transformBody  transport.TransformBody
}

func NewCarryGin(opts ...Option) *CarryGin {
	cy := &CarryGin{
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

func (cy *CarryGin) setEncoding(e *encoding.Encoding) {}
func (cy *CarryGin) setValidation(v *validator.Validate) {
	cy.validation = v
}

func (cy *CarryGin) setTransformError(e transport.TransformError) {
	cy.transformError = e
}

func (cy *CarryGin) setTransformBody(e transport.TransformBody) {
	cy.transformBody = e
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
func (cy *CarryGin) ShouldBind(c *gin.Context, v any) error {
	if err := cy.Bind(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *CarryGin) ShouldBindQuery(c *gin.Context, v any) error {
	if err := cy.BindQuery(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *CarryGin) ShouldBindUri(c *gin.Context, v any) error {
	if err := cy.BindUri(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *CarryGin) ShouldBindBodyUri(c *gin.Context, v any) error {
	if err := cy.Bind(c, v); err != nil {
		return err
	}
	if err := cy.BindUri(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *CarryGin) ShouldBindQueryUri(c *gin.Context, v any) error {
	if err := cy.BindQuery(c, v); err != nil {
		return err
	}
	if err := cy.BindUri(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *CarryGin) ShouldBindQueryBody(c *gin.Context, v any) error {
	if err := cy.BindQuery(c, v); err != nil {
		return err
	}
	if err := cy.Bind(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *CarryGin) ShouldBindQueryBodyUri(c *gin.Context, v any) error {
	if err := cy.BindQuery(c, v); err != nil {
		return err
	}
	if err := cy.Bind(c, v); err != nil {
		return err
	}
	if err := cy.BindUri(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}

func (cy *CarryGin) Error(c *gin.Context, err error) {
	var obj any
	var statusCode = http.StatusInternalServerError

	if cy.transformError != nil {
		statusCode, obj = cy.transformError.TransformError(c.Request.Context(), err)
	} else {
		obj = err.Error()
	}
	c.AbortWithStatusJSON(statusCode, obj)
}
func (cy *CarryGin) Render(c *gin.Context, v any) {
	if cy.transformBody != nil {
		v = cy.transformBody.TransformBody(c.Request.Context(), v)
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
