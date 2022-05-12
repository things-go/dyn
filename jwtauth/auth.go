package jwtauth

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
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

type Claims struct {
	Type     string            `json:"type,omitempty"`
	Scopes   []string          `json:"scopes,omitempty"`
	Metadata metadata.Metadata `json:"metadata,omitempty"`
	jwt.RegisteredClaims
}

// Config Auth config
type Config struct {
	// 支持签名算法: HS256, HS384, HS512, RS256, RS384 or RS512
	// Optional, Default HS256.
	Algorithm string
	// Secret key used for signing.
	// Required, HS256, HS384, HS512.
	Key []byte
	// Private key for asymmetric algorithms,
	// Public key for asymmetric algorithms
	// Required, RS256, RS384 or RS512.
	PrivKey, PubKey string
	// if timeout <= maxTimeout, maxTimeout = timeout + 30 * time.Minute
	Timeout    time.Duration
	MaxTimeout time.Duration
	Issuer     string
	Lookup     string
}

// Auth provides a Json-Web-Token authentication implementation.
type Auth struct {
	signingMethod jwt.SigningMethod
	encodeKey     interface{}
	decodeKey     interface{}
	timeout       time.Duration
	maxTimeout    time.Duration
	issuer        string
	lookup        *Lookup
}

// New auth with Config
func New(c Config) (*Auth, error) {
	mw := &Auth{
		timeout:    c.Timeout,
		maxTimeout: c.MaxTimeout,
		issuer:     c.Issuer,
		lookup:     NewLookup(c.Lookup),
	}

	usingAlgo := false
	switch c.Algorithm {
	case "RS256", "RS512", "RS384":
		usingAlgo = true
		mw.encodeKey = c.PrivKey
		mw.decodeKey = c.PubKey
	case "HS256", "HS512", "HS384":
	default:
		c.Algorithm = "HS256"
	}
	mw.signingMethod = jwt.GetSigningMethod(c.Algorithm)

	if !usingAlgo {
		if c.Key == nil {
			return nil, ErrMissingSecretKey
		}
		mw.encodeKey = c.Key
		mw.decodeKey = c.Key
	}
	if mw.timeout <= mw.maxTimeout {
		mw.maxTimeout = mw.timeout + 30*time.Minute
	}
	return mw, nil
}

func (sf *Auth) Timeout() time.Duration    { return sf.timeout }
func (sf *Auth) MaxTimeout() time.Duration { return sf.maxTimeout }

func (sf *Auth) Parse(tokenString string) (jwt.Claims, error) {
	tk, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if sf.signingMethod != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}
		return sf.decodeKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				err = ErrInvalidToken
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				err = ErrTokenExpired
			} else {
				err = ErrTokenParseFail
			}
		}
		return nil, err
	}
	claims, ok := tk.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	if claims == nil || claims.Subject == "" {
		return nil, errors.New("invalid subject")
	}
	subject, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || subject == 0 {
		return nil, errors.New("invalid subject")
	}
	return tk.Claims, nil
}

func (sf *Auth) ExtractToken(r *http.Request) (string, error) {
	return sf.lookup.ExtractToken(r)
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
	token, err := jwt.NewWithClaims(sf.signingMethod, &Claims{
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
	}).SignedString(sf.encodeKey)
	return token, expiresAt, err
}

// Option is jwt option.
type Option func(*options)

// options is a jwt option
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
			token, err := sf.lookup.ExtractToken(c.Request)
			if err != nil {
				o.unauthorizedFallback(c, err)
				c.Abort()
				return
			}
			claims, err := sf.Parse(token)
			if err != nil {
				o.unauthorizedFallback(c, err)
				c.Abort()
				return
			}
			ctx := NewContext(c.Request.Context(), claims)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}

func parsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	var err error
	var priv []byte

	if strings.HasPrefix(privateKey, "-----BEGIN RSA PRIVATE KEY-----") ||
		strings.HasPrefix(privateKey, "-----BEGIN PRIVATE KEY-----") {
		priv = []byte(privateKey)
	} else {
		priv, err = base64.StdEncoding.DecodeString(privateKey)
		if err != nil {
			return nil, err
		}
	}
	return jwt.ParseRSAPrivateKeyFromPEM(priv)
}

// parsePublicKey parses a public key
func parsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	var err error
	var pub []byte

	if strings.HasPrefix(publicKey, "-----BEGIN RSA PUBLIC KEY-----") ||
		strings.HasPrefix(publicKey, "-----BEGIN PUBLIC KEY-----") ||
		strings.HasPrefix(publicKey, "-----BEGIN CERTIFICATE-----") {
		pub = []byte(publicKey)
	} else {
		pub, err = base64.StdEncoding.DecodeString(publicKey)
		if err != nil {
			return nil, err
		}
	}
	return jwt.ParseRSAPublicKeyFromPEM(pub)
}
