package transport

type TranslatorError interface {
	TranslateError(err error) (statusCode int, v any)
}

type TranslatorBody interface {
	TranslateBody(v any) any
}
