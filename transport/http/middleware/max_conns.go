package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/things-go/core/limit"
	"github.com/things-go/dyn/log"
)

// MaxConns returns a middleware that limit the concurrent connections.
func MaxConns(n int) gin.HandlerFunc {
	if n <= 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	latch := limit.NewLimit(n)
	return func(c *gin.Context) {
		if latch.TryBorrow() {
			defer func() {
				if err := latch.Return(); err != nil {
					log.WithContext(c.Request.Context()).Error(err)
				}
			}()
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusServiceUnavailable)
		}
	}
}
