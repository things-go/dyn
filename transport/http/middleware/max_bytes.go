package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MaxBytes returns a middleware that limit reading of http request body.
func MaxBytes(n int64) gin.HandlerFunc {
	if n <= 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return func(c *gin.Context) {
		if c.Request.ContentLength > n {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		} else {
			c.Next()
		}
	}
}
