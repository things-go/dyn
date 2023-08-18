package transport

type ErrorTranslator interface {
	Translate(err error) error
}
