package xml

import (
	"encoding/xml"

	"github.com/things-go/dyn/core/encoding"
)

func init() {
	encoding.Register(codec{})
}

// codec is a Codec implementation with xml.
type codec struct{}

func (codec) Name() string                               { return "xml" }
func (codec) Marshal(v interface{}) ([]byte, error)      { return xml.Marshal(v) }
func (codec) Unmarshal(data []byte, v interface{}) error { return xml.Unmarshal(data, v) }
