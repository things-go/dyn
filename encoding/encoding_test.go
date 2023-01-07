package encoding

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/things-go/dyn/encoding/internal/examplepb"
)

type TestMode struct {
	Id   string `json:"id" yaml:"id" xml:"id" toml:"id" msgpack:"id"`
	Name string `json:"name" yaml:"name" xml:"name" toml:"name" msgpack:"name"`
}

var protoMessage = &examplepb.ABitOfEverything{
	SingleNested:        &examplepb.ABitOfEverything_Nested{},
	RepeatedStringValue: nil,
	MappedStringValue:   nil,
	MappedNestedValue:   nil,
	RepeatedEnumValue:   nil,
	TimestampValue:      &timestamppb.Timestamp{},
	Uuid:                "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
	Nested: []*examplepb.ABitOfEverything_Nested{
		{
			Name:   "foo",
			Amount: 12345,
		},
	},
	Uint64Value: 0xFFFFFFFFFFFFFFFF,
	EnumValue:   examplepb.NumericEnum_ONE,
	OneofValue: &examplepb.ABitOfEverything_OneofString{
		OneofString: "bar",
	},
	MapValue: map[string]examplepb.NumericEnum{
		"a": examplepb.NumericEnum_ONE,
		"b": examplepb.NumericEnum_ZERO,
	},
}

func TestBind(t *testing.T) {
	tests := []struct {
		name    string
		genReq  func() (*http.Request, error)
		want    any
		wantErr bool
	}{
		{
			"default: marshaler",
			func() (*http.Request, error) {
				marshaler := GetMarshaler(MIMEWildcard)

				b, err := marshaler.Marshal(&examplepb.Complex{
					Id:     11,
					Uint32: wrapperspb.UInt32(1234),
					Bool:   wrapperspb.Bool(true),
				})
				if err != nil {
					return nil, err
				}
				r, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(b))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/unknown")
				return r, nil
			},
			&examplepb.Complex{
				Id:     11,
				Uint32: wrapperspb.UInt32(1234),
				Bool:   wrapperspb.Bool(true),
			},
			false,
		},
		{
			"form - application/x-www-form-urlencoded",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader([]byte(`id=foo&name=bar`)))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"form - method get so it query",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com?id=foo&name=bar", nil)
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"form - MultipartForm",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "http://example.com", nil)
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "multipart/form-data")
				r.MultipartForm = &multipart.Form{
					Value: map[string][]string{
						"id":   {"foo"},
						"name": {"bar"},
					},
					File: nil,
				}
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"json",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader([]byte(`{"id":"foo"}`)))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/json")
				return r, nil
			},
			&examplepb.SimpleMessage{
				Id: "foo",
			},
			false,
		},
		{
			"proto",
			func() (*http.Request, error) {
				buf := &bytes.Buffer{}

				m := GetMarshaler("application/x-protobuf")
				err := m.NewEncoder(buf).Encode(protoMessage)
				if err != nil {
					return nil, err
				}
				r, err := http.NewRequest(http.MethodPost, "http://example.com", buf)
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-protobuf")
				return r, nil
			},
			protoMessage,
			false,
		},
		{
			"yaml",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader([]byte("id: foo\nname: bar")))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-yaml")
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"xml",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader([]byte("<TestMode><id>foo</id><name>bar</name></TestMode>")))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/xml")
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"toml",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader([]byte("id=\"foo\"\nname=\"bar\"")))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/toml")
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"msgpack",
			func() (*http.Request, error) {
				buf := &bytes.Buffer{}

				m := GetMarshaler("application/x-msgpack")
				err := m.NewEncoder(buf).Encode(&TestMode{
					Id:   "foo",
					Name: "bar",
				})
				if err != nil {
					return nil, err
				}

				r, err := http.NewRequest(http.MethodPost, "http://example.com", buf)
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-msgpack")
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.genReq()
			if err != nil {
				t.Errorf("genReq() error = %v", err)
			}
			got := alloc(reflect.TypeOf(tt.want))
			if err = Bind(req, got.Interface()); (err != nil) != tt.wantErr {
				t.Errorf("Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, ok := tt.want.(proto.Message); ok {
				if diff := proto.Equal(got.Interface().(proto.Message), tt.want.(proto.Message)); !diff {
					t.Errorf("got = %v, want %v", got, tt.want)
				}
			} else {
				require.Equal(t, got.Interface(), tt.want)
			}
		})
	}
}

func TestBindQuery(t *testing.T) {
	tests := []struct {
		name    string
		genReq  func() (*http.Request, error)
		want    any
		wantErr bool
	}{
		{
			"form - no proto",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com?id=foo&name=bar", nil)
				if err != nil {
					return nil, err
				}
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"form - proto",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com?id=11&uint32=1234&bool=true", nil)
				if err != nil {
					return nil, err
				}
				return r, nil
			},
			&examplepb.Complex{
				Id:     11,
				Uint32: wrapperspb.UInt32(1234),
				Bool:   wrapperspb.Bool(true),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.genReq()
			if err != nil {
				t.Errorf("genReq() error = %v", err)
			}
			got := alloc(reflect.TypeOf(tt.want))
			if err = BindQuery(req, got.Interface()); (err != nil) != tt.wantErr {
				t.Errorf("BindQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, ok := tt.want.(proto.Message); ok {
				if diff := proto.Equal(got.Interface().(proto.Message), tt.want.(proto.Message)); !diff {
					t.Errorf("got = %v, want %v", got, tt.want)
				}
			} else {
				require.Equal(t, got.Interface(), tt.want)
			}
		})
	}
}

func TestBindUri(t *testing.T) {
	tests := []struct {
		name    string
		genReq  func() (*http.Request, error)
		want    any
		wantErr bool
	}{
		{
			"uri - no proto",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com/foo/bar", nil)
				if err != nil {
					return nil, err
				}

				param := url.Values{}
				param.Add("id", "foo")
				param.Add("name", "bar")
				return RequestWithUri(r, param), nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"uri - proto",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com?id=11&uint32=1234&bool=true", nil)
				if err != nil {
					return nil, err
				}
				param := url.Values{}
				param.Add("id", "11")
				param.Add("uint32", "1234")
				param.Add("bool", "true")
				return RequestWithUri(r, param), nil
			},
			&examplepb.Complex{
				Id:     11,
				Uint32: wrapperspb.UInt32(1234),
				Bool:   wrapperspb.Bool(true),
			},
			false,
		},
		{
			"uri - always in context",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com?id=11&uint32=1234&bool=true", nil)
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return RequestWithUri(r, nil), nil
			},
			&examplepb.Complex{},
			false,
		},
		{
			"uri - no existing in context",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com?id=11&uint32=1234&bool=true", nil)
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return r, nil
			},
			&examplepb.Complex{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.genReq()
			if err != nil {
				t.Errorf("genReq() error = %v", err)
			}
			got := alloc(reflect.TypeOf(tt.want))
			if err = BindUri(req, got.Interface()); (err != nil) != tt.wantErr {
				t.Errorf("BindQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, ok := tt.want.(proto.Message); ok {
				if diff := proto.Equal(got.Interface().(proto.Message), tt.want.(proto.Message)); !diff {
					t.Errorf("got = %v, want %v", got, tt.want)
				}
			} else {
				require.Equal(t, got.Interface(), tt.want)
			}
		})
	}
}

// helper
func alloc(t reflect.Type) reflect.Value {
	if t == nil {
		return reflect.ValueOf(new(interface{}))
	}
	return reflect.New(t.Elem())
}
