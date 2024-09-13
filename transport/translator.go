package transport

import "context"

type TranslatorError interface {
	TranslateError(ctx context.Context, err error) (statusCode int, v any)
}

type TranslatorBody interface {
	TranslateBody(ctx context.Context, v any) any
}
