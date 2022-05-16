package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

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
	// if timeout <= maxTimeout, maxTimeout = timeout + 30 * time.Minute
	Timeout    time.Duration
	MaxTimeout time.Duration
	Lookup     string
}

// Auth provides a Json-Web-Token authentication implementation.
type Auth struct {
	provider   Provider
	timeout    time.Duration
	maxTimeout time.Duration
	lookup     *Lookup
}

// New auth with Config
func New(p Provider, c Config) *Auth {
	mw := &Auth{
		provider:   p,
		timeout:    c.Timeout,
		maxTimeout: c.MaxTimeout,
		lookup:     NewLookup(c.Lookup),
	}
	if mw.timeout <= mw.maxTimeout {
		mw.maxTimeout = mw.timeout + 30*time.Minute
	}
	return mw
}

func (sf *Auth) Timeout() time.Duration    { return sf.timeout }
func (sf *Auth) MaxTimeout() time.Duration { return sf.maxTimeout }

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

func (sf *Auth) ParseToken(token string) (*Account, error) {
	return sf.provider.ParseToken(token)
}

func (sf *Auth) GenerateToken(id string, acc *Account) (string, time.Time, error) {
	return sf.provider.GenerateToken(id, acc, sf.timeout)
}

func (sf *Auth) GenerateRefreshToken(id string, acc *Account) (string, time.Time, error) {
	return sf.provider.GenerateRefreshToken(id, acc, sf.maxTimeout)
}

// Option is Middleware option.
type Option func(*options)

// options is a Middleware option
type options struct {
	skip                 func(c *gin.Context) bool
	unauthorizedFallback func(*gin.Context, error)
}

// WithSkip set skip func
func WithSkip(f func(c *gin.Context) bool) Option {
	return func(o *options) {
		if f != nil {
			o.skip = f
		}
	}
}

// WithUnauthorizedFallback sets the fallback handler when requests are unauthorized.
func WithUnauthorizedFallback(f func(c *gin.Context, err error)) Option {
	return func(o *options) {
		if f != nil {
			o.unauthorizedFallback = f
		}
	}
}

func (sf *Auth) Middleware(opts ...Option) gin.HandlerFunc {
	o := &options{
		unauthorizedFallback: func(c *gin.Context, err error) {
			c.String(http.StatusUnauthorized, err.Error())
		},
		skip: func(c *gin.Context) bool { return false },
	}
	for _, opt := range opts {
		opt(o)
	}
	return func(c *gin.Context) {
		if !o.skip(c) {
			acc, err := sf.ParseFromRequest(c.Request)
			if err != nil {
				o.unauthorizedFallback(c, err)
				c.Abort()
				return
			}
			c.Request = c.Request.WithContext(NewContext(c.Request.Context(), acc))
		}
		c.Next()
	}
}
