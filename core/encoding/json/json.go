package json

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/things-go/dyn/core/encoding"
)

var defaultCodec = New()

func init() {
	encoding.Register(defaultCodec)
}

func Name() string                               { return defaultCodec.Name() }
func Marshal(v interface{}) ([]byte, error)      { return defaultCodec.Marshal(v) }
func Unmarshal(data []byte, v interface{}) error { return defaultCodec.Unmarshal(data, v) }

type Option func(*Codec)

func WithMarshalOptions(marshalOpts protojson.MarshalOptions) Option {
	return func(c *Codec) {
		c.marshalOpts = marshalOpts
	}
}

func WithUnmarshalOptions(unmarshalOpts protojson.UnmarshalOptions) Option {
	return func(c *Codec) {
		c.unmarshalOpts = unmarshalOpts
	}
}

// New returns a new Codec
func New(opts ...Option) Codec {
	codec := Codec{
		protojson.MarshalOptions{EmitUnpopulated: true},
		protojson.UnmarshalOptions{DiscardUnknown: true},
	}
	for _, opt := range opts {
		opt(&codec)
	}
	return codec
}

// Codec is a Codec implementation with json.
type Codec struct {
	marshalOpts   protojson.MarshalOptions
	unmarshalOpts protojson.UnmarshalOptions
}

func (Codec) Name() string { return "json" }
func (c Codec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return c.marshalOpts.Marshal(m)
	default:
		return json.Marshal(m)
	}
}
func (c Codec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return c.unmarshalOpts.Unmarshal(data, m)
	default:
		return json.Unmarshal(data, v)
	}
}
