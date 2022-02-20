package ginp

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/things-go/ginp/errors"
)

// ResponseBody 应答
type ResponseBody struct {
	Code    int32       `json:"code"`              // 业务代码
	Message string      `json:"message,omitempty"` // 消息
	Detail  string      `json:"detail,omitempty"`  // 主要用于开发调试, gin.IsDebugging() == false 时不显示
	Data    interface{} `json:"data"`              // 应用数据
}

func Response(c *gin.Context, data ...interface{}) {
	var obj interface{}

	if len(data) > 0 {
		obj = data[0]
	} else {
		obj = struct{}{}
	}
	c.JSON(http.StatusOK, obj)
}

func Abort(c *gin.Context, err error, data ...interface{}) {
	e := errors.FromError(err)
	r := ResponseBody{
		Code:    e.Code,
		Message: e.Message,
	}
	if e.Message != "" {
		r.Message = e.Message
	}

	if gin.IsDebugging() {
		r.Detail = e.Detail
	}
	if len(data) > 0 && data[0] != nil {
		r.Data = data[0]
	}

	status := 599
	switch {
	case e.Code == -1:
		status = http.StatusInternalServerError
	case e.Code < 1000:
		status = int(e.Code)
	}
	c.AbortWithStatusJSON(status, r)
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
