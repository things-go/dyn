package ginp

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/things-go/ginp/errors"
)

var disableBindValidation bool
var defaultValidator = func() *validator.Validate {
	v := validator.New()
	v.SetTagName("binding")
	return v
}()

func DisableBindValidation() { disableBindValidation = true }

func Validator() *validator.Validate { return defaultValidator }

func Validate(ctx context.Context, v interface{}) error {
	if disableBindValidation {
		return nil
	}
	return defaultValidator.StructCtx(ctx, v)
}

func ErrorEncoder(c *gin.Context, err error, isBadRequest bool) {
	if isBadRequest {
		err = errors.ErrBadRequest(err.Error())
	}
	Abort(c, err)
}
