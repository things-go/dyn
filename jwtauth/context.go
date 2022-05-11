package jwtauth

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/things-go/dyn/core/metadata"
)

type ctxAuthKey struct{}

// NewContext put auth info into context
func NewContext(ctx context.Context, claims jwt.Claims) context.Context {
	return context.WithValue(ctx, ctxAuthKey{}, claims)
}

// FromContext extract auth info from context
func FromContext(ctx context.Context) (claims jwt.Claims, ok bool) {
	claims, ok = ctx.Value(ctxAuthKey{}).(jwt.Claims)
	return
}

func FromId(ctx context.Context) string {
	v, ok := FromContext(ctx)
	if !ok {
		return ""
	}
	vv, ok := v.(*Claims)
	if !ok {
		return ""
	}
	return vv.ID
}

func FromAccount(c context.Context) (Account, bool) {
	v, ok := FromContext(c)
	if !ok {
		return Account{}, false
	}
	vv, ok := v.(*Claims)
	if !ok {
		return Account{}, false
	}
	return Account{
		Subject:  vv.Subject,
		Type:     vv.Type,
		Scopes:   vv.Scopes,
		Metadata: vv.Metadata,
	}, true
}

func FromSubject(c context.Context) int64 {
	if v, ok := FromAccount(c); ok {
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

func FromMetadata(c context.Context) metadata.Metadata {
	if v, ok := FromAccount(c); ok {
		return v.Metadata
	}
	return metadata.Metadata{}
}

func FromType(c context.Context) string {
	if v, ok := FromAccount(c); ok {
		return v.Type
	}
	return ""
}

func Subject(c *gin.Context) string {
	if v, ok := FromAccount(c); ok {
		return v.Subject
	}
	return ""
}
