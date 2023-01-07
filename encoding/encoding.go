package encoding

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/things-go/dyn/encoding/codec"
)

type ctxUriKey struct{}

const defaultMemory = 32 << 20

// Content-Type MIME of the most common data formats.
const (
	// MIMEWildcard is the fallback MIME type used for requests which do not match
	// a registered MIME type.
	MIMEWildcard = "*"
	// MIMEURI is sepcial form query.
	MIMEQuery = "__MIME__/QUERY"
	// MIMEURI is sepcial form uri.
	MIMEURI = "__MIME__/URI"

	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
	MIMETOML              = "application/toml"
)

func Bind(req *http.Request, v any) error {
	if req.Method == http.MethodGet {
		return BindQuery(req, v)
	}
	marshaller := MarshalerFromRequest(req, true)
	if contentType := marshaller.ContentType(req.Body); contentType == MIMEMultipartPOSTForm {
		m, ok := marshaller.(codec.FormMarshaler)
		if !ok {
			return fmt.Errorf("not supported marshaller(%v)", contentType)
		}
		if err := req.ParseMultipartForm(defaultMemory); err != nil {
			return err
		}
		return m.Decode(req.MultipartForm.Value, v)
	}
	return marshaller.NewDecoder(req.Body).
		Decode(v)
}

func BindQuery(req *http.Request, v any) error {
	marshaller := GetMarshaler(MIMEQuery)
	m, ok := marshaller.(codec.FormMarshaler)
	if !ok {
		return fmt.Errorf("not supported marshaller(%v)", MIMEQuery)
	}
	return m.Decode(req.URL.Query(), v)
}

func BindUri(req *http.Request, v any) error {
	raws := FromRequestUri(req)
	if raws == nil {
		return errors.New("must be request with uri in context")
	}
	marshaller := GetMarshaler(MIMEURI)
	m, ok := marshaller.(codec.FormMarshaler)
	if !ok {
		return fmt.Errorf("not supported marshaller(%v)", MIMEURI)
	}
	return m.Decode(raws, v)
}

func RequestWithUri(req *http.Request, uri url.Values) *http.Request {
	if uri == nil {
		uri = url.Values{}
	}
	ctx := context.WithValue(req.Context(), ctxUriKey{}, uri)
	return req.WithContext(ctx)
}

func FromRequestUri(req *http.Request) url.Values {
	if rv := req.Context().Value(ctxUriKey{}); rv != nil {
		return rv.(url.Values)
	}
	return nil
}
