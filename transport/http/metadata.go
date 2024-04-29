package http

import (
	"github.com/gin-gonic/gin"
)

const ExclusivelyMetadataKey = "_dyn/transport/http/metadata"

type Metadata struct {
	Service string
	Method  string
}

// MetadataInterceptor is used to store a `ExclusivelyMetadataKey`/Metadata pair exclusively for this `*gin.Context`.
func MetadataInterceptor(md Metadata) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ExclusivelyMetadataKey, md)
		c.Next()
	}
}

// GetMetadata returns the Metadata value stored in `*gin.Context`, if any.
func GetMetadata(c *gin.Context) (md Metadata, ok bool) {
	v, ok := c.Get(ExclusivelyMetadataKey)
	if ok {
		md, ok = v.(Metadata)
	}
	return md, ok
}

// MustGetMetadata returns the Metadata value stored in `*gin.Context`.
func MustGetMetadata(c *gin.Context) Metadata {
	p, ok := GetMetadata(c)
	if !ok {
		panic("http: must be set \"" + ExclusivelyMetadataKey + "\" into gin.Context but it is not!!!")
	}
	return p
}
