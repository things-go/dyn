package json

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/things-go/dyn/core/encoding"
)

var defaultCodec encoding.Codec = codec{}

func init() {
	encoding.Register(defaultCodec)
}

func Name() string                               { return defaultCodec.Name() }
func Marshal(v interface{}) ([]byte, error)      { return defaultCodec.Marshal(v) }
func Unmarshal(data []byte, v interface{}) error { return defaultCodec.Unmarshal(data, v) }

// MarshalOptions is a configurable JSON format marshaller.
var MarshalOptions = protojson.MarshalOptions{
	EmitUnpopulated: true,
}

// UnmarshalOptions is a configurable JSON format parser.
var UnmarshalOptions = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}

// codec is a Codec implementation with json.
type codec struct{}

func (codec) Name() string { return "json" }
func (codec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return MarshalOptions.Marshal(m)
	default:
		return json.Marshal(m)
	}
}
func (codec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return UnmarshalOptions.Unmarshal(data, m)
	default:
		return json.Unmarshal(data, v)
	}
}
