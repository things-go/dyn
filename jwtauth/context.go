package jwtauth

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
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
