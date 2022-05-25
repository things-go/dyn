package yaml

import (
	"gopkg.in/yaml.v3"

	"github.com/things-go/dyn/core/encoding"
)

func init() {
	encoding.Register(codec{})
}

// codec is a Codec implementation with yaml.
type codec struct{}

func (codec) Name() string { return "yaml" }
func (codec) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}
func (codec) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
