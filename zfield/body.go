package zfield

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ctxBodyKey struct{}

type bodyValue struct {
	enabled bool
	limit   int
	req     string
	resp    *strings.Builder
}

type BodyConfig struct {
	Enabled bool // enabled
	Limit   int  // <=0: mean not limit
}

func PrepareBody(cfg BodyConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		bv := &bodyValue{
			limit:   cfg.Limit,
			enabled: cfg.Enabled,
			req:     "",
			resp:    &strings.Builder{},
		}

		if bv.enabled {
			c.Writer = &bodyWriter{ResponseWriter: c.Writer, dupBody: bv.resp}
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				c.Abort()
				return
			}
			c.Request.Body.Close()
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			if bv.limit > 0 && len(body) >= bv.limit {
				bv.req = "ignore larger req body"
			} else {
				bv.req = string(body)
			}
		}
		ctx := context.WithValue(c.Request.Context(), ctxBodyKey{}, bv)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func RequestBody(c *gin.Context) zap.Field {
	v, ok := c.Request.Context().Value(ctxBodyKey{}).(*bodyValue)
	if !ok || !v.enabled {
		return zap.Skip()
	}
	return zap.String("requestBody", v.req)
}

func ResponseBody(c *gin.Context) zap.Field {
	v, ok := c.Request.Context().Value(ctxBodyKey{}).(*bodyValue)
	if !ok || !v.enabled {
		return zap.Skip()
	}
	if v.limit >= 0 && v.resp.Len() >= v.limit {
		return zap.String("responseBody", "ignore larger response body")
	}
	return zap.String("responseBody", v.resp.String())
}

func FromRequestBody(ctx context.Context) string {
	bv, ok := ctx.Value(ctxBodyKey{}).(*bodyValue)
	if !ok {
		return ""
	}
	return bv.req
}

func FromResponseBody(ctx context.Context) string {
	v, ok := ctx.Value(ctxBodyKey{}).(*bodyValue)
	if !ok {
		return ""
	}
	return v.resp.String()
}

type bodyWriter struct {
	gin.ResponseWriter
	dupBody *strings.Builder
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.dupBody.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyWriter) WriteString(s string) (int, error) {
	w.dupBody.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
