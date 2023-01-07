package xml

import (
	"encoding/xml"
	"io"

	"github.com/things-go/dyn/encoding/codec"
)

// Codec is a Codec implementation with xml.
type Codec struct{}

// ContentType always Returns "application/xml".
func (*Codec) ContentType(_ interface{}) string {
	return "application/xml"
}
func (*Codec) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}
func (*Codec) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}
func (*Codec) NewEncoder(w io.Writer) codec.Encoder {
	return xml.NewEncoder(w)
}
func (*Codec) NewDecoder(r io.Reader) codec.Decoder {
	return xml.NewDecoder(r)
}
