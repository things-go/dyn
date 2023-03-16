package ginp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/things-go/dyn/errors"
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
func (i *ConveyorGin) WithValueUri(req *http.Request, params gin.Params) *http.Request {
	return transportHttp.WithValueUri(req, params)
}
func (i *ConveyorGin) Bind(c *gin.Context, v any) error {
	return c.ShouldBind(v)
}
func (i *ConveyorGin) BindQuery(c *gin.Context, v any) error {
	return c.ShouldBindQuery(v)
}
func (i *ConveyorGin) BindUri(c *gin.Context, v any) error {
	return c.ShouldBindUri(v)
}
func (*ConveyorGin) ErrorBadRequest(c *gin.Context, err error) {
	Abort(c, errors.ErrBadRequest(err.Error()))
}
func (*ConveyorGin) Error(c *gin.Context, err error) {
	Abort(c, err)
}
func (i *ConveyorGin) Render(c *gin.Context, v any) {
	c.JSON(http.StatusOK, v)
}
func (i *ConveyorGin) Validator() *validator.Validate {
	return i.validation
}
func (i *ConveyorGin) Validate(ctx context.Context, v any) error {
	return i.validation.StructCtx(ctx, v)
}
func (i *ConveyorGin) StructCtx(ctx context.Context, v any) error {
	return i.validation.StructCtx(ctx, v)
}
func (i *ConveyorGin) Struct(v any) error {
	return i.validation.Struct(v)
}
func (i *ConveyorGin) VarCtx(ctx context.Context, v any, tag string) error {
	return i.validation.VarCtx(ctx, v, tag)
}
func (i *ConveyorGin) Var(v any, tag string) error {
	return i.validation.Var(v, tag)
}
