package ginp

import (
	"github.com/go-playground/validator/v10"
	"github.com/things-go/dyn/transport"
	"github.com/things-go/encoding"
)

type Applier interface {
	setEncoding(*encoding.Encoding)
	setValidation(*validator.Validate)
	setErrorTranslator(transport.ErrorTranslator)
	setTranslatorError(transport.TranslatorError)
	setTranslatorData(transport.TranslatorData)
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
func WithErrorTranslator(t transport.ErrorTranslator) Option {
	return func(cy Applier) {
		cy.setErrorTranslator(t)
	}
}
func WithTranslatorError(t transport.TranslatorError) Option {
	return func(cy Applier) {
		cy.setTranslatorError(t)
	}
}
func WithTranslatorData(t transport.TranslatorData) Option {
	return func(cy Applier) {
		cy.setTranslatorData(t)
	}
}
