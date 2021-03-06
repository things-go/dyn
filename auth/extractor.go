package auth

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// Extractor is an interface for extracting a token from an HTTP request.
// The ExtractToken method should return a token string or an error.
// If no token is present, you must return ErrNoTokenInRequest.
type Extractor interface {
	ExtractToken(*http.Request) (string, error)
}

// MultiExtractor tries Extractors in order until one returns a token string or an error occurs
type MultiExtractor []Extractor

func (e MultiExtractor) ExtractToken(req *http.Request) (string, error) {
	// loop over header names and return the first one that contains data
	for _, extractor := range e {
		if tok, err := extractor.ExtractToken(req); tok != "" {
			return tok, nil
		} else if !errors.Is(err, ErrMissingToken) {
			return "", err
		}
	}
	return "", ErrMissingToken
}

// HeaderExtractor is an extractor for finding a token in a header.
// Looks at each specified header in order until there's a match
type HeaderExtractor struct {
	// The key of the header
	// Required
	Key string
	// Strips 'Bearer ' prefix from bearer token string.
	// Possible value is "Bearer"
	// Optional
	Prefix string
}

func (e HeaderExtractor) ExtractToken(r *http.Request) (string, error) {
	// loop over header names and return the first one that contains data
	return stripHeadValuePrefixFromTokenString(e.Prefix)(r.Header.Get(e.Key))
}

// ArgumentExtractor extracts a token from request arguments.  This includes a POSTed form or
// GET URL arguments.
// This extractor calls `ParseMultipartForm` on the request
type ArgumentExtractor string

func (e ArgumentExtractor) ExtractToken(r *http.Request) (string, error) {
	// Make sure form is parsed
	_ = r.ParseMultipartForm(10e6)

	tk := strings.TrimSpace(r.Form.Get(string(e)))
	if tk != "" {
		return tk, nil
	}
	return "", ErrMissingToken
}

// CookieExtractor extracts a token from cookie.
type CookieExtractor string

func (e CookieExtractor) ExtractToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(string(e))
	if err != nil {
		return "", ErrMissingToken
	}
	val, _ := url.QueryUnescape(cookie.Value)
	if val = strings.TrimSpace(val); val != "" {
		return val, nil
	}
	return "", ErrMissingToken
}

// Strips like 'Bearer ' prefix from token string with header name
func stripHeadValuePrefixFromTokenString(prefix string) func(string) (string, error) {
	return func(tok string) (string, error) {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			return "", ErrMissingToken
		}
		l := len(prefix)
		if l == 0 {
			return tok, nil
		}
		// Should be a bearer token
		if len(tok) > l && strings.EqualFold(tok[:l], prefix) {
			if tok = strings.TrimSpace(tok[l+1:]); tok != "" {
				return tok, nil
			}
		}
		return "", ErrMissingToken
	}
}
