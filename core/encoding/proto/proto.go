// Package proto defines the protobuf codec. Importing this package will
// register the codec.
package proto

import (
	"github.com/things-go/dyn/core/encoding"

	"google.golang.org/protobuf/proto"
)

var defaultCodec encoding.Codec = codec{}

func init() {
	encoding.Register(codec{})
}

func Name() string                               { return defaultCodec.Name() }
func Marshal(v interface{}) ([]byte, error)      { return defaultCodec.Marshal(v) }
func Unmarshal(data []byte, v interface{}) error { return defaultCodec.Unmarshal(data, v) }

// codec is a Codec implementation with protobuf.
type codec struct{}

func (codec) Name() string { return "proto" }
func (codec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}
func (codec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
