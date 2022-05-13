package jwtauth

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/things-go/dyn/core/metadata"
)

type Account struct {
	Subject  string            `json:"subject,omitempty"`
	Type     string            `json:"type,omitempty"`
	Scopes   []string          `json:"scopes,omitempty"`
	Metadata metadata.Metadata `json:"metadata,omitempty"`
}

func (sf *Auth) GenerateToken(id string, acc Account) (string, time.Time, error) {
	return sf.generateToken(id, acc, sf.timeout)
}

func (sf *Auth) GenerateRefreshToken(id string, acc Account) (string, time.Time, error) {
	return sf.generateToken(id, acc, sf.maxTimeout)
}

func (sf *Auth) generateToken(id string, acc Account, timeout time.Duration) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(timeout)
	token, err := sf.NewWithClaims(&Claims{
		Type:     acc.Type,
		Scopes:   acc.Scopes,
		Metadata: acc.Metadata,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    sf.issuer,
			Subject:   acc.Subject,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        id,
		},
	})
	return token, expiresAt, err
}

func (sf *Auth) ParseAccountFromRequest(r *http.Request) (*Account, error) {
	claims, err := sf.ParseFromRequest(r)
	if err != nil {
		return nil, err
	}
	cc := claims.(*Claims)
	return &Account{
		Subject:  cc.Subject,
		Type:     cc.Type,
		Scopes:   cc.Scopes,
		Metadata: cc.Metadata,
	}, nil
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

func FromType(c context.Context) string {
	if v, ok := FromAccount(c); ok {
		return v.Type
	}
	return ""
}

func FromMetadata(c context.Context) metadata.Metadata {
	if v, ok := FromAccount(c); ok {
		return v.Metadata
	}
	return metadata.Metadata{}
}

func Subject(c *gin.Context) string {
	if v, ok := FromAccount(c.Request.Context()); ok {
		return v.Subject
	}
	return ""
}
