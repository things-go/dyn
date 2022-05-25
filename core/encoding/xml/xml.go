package xml

import (
	"encoding/xml"

	"github.com/things-go/dyn/core/encoding"
)

var defaultCodec encoding.Codec = codec{}

func init() {
	encoding.Register(codec{})
}

func Name() string                               { return defaultCodec.Name() }
func Marshal(v interface{}) ([]byte, error)      { return defaultCodec.Marshal(v) }
func Unmarshal(data []byte, v interface{}) error { return defaultCodec.Unmarshal(data, v) }

// codec is a Codec implementation with xml.
type codec struct{}

func (codec) Name() string { return "xml" }
func (codec) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}
func (codec) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}
