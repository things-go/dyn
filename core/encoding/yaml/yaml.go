package yaml

import (
	"gopkg.in/yaml.v3"

	"github.com/things-go/dyn/core/encoding"
)

var defaultCodec encoding.Codec = codec{}

func init() {
	encoding.Register(codec{})
}

func Name() string                               { return defaultCodec.Name() }
func Marshal(v interface{}) ([]byte, error)      { return defaultCodec.Marshal(v) }
func Unmarshal(data []byte, v interface{}) error { return defaultCodec.Unmarshal(data, v) }

// codec is a Codec implementation with yaml.
type codec struct{}

func (codec) Name() string { return "yaml" }
func (codec) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}
func (codec) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
