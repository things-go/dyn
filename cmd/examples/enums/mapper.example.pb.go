// Code generated by protoc-gen-dyn-enum. DO NOT EDIT.
// versions:
//   - protoc-gen-dyn-enum v1.0.0
//   - protoc             v5.28.1
// source: nested.proto,non_nested.proto

package enums

func EnumComment() []map[int32]string {
	return []map[int32]string{
		{
			int32(Nested_Status_Unspecified): "\"unspecified\", aaaa",
			int32(Nested_Status_Up):          "nested1",
			int32(Nested_Status_Down):        "nested2",
			int32(Nested_Status_Left):        "nested3",
			int32(Nested_Status_Right):       "nested4",
			int32(Nested_Status_End):         "end",
		},
		{
			int32(Nested_Nested1_Type_Disable): "禁用",
			int32(Nested_Nested1_Type_Enable):  "启用",
		},
		{
			int32(NonNestedStatus_NonNestedStatusAnnote_Unspecified): "unspecified",
			int32(NonNestedStatus_NonNestedStatusAnnote_Up):          "up",
			int32(NonNestedStatus_NonNestedStatusAnnote_Down):        "down",
			int32(NonNestedStatus_NonNestedStatusAnnote_Left):        "left",
			int32(NonNestedStatus_NonNestedStatusAnnote_Right):       "right",
		},
	}
}
