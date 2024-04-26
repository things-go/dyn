package transport

type ErrorTranslator interface {
	Translate(error) error
}

type RenderErrorTranslator interface {
	TranslateError(err error) (statusCode int, v any)
}

type RenderTranslator interface {
	TranslateData(v any) any
}
