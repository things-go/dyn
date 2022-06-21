package yaml

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestCodec(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		require.Equal(t, "yaml", Name())
	})
	t.Run("Marshal", func(t *testing.T) {
		var v = struct {
			A int
		}{
			A: 1,
		}
		got, err := Marshal(v)
		require.NoError(t, err)
		require.Equal(t, string(got), "a: 1\n")
	})
	t.Run("Unmarshal", func(t *testing.T) {
		for _, tt := range unmarshalTests {
			v := reflect.ValueOf(tt.value).Type()
			value := reflect.New(v)
			err := Unmarshal([]byte(tt.data), value.Interface())
			if _, ok := err.(*yaml.TypeError); !ok {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, value.Elem().Interface(), tt.value)
			}
		}
	})
}

// copy from yaml
type M map[string]any

type inlineB struct {
	B       int
	inlineC `yaml:",inline"`
}
type inlineC struct {
	C int
}
type inlineD struct {
	C *inlineC `yaml:",inline"`
	D int
}
type textUnmarshaler struct {
	S string
}

func (t *textUnmarshaler) UnmarshalText(s []byte) error {
	t.S = string(s)
	return nil
}

var unmarshalIntTest = 123
var unmarshalTests = []struct {
	data  string
	value any
}{
	{
		"",
		(*struct{})(nil),
	},
	{
		"{}", &struct{}{},
	},
	{
		"v: hi",
		map[string]string{"v": "hi"},
	},
	{
		"v: hi", map[string]any{"v": "hi"},
	},
	{
		"v: true",
		map[string]string{"v": "true"},
	},
	{
		"v: true",
		map[string]any{"v": true},
	},
	{
		"v: 10",
		map[string]any{"v": 10},
	},
	{
		"v: 0b10",
		map[string]any{"v": 2},
	},
	{
		"v: 0xA",
		map[string]any{"v": 10},
	},
	{
		"v: 4294967296",
		map[string]int64{"v": 4294967296},
	},
	{
		"v: 0.1",
		map[string]any{"v": 0.1},
	},
	{
		"v: .1",
		map[string]any{"v": 0.1},
	},
	{
		"v: .Inf",
		map[string]any{"v": math.Inf(+1)},
	},
	{
		"v: -.Inf",
		map[string]any{"v": math.Inf(-1)},
	},
	{
		"v: -10",
		map[string]any{"v": -10},
	},
	{
		"v: -.1",
		map[string]any{"v": -0.1},
	},

	// Simple values.
	{
		"123",
		&unmarshalIntTest,
	},

	// Floats from spec
	{
		"canonical: 6.8523e+5",
		map[string]any{"canonical": 6.8523e+5},
	},
	{
		"expo: 685.230_15e+03",
		map[string]any{"expo": 685.23015e+03},
	},
	{
		"fixed: 685_230.15",
		map[string]any{"fixed": 685230.15},
	},
	{
		"neginf: -.inf",
		map[string]any{"neginf": math.Inf(-1)},
	},
	{
		"fixed: 685_230.15",
		map[string]float64{"fixed": 685230.15},
	},
	// {"sexa: 190:20:30.15", map[string]any{"sexa": 0}}, // Unsupported
	// {"notanum: .NaN", map[string]any{"notanum": math.NaN()}}, // Equality of NaN fails.

	// Bools are per 1.2 spec.
	{
		"canonical: true",
		map[string]any{"canonical": true},
	},
	{
		"canonical: false",
		map[string]any{"canonical": false},
	},
	{
		"bool: True",
		map[string]any{"bool": true},
	},
	{
		"bool: False",
		map[string]any{"bool": false},
	},
	{
		"bool: TRUE",
		map[string]any{"bool": true},
	},
	{
		"bool: FALSE",
		map[string]any{"bool": false},
	},
	// For backwards compatibility with 1.1, decoding old strings into typed values still works.
	{
		"option: on",
		map[string]bool{"option": true},
	},
	{
		"option: y",
		map[string]bool{"option": true},
	},
	{
		"option: Off",
		map[string]bool{"option": false},
	},
	{
		"option: No",
		map[string]bool{"option": false},
	},
	{
		"option: other",
		map[string]bool{},
	},
	// Ints from spec
	{
		"canonical: 685230",
		map[string]any{"canonical": 685230},
	}, {
		"decimal: +685_230",
		map[string]any{"decimal": 685230},
	}, {
		"octal: 02472256",
		map[string]any{"octal": 685230},
	}, {
		"octal: -02472256",
		map[string]any{"octal": -685230},
	}, {
		"octal: 0o2472256",
		map[string]any{"octal": 685230},
	}, {
		"octal: -0o2472256",
		map[string]any{"octal": -685230},
	}, {
		"hexa: 0x_0A_74_AE",
		map[string]any{"hexa": 685230},
	}, {
		"bin: 0b1010_0111_0100_1010_1110",
		map[string]any{"bin": 685230},
	}, {
		"bin: -0b101010",
		map[string]any{"bin": -42},
	}, {
		"bin: -0b1000000000000000000000000000000000000000000000000000000000000000",
		map[string]any{"bin": -9223372036854775808},
	}, {
		"decimal: +685_230",
		map[string]int{"decimal": 685230},
	},

	// {"sexa: 190:20:30", map[string]any{"sexa": 0}}, // Unsupported

	// Nulls from spec
	{
		"empty:",
		map[string]any{"empty": nil},
	}, {
		"canonical: ~",
		map[string]any{"canonical": nil},
	}, {
		"english: null",
		map[string]any{"english": nil},
	}, {
		"~: null key",
		map[any]string{nil: "null key"},
	}, {
		"empty:",
		map[string]*bool{"empty": nil},
	},

	// Flow sequence
	{
		"seq: [A,B]",
		map[string]any{"seq": []any{"A", "B"}},
	}, {
		"seq: [A,B,C,]",
		map[string][]string{"seq": {"A", "B", "C"}},
	}, {
		"seq: [A,1,C]",
		map[string][]string{"seq": {"A", "1", "C"}},
	}, {
		"seq: [A,1,C]",
		map[string][]int{"seq": {1}},
	}, {
		"seq: [A,1,C]",
		map[string]any{"seq": []any{"A", 1, "C"}},
	},
	// Block sequence
	{
		"seq:\n - A\n - B",
		map[string]any{"seq": []any{"A", "B"}},
	}, {
		"seq:\n - A\n - B\n - C",
		map[string][]string{"seq": {"A", "B", "C"}},
	}, {
		"seq:\n - A\n - 1\n - C",
		map[string][]string{"seq": {"A", "1", "C"}},
	}, {
		"seq:\n - A\n - 1\n - C",
		map[string][]int{"seq": {1}},
	}, {
		"seq:\n - A\n - 1\n - C",
		map[string]any{"seq": []any{"A", 1, "C"}},
	},

	// Literal block scalar
	{
		"scalar: | # Comment\n\n literal\n\n \ttext\n\n",
		map[string]string{"scalar": "\nliteral\n\n\ttext\n"},
	},

	// Folded block scalar
	{
		"scalar: > # Comment\n\n folded\n line\n \n next\n line\n  * one\n  * two\n\n last\n line\n\n",
		map[string]string{"scalar": "\nfolded line\nnext line\n * one\n * two\n\nlast line\n"},
	},

	// Map inside interface with no type hints.
	{
		"a: {b: c}",
		map[any]any{"a": map[string]any{"b": "c"}},
	},
	// Non-string map inside interface with no type hints.
	{
		"a: {b: c, 1: d}",
		map[any]any{"a": map[any]any{"b": "c", 1: "d"}},
	},

	// Structs and type conversions.
	{
		"hello: world",
		&struct{ Hello string }{"world"},
	},
	{
		"a: {b: c}",
		&struct{ A struct{ B string } }{struct{ B string }{"c"}},
	},
	{
		"a: {b: c}",
		&struct{ A *struct{ B string } }{&struct{ B string }{"c"}},
	},
	{
		"a: {b: c}",
		&struct{ A map[string]string }{map[string]string{"b": "c"}},
	},
	{
		"a: {b: c}",
		&struct{ A *map[string]string }{&map[string]string{"b": "c"}},
	},
	{
		"a:",
		&struct{ A map[string]string }{},
	},
	{
		"a: 1",
		&struct{ A int }{1},
	},
	{
		"a: 1",
		&struct{ A float64 }{1},
	},
	{
		"a: 1.0",
		&struct{ A int }{1},
	},
	{
		"a: 1.0",
		&struct{ A uint }{1},
	},
	{
		"a: [1, 2]",
		&struct{ A []int }{[]int{1, 2}},
	},
	{
		"a: [1, 2]",
		&struct{ A [2]int }{[2]int{1, 2}},
	},
	{
		"a: 1",
		&struct{ B int }{0},
	},
	{
		"a: 1",
		&struct {
			B int "a" // nolint: govet
		}{1},
	},
	{
		// Some limited backwards compatibility with the 1.1 spec.
		"a: YES",
		&struct{ A bool }{true},
	},

	// Some cross type conversions
	{
		"v: 42",
		map[string]uint{"v": 42},
	},
	{
		"v: -42",
		map[string]uint{},
	},
	{
		"v: 4294967296",
		map[string]uint64{"v": 4294967296},
	},
	{
		"v: -4294967296",
		map[string]uint64{},
	},

	// int
	{
		"int_max: 2147483647",
		map[string]int{"int_max": math.MaxInt32},
	},
	{
		"int_min: -2147483648",
		map[string]int{"int_min": math.MinInt32},
	},
	{
		"int_overflow: 9223372036854775808", // math.MaxInt64 + 1
		map[string]int{},
	},

	// int64
	{
		"int64_max: 9223372036854775807",
		map[string]int64{"int64_max": math.MaxInt64},
	},
	{
		"int64_max_base2: 0b111111111111111111111111111111111111111111111111111111111111111",
		map[string]int64{"int64_max_base2": math.MaxInt64},
	},
	{
		"int64_min: -9223372036854775808",
		map[string]int64{"int64_min": math.MinInt64},
	},
	{
		"int64_neg_base2: -0b111111111111111111111111111111111111111111111111111111111111111",
		map[string]int64{"int64_neg_base2": -math.MaxInt64},
	},
	{
		"int64_overflow: 9223372036854775808", // math.MaxInt64 + 1
		map[string]int64{},
	},

	// uint
	{
		"uint_min: 0",
		map[string]uint{"uint_min": 0},
	},
	{
		"uint_max: 4294967295",
		map[string]uint{"uint_max": math.MaxUint32},
	},
	{
		"uint_underflow: -1",
		map[string]uint{},
	},

	// uint64
	{
		"uint64_min: 0",
		map[string]uint{"uint64_min": 0},
	},
	{
		"uint64_max: 18446744073709551615",
		map[string]uint64{"uint64_max": math.MaxUint64},
	},
	{
		"uint64_max_base2: 0b1111111111111111111111111111111111111111111111111111111111111111",
		map[string]uint64{"uint64_max_base2": math.MaxUint64},
	},
	{
		"uint64_maxint64: 9223372036854775807",
		map[string]uint64{"uint64_maxint64": math.MaxInt64},
	},
	{
		"uint64_underflow: -1",
		map[string]uint64{},
	},

	// float32
	{
		"float32_max: 3.40282346638528859811704183484516925440e+38",
		map[string]float32{"float32_max": math.MaxFloat32},
	},
	{
		"float32_nonzero: 1.401298464324817070923729583289916131280e-45",
		map[string]float32{"float32_nonzero": math.SmallestNonzeroFloat32},
	},
	{
		"float32_maxuint64: 18446744073709551615",
		map[string]float32{"float32_maxuint64": float32(math.MaxUint64)},
	},
	{
		"float32_maxuint64+1: 18446744073709551616",
		map[string]float32{"float32_maxuint64+1": float32(math.MaxUint64 + 1)},
	},

	// float64
	{
		"float64_max: 1.797693134862315708145274237317043567981e+308",
		map[string]float64{"float64_max": math.MaxFloat64},
	},
	{
		"float64_nonzero: 4.940656458412465441765687928682213723651e-324",
		map[string]float64{"float64_nonzero": math.SmallestNonzeroFloat64},
	},
	{
		"float64_maxuint64: 18446744073709551615",
		map[string]float64{"float64_maxuint64": float64(math.MaxUint64)},
	},
	{
		"float64_maxuint64+1: 18446744073709551616",
		map[string]float64{"float64_maxuint64+1": float64(math.MaxUint64 + 1)},
	},

	// Overflow cases.
	{
		"v: 4294967297",
		map[string]int32{},
	}, {
		"v: 128",
		map[string]int8{},
	},

	// Quoted values.
	{
		"'1': '\"2\"'",
		map[any]any{"1": "\"2\""},
	},
	{
		"v:\n- A\n- 'B\n\n  C'\n",
		map[string][]string{"v": {"A", "B\nC"}},
	},

	// Explicit tags.
	{
		"v: !!float '1.1'",
		map[string]any{"v": 1.1},
	},
	{
		"v: !!float 0",
		map[string]any{"v": float64(0)},
	},
	{
		"v: !!float -1",
		map[string]any{"v": float64(-1)},
	},
	{
		"v: !!null ''",
		map[string]any{"v": nil},
	},
	{
		"%TAG !y! tag:yaml.org,2002:\n---\nv: !y!int '1'",
		map[string]any{"v": 1},
	},

	// Non-specific tag (Issue #75)
	{
		"v: ! test",
		map[string]any{"v": "test"},
	},

	// Anchors and aliases.
	{
		"a: &x 1\nb: &y 2\nc: *x\nd: *y\n",
		&struct{ A, B, C, D int }{1, 2, 1, 2},
	},
	{
		"a: &a {c: 1}\nb: *a",
		&struct {
			A, B struct {
				C int
			}
		}{struct{ C int }{1}, struct{ C int }{1}},
	},
	{
		"a: &a [1, 2]\nb: *a",
		&struct{ B []int }{[]int{1, 2}},
	},

	// Bug #1133337
	{
		"foo: ''",
		map[string]*string{"foo": new(string)},
	},
	{
		"foo: null",
		map[string]*string{"foo": nil},
	},
	{
		"foo: null",
		map[string]string{"foo": ""},
	},
	{
		"foo: null",
		map[string]any{"foo": nil},
	},

	// Support for ~
	{
		"foo: ~",
		map[string]*string{"foo": nil},
	},
	{
		"foo: ~",
		map[string]string{"foo": ""},
	},
	{
		"foo: ~",
		map[string]any{"foo": nil},
	},

	// Ignored field
	{
		"a: 1\nb: 2\n",
		&struct {
			A int
			B int "-" // nolint: govet
		}{1, 0},
	},

	// Bug #1191981
	{
		"" +
			"%YAML 1.1\n" +
			"--- !!str\n" +
			`"Generic line break (no glyph)\n\` + "\n" +
			` Generic line break (glyphed)\n\` + "\n" +
			` Line separator\u2028\` + "\n" +
			` Paragraph separator\u2029"` + "\n",
		"" +
			"Generic line break (no glyph)\n" +
			"Generic line break (glyphed)\n" +
			"Line separator\u2028Paragraph separator\u2029",
	},

	// Struct inlining
	{
		"a: 1\nb: 2\nc: 3\n",
		&struct {
			A int
			C inlineB `yaml:",inline"`
		}{1, inlineB{2, inlineC{3}}},
	},

	// Struct inlining as a pointer.
	{
		"a: 1\nb: 2\nc: 3\n",
		&struct {
			A int
			C *inlineB `yaml:",inline"`
		}{1, &inlineB{2, inlineC{3}}},
	},
	{
		"a: 1\n",
		&struct {
			A int
			C *inlineB `yaml:",inline"`
		}{1, nil},
	},
	{
		"a: 1\nc: 3\nd: 4\n",
		&struct {
			A int
			C *inlineD `yaml:",inline"`
		}{1, &inlineD{&inlineC{3}, 4}},
	},

	// Map inlining
	{
		"a: 1\nb: 2\nc: 3\n",
		&struct {
			A int
			C map[string]int `yaml:",inline"`
		}{1, map[string]int{"b": 2, "c": 3}},
	},

	// bug 1243827
	{
		"a: -b_c",
		map[string]any{"a": "-b_c"},
	},
	{
		"a: +b_c",
		map[string]any{"a": "+b_c"},
	},
	{
		"a: 50cent_of_dollar",
		map[string]any{"a": "50cent_of_dollar"},
	},

	// issue #295 (allow scalars with colons in flow mappings and sequences)
	{
		"a: {b: https://github.com/go-yaml/yaml}",
		map[string]any{"a": map[string]any{
			"b": "https://github.com/go-yaml/yaml",
		}},
	},
	{
		"a: [https://github.com/go-yaml/yaml]",
		map[string]any{"a": []any{"https://github.com/go-yaml/yaml"}},
	},

	// Duration
	{
		"a: 3s",
		map[string]time.Duration{"a": 3 * time.Second},
	},

	// Issue #24.
	{
		"a: <foo>",
		map[string]string{"a": "<foo>"},
	},

	// Base 60 floats are obsolete and unsupported.
	{
		"a: 1:1\n",
		map[string]string{"a": "1:1"},
	},

	// Binary data.
	{
		"a: !!binary gIGC\n",
		map[string]string{"a": "\x80\x81\x82"},
	},
	{
		"a: !!binary |\n  " + strings.Repeat("kJCQ", 17) + "kJ\n  CQ\n",
		map[string]string{"a": strings.Repeat("\x90", 54)},
	},
	{
		"a: !!binary |\n  " + strings.Repeat("A", 70) + "\n  ==\n",
		map[string]string{"a": strings.Repeat("\x00", 52)},
	},

	// Issue #39.
	{
		"a:\n b:\n  c: d\n",
		map[string]struct{ B any }{"a": {map[string]any{"c": "d"}}},
	},

	// Custom map type.
	{
		"a: {b: c}",
		M{"a": M{"b": "c"}},
	},

	// Support encoding.TextUnmarshaler.
	{
		"a: 1.2.3.4\n",
		map[string]textUnmarshaler{"a": {S: "1.2.3.4"}},
	},
	{
		"a: 2015-02-24T18:19:39Z\n",
		map[string]textUnmarshaler{"a": {"2015-02-24T18:19:39Z"}},
	},

	// Timestamps
	{
		// Date only.
		"a: 2015-01-01\n",
		map[string]time.Time{"a": time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
	},
	{
		// RFC3339
		"a: 2015-02-24T18:19:39.12Z\n",
		map[string]time.Time{"a": time.Date(2015, 2, 24, 18, 19, 39, .12e9, time.UTC)},
	},
	{
		// RFC3339 with short dates.
		"a: 2015-2-3T3:4:5Z",
		map[string]time.Time{"a": time.Date(2015, 2, 3, 3, 4, 5, 0, time.UTC)},
	},
	{
		// ISO8601 lower case t
		"a: 2015-02-24t18:19:39Z\n",
		map[string]time.Time{"a": time.Date(2015, 2, 24, 18, 19, 39, 0, time.UTC)},
	},
	{
		// space separate, no time zone
		"a: 2015-02-24 18:19:39\n",
		map[string]time.Time{"a": time.Date(2015, 2, 24, 18, 19, 39, 0, time.UTC)},
	},
	// Some cases not currently handled. Uncomment these when
	// the code is fixed.
	//	{
	//		// space separated with time zone
	//		"a: 2001-12-14 21:59:43.10 -5",
	//		map[string]any{"a": time.Date(2001, 12, 14, 21, 59, 43, .1e9, time.UTC)},
	//	},
	//	{
	//		// arbitrary whitespace between fields
	//		"a: 2001-12-14 \t\t \t21:59:43.10 \t Z",
	//		map[string]any{"a": time.Date(2001, 12, 14, 21, 59, 43, .1e9, time.UTC)},
	//	},
	{
		// explicit string tag
		"a: !!str 2015-01-01",
		map[string]any{"a": "2015-01-01"},
	},
	{
		// explicit timestamp tag on quoted string
		"a: !!timestamp \"2015-01-01\"",
		map[string]time.Time{"a": time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
	},
	{
		// explicit timestamp tag on unquoted string
		"a: !!timestamp 2015-01-01",
		map[string]time.Time{"a": time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
	},
	{
		// quoted string that's a valid timestamp
		"a: \"2015-01-01\"",
		map[string]any{"a": "2015-01-01"},
	},
	{
		// explicit timestamp tag into interface.
		"a: !!timestamp \"2015-01-01\"",
		map[string]any{"a": time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
	},
	{
		// implicit timestamp tag into interface.
		"a: 2015-01-01",
		map[string]any{"a": time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
	},

	// Encode empty lists as zero-length slices.
	{
		"a: []",
		&struct{ A []int }{[]int{}},
	},

	// UTF-16-LE
	{
		"\xff\xfe\xf1\x00o\x00\xf1\x00o\x00:\x00 \x00v\x00e\x00r\x00y\x00 \x00y\x00e\x00s\x00\n\x00",
		M{"ñoño": "very yes"},
	},
	// UTF-16-LE with surrogate.
	{
		"\xff\xfe\xf1\x00o\x00\xf1\x00o\x00:\x00 \x00v\x00e\x00r\x00y\x00 \x00y\x00e\x00s\x00 \x00=\xd8\xd4\xdf\n\x00",
		M{"ñoño": "very yes 🟔"},
	},

	// UTF-16-BE
	{
		"\xfe\xff\x00\xf1\x00o\x00\xf1\x00o\x00:\x00 \x00v\x00e\x00r\x00y\x00 \x00y\x00e\x00s\x00\n",
		M{"ñoño": "very yes"},
	},
	// UTF-16-BE with surrogate.
	{
		"\xfe\xff\x00\xf1\x00o\x00\xf1\x00o\x00:\x00 \x00v\x00e\x00r\x00y\x00 \x00y\x00e\x00s\x00 \xd8=\xdf\xd4\x00\n",
		M{"ñoño": "very yes 🟔"},
	},

	// This *is* in fact a float number, per the spec. #171 was a mistake.
	{
		"a: 123456e1\n",
		M{"a": 123456e1},
	},
	{
		"a: 123456E1\n",
		M{"a": 123456e1},
	},
	// yaml-test-suite 3GZX: Spec Example 7.1. Alias Nodes
	{
		"First occurrence: &anchor Foo\nSecond occurrence: *anchor\nOverride anchor: &anchor Bar\nReuse anchor: *anchor\n",
		map[string]any{
			"First occurrence":  "Foo",
			"Second occurrence": "Foo",
			"Override anchor":   "Bar",
			"Reuse anchor":      "Bar",
		},
	},
	// Single document with garbage following it.
	{
		"---\nhello\n...\n}not yaml",
		"hello",
	},

	// Comment scan exhausting the input buffer (issue #469).
	{
		"true\n#" + strings.Repeat(" ", 512*3),
		"true",
	},
	{
		"true #" + strings.Repeat(" ", 512*3),
		"true",
	},

	// CRLF
	{
		"a: b\r\nc:\r\n- d\r\n- e\r\n",
		map[string]any{
			"a": "b",
			"c": []any{"d", "e"},
		},
	},
}
