package auth

import (
	"net/http"
	"time"

	"github.com/things-go/clip/lookup"
	"github.com/things-go/clip/metadata"
)

type Provider interface {
	GenerateToken(id string, acc *Account, timeout time.Duration) (string, time.Time, error)
	GenerateRefreshToken(id string, acc *Account, timeout time.Duration) (string, time.Time, error)
	ParseToken(token string) (*Account, error)
}

type Account struct {
	// Subject of the account
	Subject string `json:"subject,omitempty"`
	// Type of the account, client, service, user
	Type string `json:"type,omitempty"`
	// Issuer of the account
	Issuer string `json:"issuer,omitempty"`
	// Scopes the account has access to
	Scopes []string `json:"scopes,omitempty"`
	// Metadata Any other associated metadata
	Metadata metadata.Metadata `json:"metadata,omitempty"`
}

// Config Auth config
type Config struct {
	// Timeout token valid time
	// if timeout <= refreshTimeout, refreshTimeout = timeout + 30 * time.Minute
	Timeout time.Duration
	// RefreshTimeout refresh token valid time.
	RefreshTimeout time.Duration
	// Issuer of the account
	Issuer string
	// Lookup used to extract token from the http request
	Lookup string
}

// Auth provides a Json-Web-Token authentication implementation.
type Auth struct {
	provider       Provider
	timeout        time.Duration
	refreshTimeout time.Duration
	// issuer of the account
	issuer string
	lookup *lookup.Lookup
}

// New auth with Config
func New(p Provider, c Config) *Auth {
	mw := &Auth{
		provider:       p,
		timeout:        c.Timeout,
		refreshTimeout: c.RefreshTimeout,
		issuer:         c.Issuer,
		lookup:         lookup.NewLookup(c.Lookup),
	}
	if mw.timeout <= mw.refreshTimeout {
		mw.refreshTimeout = mw.timeout + 30*time.Minute
	}
	return mw
}

// Timeout token valid time
func (sf *Auth) Timeout() time.Duration { return sf.timeout }

// MaxTimeout refresh timeout
func (sf *Auth) MaxTimeout() time.Duration { return sf.refreshTimeout }

// ParseToken parse token
func (sf *Auth) ParseToken(token string) (*Account, error) {
	return sf.provider.ParseToken(token)
}

// GenerateToken generate token
func (sf *Auth) GenerateToken(id string, acc *Account) (string, time.Time, error) {
	if acc.Issuer == "" {
		acc.Issuer = sf.issuer
	}
	return sf.provider.GenerateToken(id, acc, sf.timeout)
}

// GenerateRefreshToken generate refresh token
func (sf *Auth) GenerateRefreshToken(id string, acc *Account) (string, time.Time, error) {
	return sf.provider.GenerateRefreshToken(id, acc, sf.refreshTimeout)
}

// ExtractToken extract token from http request
func (sf *Auth) ExtractToken(r *http.Request) (string, error) {
	return sf.lookup.ExtractToken(r)
}

// ParseFromRequest parse token to account from http request
func (sf *Auth) ParseFromRequest(r *http.Request) (*Account, error) {
	token, err := sf.ExtractToken(r)
	if err != nil {
		return nil, err
	}
	return sf.ParseToken(token)
}
