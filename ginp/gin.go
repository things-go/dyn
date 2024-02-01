package ginp

import (
	"context"
	stdErrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/things-go/dyn/errors"
	"github.com/things-go/dyn/transport"
	transportHttp "github.com/things-go/dyn/transport/http"
)

var _ transportHttp.Carrier = (*GinCarry)(nil)

type GinCarry struct {
	validation *validator.Validate
	// translate error
	translate transport.ErrorTranslator
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
func (*GinCarry) WithValueUri(req *http.Request, params gin.Params) *http.Request {
	return transportHttp.WithValueUri(req, params)
}
func (*GinCarry) Bind(cg *gin.Context, v any) error {
	return cg.ShouldBind(v)
}
func (*GinCarry) BindQuery(cg *gin.Context, v any) error {
	return cg.ShouldBindQuery(v)
}
func (*GinCarry) BindUri(cg *gin.Context, v any) error {
	return cg.ShouldBindUri(v)
}
func (cy *GinCarry) Error(cg *gin.Context, err error) {
	if cy.translate != nil {
		err = cy.translate.Translate(err)
	}
	if e := new(validator.ValidationErrors); stdErrors.As(err, e) {
		err = errors.ErrBadRequest(err.Error())
	}
	Abort(cg, err)
}
func (*GinCarry) Render(cg *gin.Context, v any) {
	Response(cg, v)
}
func (cg *GinCarry) Validator() *validator.Validate {
	return cg.validation
}
func (cg *GinCarry) Validate(ctx context.Context, v any) error {
	return cg.validation.StructCtx(ctx, v)
}
func (cg *GinCarry) StructCtx(ctx context.Context, v any) error {
	return cg.validation.StructCtx(ctx, v)
}
func (cg *GinCarry) Struct(v any) error {
	return cg.validation.Struct(v)
}
func (cg *GinCarry) VarCtx(ctx context.Context, v any, tag string) error {
	return cg.validation.VarCtx(ctx, v, tag)
}
func (cg *GinCarry) Var(v any, tag string) error {
	return cg.validation.Var(v, tag)
}
