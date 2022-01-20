package ginp

import (
	"context"

	"github.com/gin-gonic/gin"
)

type errorCtxKey struct{}
type titleCKey struct{}

// ContextError context error
func ContextError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	ctx := c.Request.Context()
	if e, ok := ctx.Value(errorCtxKey{}).(error); ok && e != nil {
		err = e
	}
	ctx = context.WithValue(ctx, errorCtxKey{}, err)
	c.Request = c.Request.WithContext(ctx)
}

// FromError error from context
func FromError(ctx context.Context) error {
	err, ok := ctx.Value(errorCtxKey{}).(error)
	if !ok {
		return nil
	}
	return err
}

// Error error
func Error(c *gin.Context) error {
	return FromError(c.Request.Context())
}

// ContextTitle context title
func ContextTitle(c *gin.Context, title string) {
	ctx := context.WithValue(c.Request.Context(), titleCKey{}, title)
	c.Request = c.Request.WithContext(ctx)
}

// FromTitle title from context
func FromTitle(ctx context.Context) string {
	title, ok := ctx.Value(titleCKey{}).(string)
	if !ok {
		return ""
	}
	return title
}

// Title from gin.Context
func Title(c *gin.Context) string {
	return FromTitle(c.Request.Context())
}
