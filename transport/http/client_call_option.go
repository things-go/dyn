package http

import (
	"net/http"
)

// CallSettings allow fine-grained control over how calls are made.
type CallSettings struct {
	// Content-Type
	contentType string
	// Accept
	accept string
	// custom header
	header http.Header
	// Path overwrite api call
	Path string
	// no auth
	noAuth bool
}

// CallOption is an option used by Invoke to control behaviors of RPC calls.
// CallOption works by modifying relevant fields of CallSettings.
type CallOption func(cs *CallSettings)

// WithCoContentType use encoding.MIMExxx
func WithCoContentType(contentType string) CallOption {
	return func(cs *CallSettings) {
		cs.contentType = contentType
	}
}

// WithCoAccept use encoding.MIMExxx
func WithCoAccept(accept string) CallOption {
	return func(cs *CallSettings) {
		cs.accept = accept
	}
}

// WithCoPath
func WithCoPath(path string) CallOption {
	return func(cs *CallSettings) {
		cs.Path = path
	}
}

// WithCoHeader
func WithCoHeader(k, v string) CallOption {
	return func(cs *CallSettings) {
		cs.header.Add(k, v)
	}
}

// WithCoNoAuth
func WithCoNoAuth() CallOption {
	return func(cs *CallSettings) {
		cs.noAuth = true
	}
}
