package encoding

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

type invalidCodec struct{}

func (c invalidCodec) Marshal(v interface{}) ([]byte, error) {
	panic("implement me")
}

func (c invalidCodec) Unmarshal(data []byte, v interface{}) error {
	panic("implement me")
}

func (c invalidCodec) Name() string {
	return ""
}

// codec is a Codec implementation with xml.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func (codec) Name() string { return "xml" }

func TestRegisterCodec(t *testing.T) {
	require.Panics(t, func() { Register(nil) })
	require.Panics(t, func() { Register(invalidCodec{}) })

	cdc := codec{}
	Register(cdc)
	require.Equal(t, cdc, GetCodec(cdc.Name()))
	require.Equal(t, []string{cdc.Name()}, Codecs())

	require.Panics(t, func() { Register(cdc) }, "should register called twice for driver")
}
