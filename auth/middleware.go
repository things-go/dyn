package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Option is Middleware option.
type Option func(*options)

// options is a Middleware option
type options struct {
	skip                 func(c *gin.Context) bool
	unauthorizedFallback func(*gin.Context, error)
}

// WithSkip set skip func
func WithSkip(f func(c *gin.Context) bool) Option {
	return func(o *options) {
		if f != nil {
			o.skip = f
		}
	}
}

// WithUnauthorizedFallback sets the fallback handler when requests are unauthorized.
func WithUnauthorizedFallback(f func(c *gin.Context, err error)) Option {
	return func(o *options) {
		if f != nil {
			o.unauthorizedFallback = f
		}
	}
}

func (sf *Auth) Middleware(opts ...Option) gin.HandlerFunc {
	o := &options{
		unauthorizedFallback: func(c *gin.Context, err error) {
			c.String(http.StatusUnauthorized, err.Error())
		},
		skip: func(c *gin.Context) bool { return false },
	}
	for _, opt := range opts {
		opt(o)
	}
	return func(c *gin.Context) {
		if !o.skip(c) {
			acc, err := sf.ParseFromRequest(c.Request)
			if err != nil {
				o.unauthorizedFallback(c, err)
				c.Abort()
				return
			}
			c.Request = c.Request.WithContext(NewContext(c.Request.Context(), acc))
		}
		c.Next()
	}
}

func Subject(c *gin.Context) string {
	if v, ok := FromContext(c.Request.Context()); ok {
		return v.Subject
	}
	return ""
}
