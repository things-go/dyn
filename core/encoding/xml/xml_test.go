package xml

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var nilStruct *int

type Plain struct {
	V interface{}
}

var marshalTests = []struct {
	Value          interface{}
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
}

func TestCodec_Marshal(t *testing.T) {
	for _, tt := range marshalTests {
		data, err := (codec{}).Marshal(tt.Value)
		if err != nil {
			t.Errorf("marshal(%#v): %s", tt.Value, err)
		}
		if got, want := string(data), tt.ExpectXML; got != want {
			if strings.Contains(want, "\n") {
				t.Errorf("marshal(%#v):\nHAVE:\n%s\nWANT:\n%s", tt.Value, got, want)
			} else {
				t.Errorf("marshal(%#v):\nhave %#q\nwant %#q", tt.Value, got, want)
			}
		}
	}
}

func TestCodec_Unmarshal(t *testing.T) {
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
		err := (codec{}).Unmarshal([]byte(test.ExpectXML), dest)

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
}
