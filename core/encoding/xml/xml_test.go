package xml

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var nilStruct *int

type Plain struct {
	V any
}
type PresenceTest struct {
	Exists *struct{}
}

var marshalTests = []struct {
	Value          any
	ExpectXML      string
	MarshalOnly    bool
	MarshalError   string
	UnmarshalOnly  bool
	UnmarshalError string
}{
	// Test nil marshals to nothing
	{Value: nil, ExpectXML: ``, MarshalOnly: true},
	{Value: nilStruct, ExpectXML: ``, MarshalOnly: true},

	// Test value types
	{Value: &Plain{true}, ExpectXML: `<Plain><V>true</V></Plain>`},
	{Value: &Plain{false}, ExpectXML: `<Plain><V>false</V></Plain>`},
	{Value: &Plain{int(42)}, ExpectXML: `<Plain><V>42</V></Plain>`},
	{Value: &Plain{int8(42)}, ExpectXML: `<Plain><V>42</V></Plain>`},
	{Value: &Plain{int16(42)}, ExpectXML: `<Plain><V>42</V></Plain>`},
	{Value: &Plain{int32(42)}, ExpectXML: `<Plain><V>42</V></Plain>`},
	{Value: &Plain{uint(42)}, ExpectXML: `<Plain><V>42</V></Plain>`},
	{Value: &Plain{uint8(42)}, ExpectXML: `<Plain><V>42</V></Plain>`},
	{Value: &Plain{uint16(42)}, ExpectXML: `<Plain><V>42</V></Plain>`},
	{Value: &Plain{uint32(42)}, ExpectXML: `<Plain><V>42</V></Plain>`},
	{Value: &Plain{float32(1.25)}, ExpectXML: `<Plain><V>1.25</V></Plain>`},
	{Value: &Plain{float64(1.25)}, ExpectXML: `<Plain><V>1.25</V></Plain>`},
	{Value: &Plain{uintptr(0xFFDD)}, ExpectXML: `<Plain><V>65501</V></Plain>`},
	{Value: &Plain{"gopher"}, ExpectXML: `<Plain><V>gopher</V></Plain>`},
	{Value: &Plain{[]byte("gopher")}, ExpectXML: `<Plain><V>gopher</V></Plain>`},
	{Value: &Plain{"</>"}, ExpectXML: `<Plain><V>&lt;/&gt;</V></Plain>`},
	{Value: &Plain{[]byte("</>")}, ExpectXML: `<Plain><V>&lt;/&gt;</V></Plain>`},
	{Value: &Plain{[3]byte{'<', '/', '>'}}, ExpectXML: `<Plain><V>&lt;/&gt;</V></Plain>`},
	{Value: &Plain{[]int{1, 2, 3}}, ExpectXML: `<Plain><V>1</V><V>2</V><V>3</V></Plain>`},
	{Value: &Plain{[3]int{1, 2, 3}}, ExpectXML: `<Plain><V>1</V><V>2</V><V>3</V></Plain>`},

	// Test time.
	{
		Value:     &Plain{time.Unix(1e9, 123456789).UTC()},
		ExpectXML: `<Plain><V>2001-09-09T01:46:40.123456789Z</V></Plain>`,
	},

	// A pointer to struct{} may be used to test for an element's presence.
	{
		Value:     &PresenceTest{new(struct{})},
		ExpectXML: `<PresenceTest><Exists></Exists></PresenceTest>`,
	},
	{
		Value:     &PresenceTest{},
		ExpectXML: `<PresenceTest></PresenceTest>`,
	},
}

func TestCodec(t *testing.T) {
	t.Run("Name", func(t *testing.T) {
		require.Equal(t, "xml", Name())
	})
	t.Run("Marshal", func(t *testing.T) {
		for _, tt := range marshalTests {
			data, err := Marshal(tt.Value)
			assert.NoError(t, err)
			assert.Equal(t, tt.ExpectXML, string(data))
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		for i, test := range marshalTests {
			if test.MarshalOnly {
				continue
			}
			if _, ok := test.Value.(*Plain); ok {
				continue
			}
			if test.ExpectXML == `<top>`+
				`<x><b xmlns="space">b</b>`+
				`<b xmlns="space1">b1</b></x>`+
				`</top>` {
				// TODO(rogpeppe): re-enable this test in
				// https://go-review.googlesource.com/#/c/5910/
				continue
			}

			vt := reflect.TypeOf(test.Value)
			dest := reflect.New(vt.Elem()).Interface()
			err := Unmarshal([]byte(test.ExpectXML), dest)

			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				if err != nil {
					if test.UnmarshalError == "" {
						t.Errorf("unmarshal(%#v): %s", test.ExpectXML, err)
						return
					}
					if !strings.Contains(err.Error(), test.UnmarshalError) {
						t.Errorf("unmarshal(%#v): %s, want %q", test.ExpectXML, err, test.UnmarshalError)
					}
					return
				}
				if got, want := dest, test.Value; !reflect.DeepEqual(got, want) {
					t.Errorf("unmarshal(%q):\nhave %#v\nwant %#v", test.ExpectXML, got, want)
				}
			})
		}
	})
}
