// Package proto defines the protobuf codec. Importing this package will
// register the codec.
package proto

import (
	"github.com/things-go/dyn/core/encoding"

	"google.golang.org/protobuf/proto"
)

func init() {
	encoding.Register(codec{})
}

// codec is a Codec implementation with protobuf.
type codec struct{}

func (codec) Name() string { return "proto" }
func (codec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}
func (codec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
