package xml

import (
	"encoding/xml"

	"github.com/things-go/dyn/core/encoding"
)

var defaultCodec = New()

func init() {
	encoding.Register(defaultCodec)
}

func Name() string                       { return defaultCodec.Name() }
func Marshal(v any) ([]byte, error)      { return defaultCodec.Marshal(v) }
func Unmarshal(data []byte, v any) error { return defaultCodec.Unmarshal(data, v) }

// New returns a new Codec
func New() Codec { return Codec{} }

// Codec is a Codec implementation with xml.
type Codec struct{}

func (Codec) Name() string { return "xml" }
func (Codec) Marshal(v any) ([]byte, error) {
	return xml.Marshal(v)
}
func (Codec) Unmarshal(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}
