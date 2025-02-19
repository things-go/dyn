// Code generated by protoc-gen-dyn-enum. DO NOT EDIT.
// versions:
//   - protoc-gen-dyn-enum v1.0.0
//   - protoc                v5.28.1
// source: nested.proto

package enums

import (
	"strconv"
)

// Enum value label/mapping for Nested_Status.
var (
	__Nested_Status_xx_Label = map[Nested_Status]string{
		0:   "custom label",
		1:   "nested1",
		2:   "nested2",
		3:   "nested3",
		4:   "nested4",
		999: "end",
	}
	__Nested_Status_xx_Value = map[string]Nested_Status{
		"custom label": 0,
		"nested1":      1,
		"nested2":      2,
		"nested3":      3,
		"nested4":      4,
		"end":          999,
	}
)

// EnumCount the number of enum value.
func (Nested_Status) EnumCount() int {
	return 6
}

// IntoNumber returns the enum value as an integer.
func (x Nested_Status) IntoNumber() int32 {
	return int32(x)
}

// IntoNumberString returns the enum value as an integer string.
func (x Nested_Status) IntoNumberString() string {
	return strconv.FormatInt(int64(x), 10)
}

// EnumLabel the label of enum value.
// Status 状态值
// [0:custom label,1:nested1,2:nested2,3:nested3,4:nested4,999:end]
func (x Nested_Status) EnumLabel() string {
	return __Nested_Status_xx_Label[x]
}

// ParseEnumLabel parse the label.
// Status 状态值
// [0:custom label,1:nested1,2:nested2,3:nested3,4:nested4,999:end]
func (e *Nested_Status) ParseEnumLabel(s string) {
	*e = __Nested_Status_xx_Value[s]
}

// Enum value label/mapping for Nested_Nested1_Type.
var (
	__Nested_Nested1_Type_xx_Label = map[Nested_Nested1_Type]string{
		0: "禁用",
		1: "启用",
	}
	__Nested_Nested1_Type_xx_Value = map[string]Nested_Nested1_Type{
		"禁用": 0,
		"启用": 1,
	}
)

// EnumCount the number of enum value.
func (Nested_Nested1_Type) EnumCount() int {
	return 2
}

// IntoNumber returns the enum value as an integer.
func (x Nested_Nested1_Type) IntoNumber() int32 {
	return int32(x)
}

// IntoNumberString returns the enum value as an integer string.
func (x Nested_Nested1_Type) IntoNumberString() string {
	return strconv.FormatInt(int64(x), 10)
}

// EnumLabel the label of enum value.
// Type 类型
// [0:禁用,1:启用]
func (x Nested_Nested1_Type) EnumLabel() string {
	return __Nested_Nested1_Type_xx_Label[x]
}

// ParseEnumLabel parse the label.
// Type 类型
// [0:禁用,1:启用]
func (e *Nested_Nested1_Type) ParseEnumLabel(s string) {
	*e = __Nested_Nested1_Type_xx_Value[s]
}
