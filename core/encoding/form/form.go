package form

import (
	"net/url"
	"reflect"

	"github.com/things-go/dyn/core/encoding"

	"github.com/go-playground/form/v4"
	"google.golang.org/protobuf/proto"
)

var defaultCodec = New()

func init() {
	encoding.Register(defaultCodec)
}
func ReplaceDefaultCodec(codec Codec)            { defaultCodec = codec }
func Name() string                               { return defaultCodec.Name() }
func Marshal(v interface{}) ([]byte, error)      { return defaultCodec.Marshal(v) }
func Unmarshal(data []byte, v interface{}) error { return defaultCodec.Unmarshal(data, v) }
func Encode(v interface{}) (url.Values, error)   { return defaultCodec.Encode(v) }
func Decode(vs url.Values, v interface{}) error  { return defaultCodec.Decode(vs, v) }

type Codec struct {
	encoder      *form.Encoder
	decoder      *form.Decoder
	disableProto bool
}

type Option func(*Codec)

func WithEncoder(encoder *form.Encoder) Option {
	return func(c *Codec) {
		if encoder != nil {
			c.encoder = encoder
		}
	}
}

func WithDecoder(decoder *form.Decoder) Option {
	return func(c *Codec) {
		if decoder != nil {
			c.decoder = decoder
		}

	}
}

func WithDisableProto() Option {
	return func(c *Codec) {
		c.disableProto = true
	}
}

// New returns a new Codec,
// default tag name is "form",
// proto use protoJSON tag
func New(opts ...Option) Codec {
	encoder := form.NewEncoder()
	encoder.SetTagName("form")
	decoder := form.NewDecoder()
	decoder.SetTagName("form")
	codec := Codec{
		encoder,
		decoder,
		false,
	}
	for _, opt := range opts {
		opt(&codec)
	}
	return codec
}

func (Codec) Name() string { return "x-www-form-urlencoded" }

func (c Codec) Marshal(v interface{}) ([]byte, error) {
	vs, err := c.Encode(v)
	if err != nil {
		return nil, err
	}
	return []byte(vs.Encode()), nil
}

func (c Codec) Unmarshal(data []byte, v interface{}) error {
	vs, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}
	return c.Decode(vs, v)
}

func (c Codec) Encode(v interface{}) (url.Values, error) {
	var vs url.Values
	var err error

	if c.disableProto {
		vs, err = c.encoder.Encode(v)
	} else if m, ok := v.(proto.Message); ok {
		vs, err = EncodeValues(m)
	} else {
		vs, err = c.encoder.Encode(v)
	}
	if err != nil {
		return nil, err
	}
	for k, vv := range vs {
		if len(vv) == 0 {
			delete(vs, k)
		}
	}
	return vs, nil
}

func (c Codec) Decode(vs url.Values, v interface{}) error {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	if !c.disableProto {
		if m, ok := v.(proto.Message); ok {
			return DecodeValues(m, vs)
		}
		if m, ok := reflect.Indirect(reflect.ValueOf(v)).Interface().(proto.Message); ok {
			return DecodeValues(m, vs)
		}
	}
	return c.decoder.Decode(v, vs)
}
