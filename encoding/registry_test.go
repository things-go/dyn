package encoding

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/things-go/dyn/encoding/codec"
)

var marshalers = []dummyMarshaler{0, 1}

func init() {
	RegisterMarshaler("application/x-0", &marshalers[0])
	RegisterMarshaler("application/x-1", &marshalers[1])
}

func TestForRequest_Wildcard(t *testing.T) {
	r, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf(`http.NewRequest("GET", "http://example.com", nil) failed with %v; want success`, err)
	}

	r.Header.Set("Accept", "application/unknown")
	r.Header.Set("Content-Type", "application/unknown")
	in := MarshalerFromRequest(r, true)
	if _, ok := in.(*HTTPBodyCodec); !ok {
		t.Errorf("in = %#v; want a HTTPBodyCodec", in)
	}
	out := MarshalerFromRequest(r, false)
	if _, ok := out.(*HTTPBodyCodec); !ok {
		t.Errorf("out = %#v; want a HTTPBodyCodec", in)
	}
}

func TestForRequest_NotWildcard(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		accept      string
		wantIn      codec.Marshaler
		wantOut     codec.Marshaler
	}{
		// You can specify a marshaler for a specific MIME type.
		// The output marshaler follows the input one unless specified.
		{
			name:        "",
			contentType: "application/x-0",
			accept:      "application/x-0",
			wantIn:      &marshalers[0],
			wantOut:     &marshalers[0],
		},
		// You can also separately specify an output marshaler
		{
			name:        "",
			contentType: "application/x-0",
			accept:      "application/x-1",
			wantIn:      &marshalers[0],
			wantOut:     &marshalers[1],
		},
		{
			name:        "",
			contentType: "application/x-1; charset=UTF-8",
			accept:      "application/x-1",
			wantIn:      &marshalers[1],
			wantOut:     &marshalers[1],
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := http.NewRequest("GET", "http://example.com", nil)
			if err != nil {
				t.Fatalf(`http.NewRequest("GET", "http://example.com", nil) failed with %v; want success`, err)
			}
			r.Header.Set("Accept", test.accept)
			r.Header.Set("Content-Type", test.contentType)
			in := MarshalerFromRequest(r, true)
			if got, want := in, test.wantIn; got != want {
				t.Errorf("in = %#v; want %#v", got, want)
			}
			out := MarshalerFromRequest(r, false)
			if got, want := out, test.wantOut; got != want {
				t.Errorf("out = %#v; want %#v", got, want)
			}
		})
	}
}

type dummyMarshaler int

func (dummyMarshaler) ContentType(_ interface{}) string { return "" }
func (dummyMarshaler) Marshal(interface{}) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (dummyMarshaler) Unmarshal([]byte, interface{}) error {
	return errors.New("not implemented")
}

func (dummyMarshaler) NewDecoder(r io.Reader) codec.Decoder {
	return dummyDecoder{}
}
func (dummyMarshaler) NewEncoder(w io.Writer) codec.Encoder {
	return dummyEncoder{}
}

func (m dummyMarshaler) GoString() string {
	return fmt.Sprintf("dummyMarshaler(%d)", m)
}

type dummyDecoder struct{}

func (dummyDecoder) Decode(interface{}) error {
	return errors.New("not implemented")
}

type dummyEncoder struct{}

func (dummyEncoder) Encode(interface{}) error {
	return errors.New("not implemented")
}
