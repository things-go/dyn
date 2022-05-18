package gctx

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ctxGinKey struct{}

func WithValueGin(ctx context.Context, c *gin.Context) context.Context {
	return context.WithValue(ctx, ctxGinKey{}, c)
}

func FromGin(ctx context.Context) (*gin.Context, bool) {
	c, ok := ctx.Value(ctxGinKey{}).(*gin.Context)
	return c, ok
}

func MustFromGin(ctx context.Context) *gin.Context {
	c, ok := ctx.Value(ctxGinKey{}).(*gin.Context)
	if !ok {
		panic("must be set gin into context but it is not!!!")
	}
	return c
}

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(WithValueGin(c.Request.Context(), c))
		c.Next()
	}
}
