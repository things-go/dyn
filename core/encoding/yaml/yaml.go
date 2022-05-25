package yaml

import (
	"gopkg.in/yaml.v3"

	"github.com/things-go/dyn/core/encoding"
)

var defaultCodec = New()

func init() {
	encoding.Register(defaultCodec)
}

func Name() string                               { return defaultCodec.Name() }
func Marshal(v interface{}) ([]byte, error)      { return defaultCodec.Marshal(v) }
func Unmarshal(data []byte, v interface{}) error { return defaultCodec.Unmarshal(data, v) }

func New() Codec { return Codec{} }

// Codec is a Codec implementation with yaml.
type Codec struct{}

func (Codec) Name() string { return "yaml" }
func (Codec) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}
func (Codec) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
