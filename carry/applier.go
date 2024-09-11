package carry

import (
	"github.com/go-playground/validator/v10"
	"github.com/things-go/dyn/transport"
	"github.com/things-go/encoding"
)

type Applier interface {
	setEncoding(*encoding.Encoding)
	setValidation(*validator.Validate)
	setTranslatorError(transport.TranslatorError)
	setTranslatorBody(transport.TranslatorBody)
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

func WithTranslatorError(t transport.TranslatorError) Option {
	return func(cy Applier) {
		cy.setTranslatorError(t)
	}
}
func WithTranslatorBody(t transport.TranslatorBody) Option {
	return func(cy Applier) {
		cy.setTranslatorBody(t)
	}
}
