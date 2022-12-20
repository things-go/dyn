package ginp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/things-go/dyn/errors"
	dhttp "github.com/things-go/dyn/transport/http"
)

func Response(c *gin.Context, data ...any) {
	var obj any

	if len(data) > 0 {
		obj = data[0]
	} else {
		obj = struct{}{}
	}
	c.JSON(http.StatusOK, obj)
}

func Abort(c *gin.Context, err error) {
	e := errors.FromError(err)

	status := 599
	switch {
	case e.Code == -1:
		status = http.StatusInternalServerError
	case e.Code < 1000:
		status = int(e.Code)
	}
	c.AbortWithStatusJSON(status, e)
}

type Implemented struct{}

func (*Implemented) Validate(ctx context.Context, v any) error {
	return dhttp.Validate(ctx, v)
}

func (*Implemented) ErrorEncoder(c *gin.Context, err error, isBadRequest bool) {
	if isBadRequest {
		err = errors.ErrBadRequest(err.Error())
	}
	Abort(c, err)
}
