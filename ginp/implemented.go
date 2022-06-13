package ginp

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/things-go/dyn/binding"
	"github.com/things-go/dyn/errors"
)

type Implemented struct{}

func (*Implemented) Validate(ctx context.Context, v interface{}) error {
	return binding.Validate(ctx, v)
}

func (*Implemented) ErrorEncoder(c *gin.Context, err error, isBadRequest bool) {
	if isBadRequest {
		err = errors.ErrBadRequest(err.Error())
	}
	Abort(c, err)
}
