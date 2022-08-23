package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/things-go/clip/limit"
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
					log.Println(err)
				}
			}()
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusServiceUnavailable)
		}
	}
}
