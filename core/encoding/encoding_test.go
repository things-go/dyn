package encoding

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

type invalidCodec struct{}

func (c invalidCodec) Name() string                       { return "" }
func (c invalidCodec) Marshal(v any) ([]byte, error)      { panic("implement me") }
func (c invalidCodec) Unmarshal(data []byte, v any) error { panic("implement me") }

type codec struct{}

func (codec) Name() string                       { return "xml" }
func (codec) Marshal(v any) ([]byte, error)      { return xml.Marshal(v) }
func (codec) Unmarshal(data []byte, v any) error { return xml.Unmarshal(data, v) }

func TestRegisterCodec(t *testing.T) {
	require.Panics(t, func() { Register(nil) })
	require.Panics(t, func() { Register(invalidCodec{}) })

	cdc := codec{}
	Register(cdc)
	require.Equal(t, cdc, GetCodec(cdc.Name()))
	require.Equal(t, []string{cdc.Name()}, Codecs())
}
