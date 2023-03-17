package ginp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/things-go/dyn/genproto/errors"
	transportHttp "github.com/things-go/dyn/transport/http"
)

var _ transportHttp.Convey = (*ConveyorGin)(nil)

type ConveyorGin struct {
	validation *validator.Validate
}

func NewConveyorGin() *ConveyorGin {
	return &ConveyorGin{
		validation: func() *validator.Validate {
			v := validator.New()
			v.SetTagName("binding")
			return v
		}(),
	}
}
func (*ConveyorGin) WithValueUri(req *http.Request, params gin.Params) *http.Request {
	return transportHttp.WithValueUri(req, params)
}
func (*ConveyorGin) Bind(cg *gin.Context, v any) error {
	return cg.ShouldBind(v)
}
func (*ConveyorGin) BindQuery(cg *gin.Context, v any) error {
	return cg.ShouldBindQuery(v)
}
func (*ConveyorGin) BindUri(cg *gin.Context, v any) error {
	return cg.ShouldBindUri(v)
}
func (*ConveyorGin) ErrorBadRequest(cg *gin.Context, err error) {
	Abort(cg, errors.ErrBadRequest(err.Error()))
}
func (*ConveyorGin) Error(cg *gin.Context, err error) {
	Abort(cg, err)
}
func (*ConveyorGin) Render(cg *gin.Context, v any) {
	Response(cg, v)
}
func (cg *ConveyorGin) Validator() *validator.Validate {
	return cg.validation
}
func (cg *ConveyorGin) Validate(ctx context.Context, v any) error {
	return cg.validation.StructCtx(ctx, v)
}
func (cg *ConveyorGin) StructCtx(ctx context.Context, v any) error {
	return cg.validation.StructCtx(ctx, v)
}
func (cg *ConveyorGin) Struct(v any) error {
	return cg.validation.Struct(v)
}
func (cg *ConveyorGin) VarCtx(ctx context.Context, v any, tag string) error {
	return cg.validation.VarCtx(ctx, v, tag)
}
func (cg *ConveyorGin) Var(v any, tag string) error {
	return cg.validation.Var(v, tag)
}
