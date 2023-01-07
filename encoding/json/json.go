package json

import (
	"encoding/json"
	"io"

	"github.com/things-go/dyn/encoding/codec"
)

// Codec is a Marshaler which marshals/unmarshals into/from JSON
// with the standard "encoding/json" package of Golang.
// Although it is generally faster for simple proto messages than JSONPb,
// it does not support advanced features of protobuf, e.g. map, oneof, ....
//
// The NewEncoder and NewDecoder types return *json.Encoder and
// *json.Decoder respectively.
type Codec struct{}

// ContentType always Returns "application/json".
func (*Codec) ContentType(_ interface{}) string {
	return "application/json"
}
func (j *Codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func (j *Codec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
func (j *Codec) NewDecoder(r io.Reader) codec.Decoder {
	return json.NewDecoder(r)
}
func (j *Codec) NewEncoder(w io.Writer) codec.Encoder {
	return json.NewEncoder(w)
}
func (j *Codec) Delimiter() []byte {
	return []byte("\n")
}
