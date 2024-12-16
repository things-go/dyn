package transport

import "context"

type TransformError interface {
	TransformError(ctx context.Context, err error) (statusCode int, v any)
}

type TransformBody interface {
	TransformBody(ctx context.Context, v any) any
}
