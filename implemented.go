package ginp

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Implemented struct{}

func (*Implemented) Validate(ctx context.Context, v interface{}) error {
	return Validate(ctx, v)
}

func (*Implemented) ErrorEncoder(c *gin.Context, err error, isBadRequest bool) {
	ErrorEncoder(c, err, isBadRequest)
}
