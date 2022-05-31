package auth

import (
	"net/http"
	"time"

	"github.com/things-go/dyn/core/metadata"
)

type Provider interface {
	GenerateToken(id string, acc *Account, timeout time.Duration) (string, time.Time, error)
	GenerateRefreshToken(id string, acc *Account, timeout time.Duration) (string, time.Time, error)
	ParseToken(token string) (*Account, error)
}

type Account struct {
	Subject  string            `json:"subject,omitempty"`
	Type     string            `json:"type,omitempty"`
	Scopes   []string          `json:"scopes,omitempty"`
	Metadata metadata.Metadata `json:"metadata,omitempty"`
}

// Config Auth config
type Config struct {
	// if timeout <= refreshTimeout, refreshTimeout = timeout + 30 * time.Minute
	Timeout        time.Duration
	RefreshTimeout time.Duration
	Lookup         string
}

// Auth provides a Json-Web-Token authentication implementation.
type Auth struct {
	provider       Provider
	timeout        time.Duration
	refreshTimeout time.Duration
	lookup         *Lookup
}

// New auth with Config
func New(p Provider, c Config) *Auth {
	mw := &Auth{
		provider:       p,
		timeout:        c.Timeout,
		refreshTimeout: c.RefreshTimeout,
		lookup:         NewLookup(c.Lookup),
	}
	if mw.timeout <= mw.refreshTimeout {
		mw.refreshTimeout = mw.timeout + 30*time.Minute
	}
	return mw
}

func (sf *Auth) Timeout() time.Duration    { return sf.timeout }
func (sf *Auth) MaxTimeout() time.Duration { return sf.refreshTimeout }

func (sf *Auth) ParseToken(token string) (*Account, error) {
	return sf.provider.ParseToken(token)
}

func (sf *Auth) GenerateToken(id string, acc *Account) (string, time.Time, error) {
	return sf.provider.GenerateToken(id, acc, sf.timeout)
}

func (sf *Auth) GenerateRefreshToken(id string, acc *Account) (string, time.Time, error) {
	return sf.provider.GenerateRefreshToken(id, acc, sf.refreshTimeout)
}

func (sf *Auth) ExtractToken(r *http.Request) (string, error) {
	return sf.lookup.ExtractToken(r)
}

func (sf *Auth) ParseFromRequest(r *http.Request) (*Account, error) {
	token, err := sf.ExtractToken(r)
	if err != nil {
		return nil, err
	}
	return sf.ParseToken(token)
}
