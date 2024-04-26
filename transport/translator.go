package transport

type ErrorTranslator interface {
	Translate(error) error
}

type TranslatorError interface {
	TranslateError(err error) (statusCode int, v any)
}

type TranslatorData interface {
	TranslateData(v any) any
}
