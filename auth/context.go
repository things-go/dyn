package auth

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/things-go/dyn/core/metadata"
)

type ctxAuthKey struct{}

// NewContext put auth info into context
func NewContext(ctx context.Context, acc *Account) context.Context {
	return context.WithValue(ctx, ctxAuthKey{}, acc)
}

// FromContext extract auth info from context
func FromContext(ctx context.Context) (acc *Account, ok bool) {
	acc, ok = ctx.Value(ctxAuthKey{}).(*Account)
	return
}

func FromSubject(c context.Context) int64 {
	if v, ok := FromContext(c); ok {
		sub, err := strconv.ParseInt(v.Subject, 10, 64)
		if err == nil {
			return sub
		}
	}
	return 0
}

func MustFromSubject(ctx context.Context) int64 {
	if sub := FromSubject(ctx); sub != 0 {
		return sub
	}
	panic("auth: account info must in context, user must auth")
}

func FromType(c context.Context) string {
	if v, ok := FromContext(c); ok {
		return v.Type
	}
	return ""
}

func FromMetadata(c context.Context) metadata.Metadata {
	if v, ok := FromContext(c); ok {
		return v.Metadata
	}
	return metadata.Metadata{}
}

func Subject(c *gin.Context) string {
	if v, ok := FromContext(c.Request.Context()); ok {
		return v.Subject
	}
	return ""
}
