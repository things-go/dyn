package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/things-go/dyn/core/metadata"
)

type Claims struct {
	Type     string            `json:"type,omitempty"`
	Scopes   []string          `json:"scopes,omitempty"`
	Metadata metadata.Metadata `json:"metadata,omitempty"`
	jwt.RegisteredClaims
}

type JwtConfig struct {
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
	Issuer          string
}

type JwtProvider struct {
	signingMethod jwt.SigningMethod
	encodeKey     any
	decodeKey     any
	issuer        string
}

func NewJwtProvider(c JwtConfig) (Provider, error) {
	var err error
	mw := &JwtProvider{
		issuer: c.Issuer,
	}

	usingAlgo := false
	switch c.Algorithm {
	case "RS256", "RS512", "RS384":
		usingAlgo = true
		mw.encodeKey, err = parsePrivateKey(c.PrivKey)
		if err != nil {
			return nil, ErrInvalidPrivKey
		}
		mw.decodeKey, err = parsePublicKey(c.PubKey)
		if err != nil {
			return nil, ErrInvalidPubKey
		}
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
	return mw, nil
}
func (sf *JwtProvider) GenerateToken(id string, acc *Account, timeout time.Duration) (string, time.Time, error) {
	return sf.generateToken(id, acc, timeout)
}
func (sf *JwtProvider) GenerateRefreshToken(id string, acc *Account, timeout time.Duration) (string, time.Time, error) {
	return sf.generateToken(id, acc, timeout)
}
func (sf *JwtProvider) generateToken(id string, acc *Account, timeout time.Duration) (string, time.Time, error) {
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

func (sf *JwtProvider) ParseToken(tokenString string) (*Account, error) {
	tk, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if sf.signingMethod != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}
		return sf.decodeKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				err = ErrInvalidToken
			case ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0:
				err = ErrTokenExpired
			default:
				err = ErrTokenParseFail
			}
		}
		return nil, err
	}
	if !tk.Valid {
		return nil, ErrInvalidToken
	}
	claims, ok := tk.Claims.(*Claims)
	if !ok || claims == nil {
		return nil, jwt.ErrTokenInvalidClaims
	}
	if claims.Issuer != sf.issuer {
		return nil, jwt.ErrTokenInvalidIssuer
	}
	if claims.Subject == "" {
		return nil, ErrInvalidToken
	}
	cc := tk.Claims.(*Claims)
	if cc.Metadata != nil && cc.ID != "" {
		cc.Metadata.Set(TokenUniqueId, cc.ID)
	}
	return &Account{
		Subject:  cc.Subject,
		Type:     cc.Type,
		Scopes:   cc.Scopes,
		Metadata: cc.Metadata,
	}, nil
}

// helper
func parsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	priv, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		priv, err = os.ReadFile(privateKey)
		if err != nil {
			priv = []byte(privateKey)
		}
	}
	return jwt.ParseRSAPrivateKeyFromPEM(priv)
}

func parsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	pub, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		pub, err = os.ReadFile(publicKey)
		if err != nil {
			pub = []byte(publicKey)
		}
	}
	return jwt.ParseRSAPublicKeyFromPEM(pub)
}
