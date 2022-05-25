package json

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

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

func TestJSON_Marshal(t *testing.T) {
	tests := []struct {
		input  interface{}
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
			input:  &testData.TestModel{Id: 1, Name: "golang", Hobby: []string{"1", "2"}, Attrs: map[string]string{"key": "value"}},
			expect: `{"id":"1","name":"golang","hobby":["1","2"],"attrs":{"key":"value"}}`,
		},
		{
			input:  &mock{value: Gopher},
			expect: `"gopher"`,
		},
	}
	cdc := new(codec)
	for _, v := range tests {
		data, err := cdc.Marshal(v.input)
		assert.NoError(t, err)
		got := strings.ReplaceAll(string(data), " ", "")
		assert.Equal(t, got, v.expect)
	}
}

func TestJSON_Unmarshal(t *testing.T) {
	p := &testData.TestModel{}
	tests := []struct {
		input  string
		expect interface{}
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
			input:  `{"id":"1","name":"golang","hobby":["1","2"],"attrs":{}}`,
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
	cdc := new(codec)
	for _, v := range tests {
		wantB := []byte(v.input)

		err := cdc.Unmarshal(wantB, v.expect)
		assert.NoError(t, err, "Unmarshal")

		gotB, err := cdc.Marshal(v.expect)
		assert.NoError(t, err, "Marshal")

		got := strings.ReplaceAll(string(gotB), " ", "")
		want := strings.ReplaceAll(string(wantB), " ", "")
		assert.Equal(t, got, want)
	}
}
