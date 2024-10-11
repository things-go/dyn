// Code generated by protoc-gen-saber-enum. DO NOT EDIT.
// versions:
//   - protoc-gen-saber-enum v0.5.0
//   - protoc                v5.28.1
// source: nested.proto

package enums

import (
	"strconv"
)

// Enum value mapping for Nested_Status.
var (
	__Nested_StatusMapping_Desc = map[Nested_Status]string{
		0:   "\"unspecified\", aaaa",
		1:   "nested1",
		2:   "nested2",
		3:   "nested3",
		4:   "nested4",
		999: "end",
	}
	__Nested_StatusMapping_Value = map[string]Nested_Status{
		"\"unspecified\", aaaa": 0,
		"nested1":               1,
		"nested2":               2,
		"nested3":               3,
		"nested4":               4,
		"end":                   999,
	}
)

// IntoNumber returns the enum value as an integer.
func (x Nested_Status) IntoNumber() int32 {
	return int32(x)
}

// IntoNumberString returns the enum value as an integer string.
func (x Nested_Status) IntoNumberString() string {
	return strconv.FormatInt(int64(x), 10)
}

// MappingDescriptor mapping description.
// Status 状态值
// [0:\"unspecified\", aaaa,1:nested1,2:nested2,3:nested3,4:nested4,999:end]
func (x Nested_Status) MappingDescriptor() string {
	return __Nested_StatusMapping_Desc[x]
}

// EnumCount the number of enum values.
func (Nested_Status) EnumCount() int {
	return 6
}

// GetNested_StatusValue get mapping value
// Status 状态值
// [0:\"unspecified\", aaaa,1:nested1,2:nested2,3:nested3,4:nested4,999:end]
func GetNested_StatusValue(s string) int {
	return int(__Nested_StatusMapping_Value[s])
}

// Enum value mapping for Nested_Nested1_Type.
var (
	__Nested_Nested1_TypeMapping_Desc = map[Nested_Nested1_Type]string{
		0: "禁用",
		1: "启用",
	}
	__Nested_Nested1_TypeMapping_Value = map[string]Nested_Nested1_Type{
		"禁用": 0,
		"启用": 1,
	}
)

// IntoNumber returns the enum value as an integer.
func (x Nested_Nested1_Type) IntoNumber() int32 {
	return int32(x)
}

// IntoNumberString returns the enum value as an integer string.
func (x Nested_Nested1_Type) IntoNumberString() string {
	return strconv.FormatInt(int64(x), 10)
}

// MappingDescriptor mapping description.
// Type 类型
// [0:禁用,1:启用]
func (x Nested_Nested1_Type) MappingDescriptor() string {
	return __Nested_Nested1_TypeMapping_Desc[x]
}

// EnumCount the number of enum values.
func (Nested_Nested1_Type) EnumCount() int {
	return 2
}

// GetNested_Nested1_TypeValue get mapping value
// Type 类型
// [0:禁用,1:启用]
func GetNested_Nested1_TypeValue(s string) int {
	return int(__Nested_Nested1_TypeMapping_Value[s])
}
