package auth

import (
	"errors"
)

var (
	// ErrMissingToken can be thrown by follow
	// if authing with a HTTP header, the Auth header needs to be set
	// if authing with URL Query, the query token variable is empty
	// if authing with a cookie, the token cookie is empty
	ErrMissingToken = errors.New("no token present in request")
	// ErrInvalidToken indicates auth token is invalid
	ErrInvalidToken = errors.New("invalid token provided")
	// ErrTokenExpired indicates auth token is expired
	ErrTokenExpired = errors.New("token expired")
	// ErrInvalidToken indicates auth token is invalid
	ErrTokenParseFail = errors.New("parse JWT token failed")
	// ErrInvalidSigningAlgorithm indicates signing algorithm is invalid,
	// needs to be HS256, HS384, HS512, RS256, RS384 or RS512
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")

	// ErrInvalidPubKey indicates the the given public key is invalid
	ErrInvalidPubKey = errors.New("public key invalid")
	// ErrInvalidPrivKey indicates that the given private key is invalid
	ErrInvalidPrivKey = errors.New("private key invalid")
	// ErrMissingSecretKey indicates Secret key is required
	ErrMissingSecretKey = errors.New("secret key is required")
)
