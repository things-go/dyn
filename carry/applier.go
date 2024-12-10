package carry

import (
	"github.com/go-playground/validator/v10"
	"github.com/things-go/dyn/transport"
	"github.com/things-go/encoding"
)

type Applier interface {
	setEncoding(*encoding.Encoding)
	setValidation(*validator.Validate)
	setTransformError(transport.TransformError)
	setTransformBody(transport.TransformBody)
}

type Option func(Applier)

func WithEncoding(e *encoding.Encoding) Option {
	return func(cy Applier) {
		cy.setEncoding(e)
	}
}
func WithValidation(v *validator.Validate) Option {
	return func(cy Applier) {
		cy.setValidation(v)
	}
}

func WithTransformError(t transport.TransformError) Option {
	return func(cy Applier) {
		cy.setTransformError(t)
	}
}
func WithTransformBody(t transport.TransformBody) Option {
	return func(cy Applier) {
		cy.setTransformBody(t)
	}
}
