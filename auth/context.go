package auth

import (
	"context"
	"strconv"

	"github.com/things-go/clip/metadata"
)

const TokenUniqueId = "dyn:auth:uniqueId" // nolint: revive

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

func FromSubject(ctx context.Context) int64 {
	if v, ok := FromContext(ctx); ok {
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

func FromId(ctx context.Context) string {
	md := FromMetadata(ctx)
	return md.Get(TokenUniqueId)
}

func FromType(ctx context.Context) string {
	if v, ok := FromContext(ctx); ok {
		return v.Type
	}
	return ""
}

func FromScopes(ctx context.Context) []string {
	if v, ok := FromContext(ctx); ok {
		return v.Scopes
	}
	return nil
}

func FromMetadata(ctx context.Context) metadata.Metadata {
	if v, ok := FromContext(ctx); ok && v.Metadata != nil {
		return v.Metadata
	}
	return metadata.Metadata{}
}
