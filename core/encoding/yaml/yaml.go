package yaml

import (
	"gopkg.in/yaml.v3"

	"github.com/things-go/dyn/core/encoding"
)

var defaultCodec = New()

func init() {
	encoding.Register(defaultCodec)
}

func Name() string                       { return defaultCodec.Name() }
func Marshal(v any) ([]byte, error)      { return defaultCodec.Marshal(v) }
func Unmarshal(data []byte, v any) error { return defaultCodec.Unmarshal(data, v) }

// Codec is a Codec implementation with yaml.
type Codec struct{}

// New returns a new Codec
func New() Codec { return Codec{} }

func (Codec) Name() string { return "yaml" }
func (Codec) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}
func (Codec) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}
