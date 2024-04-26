package http

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func UrlValues(params gin.Params) url.Values {
	vars := make(url.Values, len(params))
	for _, p := range params {
		vars.Add(p.Key, p.Value)
	}
	return vars
}
