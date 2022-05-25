// Package proto defines the protobuf codec. Importing this package will
// register the codec.
package proto

import (
	"github.com/things-go/dyn/core/encoding"

	"google.golang.org/protobuf/proto"
)

var defaultCodec = New()

func init() {
	encoding.Register(defaultCodec)
}

func Name() string                               { return defaultCodec.Name() }
func Marshal(v interface{}) ([]byte, error)      { return defaultCodec.Marshal(v) }
func Unmarshal(data []byte, v interface{}) error { return defaultCodec.Unmarshal(data, v) }

// New returns a new Codec
func New() Codec { return Codec{} }

// Codec is a Codec implementation with protobuf.
type Codec struct{}

func (Codec) Name() string { return "proto" }
func (Codec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}
func (Codec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
