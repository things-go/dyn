package toml

import (
	"io"

	"github.com/pelletier/go-toml/v2"

	"github.com/things-go/dyn/encoding/codec"
)

// Codec is a Codec implementation with yaml.
type Codec struct{}

// ContentType always Returns "application/yaml".
func (*Codec) ContentType(_ interface{}) string {
	return "application/toml"
}
func (*Codec) Marshal(v interface{}) ([]byte, error) {
	return toml.Marshal(v)
}
func (*Codec) Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}
func (*Codec) NewDecoder(r io.Reader) codec.Decoder {
	return toml.NewDecoder(r)
}
func (*Codec) NewEncoder(w io.Writer) codec.Encoder {
	return toml.NewEncoder(w)
}
