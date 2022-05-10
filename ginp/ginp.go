package ginp

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/things-go/ginp/errors"
)

func Response(c *gin.Context, data ...interface{}) {
	var obj interface{}

	if len(data) > 0 {
		obj = data[0]
	} else {
		obj = struct{}{}
	}
	c.JSON(http.StatusOK, obj)
}

func Abort(c *gin.Context, err error) {
	e := errors.FromError(err)
	if !gin.IsDebugging() {
		e.Detail = ""
	}

	status := 599
	switch {
	case e.Code == -1:
		status = http.StatusInternalServerError
	case e.Code < 1000:
		status = int(e.Code)
	}
	c.AbortWithStatusJSON(status, e)
}

func ErrorEncoder(c *gin.Context, err error, isBadRequest bool) {
	if isBadRequest {
		err = errors.ErrBadRequest(err.Error())
	}
	Abort(c, err)
}

// Attachment application/octet-stream;charset=utf-8
func Attachment(c *gin.Context, filename string, data []byte) {
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/octet-stream;charset=utf-8", data)
}
