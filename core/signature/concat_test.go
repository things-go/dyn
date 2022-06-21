package signature

import (
	"encoding/json"
	"errors"
	"html/template"
	"testing"
)

func Test_ConcatMapWithSort(t *testing.T) {
	tests := []struct {
		name string
		mp   map[string]string
		sep1 string
		sep2 string
		want string
	}{
		{
			"empty",
			map[string]string{},
			"=",
			"&",
			"",
		},
		{
			"",
			map[string]string{
				"b": "1",
				"d": "a",
				"a": "10",
			},
			"=",
			"&",
			"a=10&b=1&d=a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConcatMapWithSort(tt.mp, tt.sep1, tt.sep2); got != tt.want {
				t.Errorf("ConcatMapWithSort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ConcatMap(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		hasBrace bool
		want     string
	}{
		{
			"",
			nil,
			false,
			"",
		},
		{
			"",
			map[string]any{},
			true,
			"",
		},
		{
			"",
			map[string]any{
				"b": 1,
				"d": "a",
				"a": 10,
				"e": map[string]any{
					"b": 1,
					"d": "a",
					"a": 10,
				},
				"f": []any{1, 6, 0},
			},
			false,
			"a=10&b=1&d=a&e={a=10&b=1&d=a}&f=[1,6,0]",
		},
		{
			"",
			map[string]any{
				"b": 1,
				"d": "a",
				"a": 10,
			},
			true,
			"{a=10&b=1&d=a}",
		},
		{
			"",
			map[string]any{
				"b": 1,
				"d": "a",
				"a": 10,
				"e": []map[string]any{
					{
						"b": 1,
						"d": "a",
						"a": 10,
					},
					{
						"b": 2,
						"d": "b",
						"a": 11,
					},
				},
				"f": []any{1, 6, 0},
			},
			false,
			"a=10&b=1&d=a&e=[{a=10&b=1&d=a},{a=11&b=2&d=b}]&f=[1,6,0]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConcatMap(tt.m, tt.hasBrace); got != tt.want {
				t.Errorf("ConcatMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ConcatArray(t *testing.T) {
	tests := []struct {
		name string
		arr  any
		want string
	}{
		{
			"",
			nil,
			"",
		},
		{
			"",
			[]any{},
			"",
		},
		{
			"",
			[]any{"d", "a", "f", "1"},
			"[d,a,f,1]",
		},
		{
			"",
			[]any{5, 1, 0, 2},
			"[5,1,0,2]",
		},
		{
			"",
			[]string{"5", "1", "0", "2"},
			"[5,1,0,2]",
		},
		{
			"",
			[]int{5, 1, 0, 2},
			"[5,1,0,2]",
		},
		{
			"",
			[]float64{5, 1, 0, 2},
			"[5,1,0,2]",
		},
		{
			"",
			[]int64{5, 1, 0, 2},
			"[5,1,0,2]",
		},
		{
			"",
			"aaa",
			"aaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConcatArray(tt.arr); got != tt.want {
				t.Errorf("ConcatArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toString(t *testing.T) {
	type Key struct {
		k string
	}
	key := &Key{"foo"}

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{"", int(8), "8"},
		{"", int8(8), "8"},
		{"", int16(8), "8"},
		{"", int32(8), "8"},
		{"", int64(8), "8"},
		{"", uint(8), "8"},
		{"", uint8(8), "8"},
		{"", uint16(8), "8"},
		{"", uint32(8), "8"},
		{"", uint64(8), "8"},
		{"", float32(8.31), "8.31"},
		{"", float64(8.31), "8.31"},
		{"", true, "true"},
		{"", false, "false"},
		{"", nil, ""},
		{"", []byte("one time"), "one time"},
		{"", "one more time", "one more time"},
		{"", template.HTML("one time"), "one time"},
		{"", template.URL("http://somehost.foo"), "http://somehost.foo"},
		{"", template.JS("(1+2)"), "(1+2)"},
		{"", template.CSS("a"), "a"},
		{"", template.HTMLAttr("a"), "a"},
		// errors
		{"", testing.T{}, ""},
		{"", key, ""},
		{"", fmtStringer{}, "fmtStringer"},
		{"", json.Number("100"), "100"},
		{"", errors.New("100"), "100"},
		{"", []any{1, 8, 3}, "[1,8,3]"},
		{"", map[string]any{"d": 1, "a": 8, "c": 3}, "{a=8&c=3&d=1}"},
		{"", []int{1, 8, 3}, "[1,8,3]"},
		{"", []int{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toString(tt.input); got != tt.want {
				t.Errorf("ConcatArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

type fmtStringer struct{}

func (fmtStringer) String() string {
	return "fmtStringer"
}

func BenchmarkConcatSortMap(b *testing.B) {
	mp := map[string]string{
		"b": "1",
		"d": "a",
		"a": "10",
	}
	for i := 0; i < b.N; i++ {
		ConcatMapWithSort(mp, "=", "&")
	}
}

func BenchmarkConcatMap(b *testing.B) {
	mp := map[string]any{
		"b": "1",
		"d": "a",
		"f": []any{1, 6, 0},
	}
	for i := 0; i < b.N; i++ {
		ConcatMap(mp, false)
	}
}

func BenchmarkConcatArray(b *testing.B) {
	arr := []any{"d", "a", "f", "1"}
	for i := 0; i < b.N; i++ {
		ConcatArray(arr)
	}
}
