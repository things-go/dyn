package http

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

// DisableBindValidation disable bind validation
func DisableBindValidation() { disableBindValidation = true }

// Validator return default Validator
func Validator() *validator.Validate { return defaultValidator }

// Validate like StructCtx but control by disableBindValidation.
func Validate(ctx context.Context, v any) error {
	if disableBindValidation {
		return nil
	}
	return defaultValidator.StructCtx(ctx, v)
}

// StructCtx validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified
// and also allows passing of context.Context for contextual validation information.
func StructCtx(ctx context.Context, v any) error {
	return defaultValidator.StructCtx(ctx, v)
}

// Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified.
func Struct(v any) error {
	return defaultValidator.Struct(v)
}

// VarCtx validates a single variable using tag style validation and allows passing of contextual
// validation information via context.Context.
func VarCtx(ctx context.Context, v any, tag string) error {
	return defaultValidator.VarCtx(ctx, v, tag)
}

// StructCtx validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified
// and also allows passing of context.Context for contextual validation information.
func Var(v any, tag string) error {
	return defaultValidator.Var(v, tag)
}
