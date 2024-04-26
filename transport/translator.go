package transport

type ErrorTranslator interface {
	Translate(err error) error
}

type RenderTranslator interface {
	TranslateData(v any) any
}
