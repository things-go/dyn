package binding

import (
	"context"

	"github.com/go-playground/validator/v10"
)

var disableBindValidation bool
var defaultValidator = func() *validator.Validate {
	v := validator.New()
	v.SetTagName("binding")
	return v
}()

func DisableBindValidation() { disableBindValidation = true }

func Validator() *validator.Validate { return defaultValidator }

func Validate(ctx context.Context, v any) error {
	if disableBindValidation {
		return nil
	}
	return defaultValidator.StructCtx(ctx, v)
}

func ValidateStruct(v any) error { return defaultValidator.Struct(v) }
