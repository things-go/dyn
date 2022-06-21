package json

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	testData "github.com/things-go/dyn/testdata/encoding"
)

type testEmbed struct {
	A int `json:"a"`
	B int `json:"b"`
	C int `json:"c"`
}

type testMessage struct {
	A string     `json:"a"`
	B string     `json:"b"`
	C string     `json:"c"`
	D *testEmbed `json:"d,omitempty"`
}

const (
	Unknown = iota
	Gopher
	Ruster
)

type mock struct {
	value int
}

func (a *mock) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		a.value = Unknown
	case "gopher":
		a.value = Gopher
	case "ruster":
		a.value = Ruster
	}

	return nil
}

func (a *mock) MarshalJSON() ([]byte, error) {
	var s string
	switch a.value {
	default:
		s = "unknown"
	case Gopher:
		s = "gopher"
	case Ruster:
		s = "ruster"
	}

	return json.Marshal(s)
}

func TestJSON(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		require.Equal(t, "json", New().Name())
	})
	t.Run("Marshal", func(t *testing.T) {
		tests := []struct {
			codec  Codec
			input  any
			expect string
		}{
			{
				codec:  New(),
				input:  &testMessage{},
				expect: `{"a":"","b":"","c":""}`,
			},
			{
				codec:  New(),
				input:  &testMessage{A: "a", B: "b", C: "c"},
				expect: `{"a":"a","b":"b","c":"c"}`,
			},
			{
				codec:  New(WithMarshalOptions(protojson.MarshalOptions{EmitUnpopulated: true}), WithUnmarshalOptions(protojson.UnmarshalOptions{DiscardUnknown: true})),
				input:  &testData.TestModel{Id: 1, Name: "golang", Hobby: []string{"1", "2"}, SnakeCase: map[string]string{"key": "value"}},
				expect: `{"id":"1","name":"golang","hobby":["1","2"],"snakeCase":{"key":"value"}}`,
			},
			{
				codec:  New(WithDisableProtoJSON()),
				input:  &testData.TestModel{Id: 1, Name: "golang", Hobby: []string{"1", "2"}, SnakeCase: map[string]string{"key": "value"}},
				expect: `{"id":1,"name":"golang","hobby":["1","2"],"snake_case":{"key":"value"}}`,
			},
			{
				codec:  New(),
				input:  &mock{value: Gopher},
				expect: `"gopher"`,
			},
		}
		for _, v := range tests {
			data, err := v.codec.Marshal(v.input)
			assert.NoError(t, err)
			got := strings.ReplaceAll(string(data), " ", "")
			assert.Equal(t, got, v.expect)
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		p := &testData.TestModel{}
		tests := []struct {
			codec  Codec
			input  string
			expect any
		}{
			{
				codec:  New(),
				input:  `{"a":"","b":"","c":"1111"}`,
				expect: &testMessage{},
			},
			{
				codec:  New(),
				input:  `{"a":"a","b":"b","c":"c"}`,
				expect: &testMessage{},
			},
			{
				codec:  New(WithMarshalOptions(protojson.MarshalOptions{EmitUnpopulated: true}), WithUnmarshalOptions(protojson.UnmarshalOptions{DiscardUnknown: true})),
				input:  `{"id":"1","name":"golang","hobby":["1","2"],"snakeCase":{}}`,
				expect: &testData.TestModel{},
			},
			{
				codec:  New(WithDisableProtoJSON()),
				input:  `{"id":1,"name":"golang","hobby":["1","2"]}`,
				expect: &testData.TestModel{},
			},
			{
				codec:  New(),
				input:  `{"id":1,"name":"golang","hobby":["1","2"]}`,
				expect: &p,
			},
			{
				codec:  New(),
				input:  `"ruster"`,
				expect: &mock{},
			},
		}

		for i, v := range tests {
			wantB := []byte(v.input)

			err := v.codec.Unmarshal(wantB, v.expect)
			assert.NoError(t, err, "Unmarshal: %d", i)

			gotB, err := v.codec.Marshal(v.expect)
			assert.NoError(t, err, "Marshal: %d", i)

			got := strings.ReplaceAll(string(gotB), " ", "")
			want := strings.ReplaceAll(string(wantB), " ", "")
			assert.Equal(t, want, got, "Marshal/Unmarshal: %d", i)
		}
	})
}

func TestJSON_default(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		require.Equal(t, "json", Name())
	})
	t.Run("Marshal", func(t *testing.T) {
		tests := []struct {
			input  any
			expect string
		}{
			{
				input:  &testMessage{},
				expect: `{"a":"","b":"","c":""}`,
			},
			{
				input:  &testMessage{A: "a", B: "b", C: "c"},
				expect: `{"a":"a","b":"b","c":"c"}`,
			},
			{
				input:  &testData.TestModel{Id: 1, Name: "golang", Hobby: []string{"1", "2"}, SnakeCase: map[string]string{"key": "value"}},
				expect: `{"id":"1","name":"golang","hobby":["1","2"],"snakeCase":{"key":"value"}}`,
			},
			{
				input:  &mock{value: Gopher},
				expect: `"gopher"`,
			},
		}
		for _, v := range tests {
			data, err := Marshal(v.input)
			assert.NoError(t, err)
			got := strings.ReplaceAll(string(data), " ", "")
			assert.Equal(t, got, v.expect)
		}
	})
	t.Run("Marshal", func(t *testing.T) {
		p := &testData.TestModel{}
		tests := []struct {
			input  string
			expect any
		}{
			{
				input:  `{"a":"","b":"","c":"1111"}`,
				expect: &testMessage{},
			},
			{
				input:  `{"a":"a","b":"b","c":"c"}`,
				expect: &testMessage{},
			},
			{
				input:  `{"id":"1","name":"golang","hobby":["1","2"],"snakeCase":{}}`,
				expect: &testData.TestModel{},
			},
			{
				input:  `{"id":1,"name":"golang","hobby":["1","2"]}`,
				expect: &p,
			},
			{
				input:  `"ruster"`,
				expect: &mock{},
			},
		}

		for _, v := range tests {
			wantB := []byte(v.input)

			err := Unmarshal(wantB, v.expect)
			assert.NoError(t, err, "Unmarshal")

			gotB, err := Marshal(v.expect)
			assert.NoError(t, err, "Marshal")

			got := strings.ReplaceAll(string(gotB), " ", "")
			want := strings.ReplaceAll(string(wantB), " ", "")
			assert.Equal(t, got, want)
		}
	})
}
