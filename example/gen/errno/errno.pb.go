// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.27.0
// source: errno/errno.proto

package errno

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// #[errno]
type ErrorReason int32

const (
	// 服务器错误
	// #[errno(code=500)]
	ErrorReason_internal_server ErrorReason = 0
	// 操作超时
	ErrorReason_timeout ErrorReason = 1000
)

// Enum value maps for ErrorReason.
var (
	ErrorReason_name = map[int32]string{
		0:    "internal_server",
		1000: "timeout",
	}
	ErrorReason_value = map[string]int32{
		"internal_server": 0,
		"timeout":         1000,
	}
)

func (x ErrorReason) Enum() *ErrorReason {
	p := new(ErrorReason)
	*p = x
	return p
}

func (x ErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_errno_errno_proto_enumTypes[0].Descriptor()
}

func (ErrorReason) Type() protoreflect.EnumType {
	return &file_errno_errno_proto_enumTypes[0]
}

func (x ErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorReason.Descriptor instead.
func (ErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_errno_errno_proto_rawDescGZIP(), []int{0}
}

var File_errno_errno_proto protoreflect.FileDescriptor

var file_errno_errno_proto_rawDesc = []byte{
	0x0a, 0x11, 0x65, 0x72, 0x72, 0x6e, 0x6f, 0x2f, 0x65, 0x72, 0x72, 0x6e, 0x6f, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x65, 0x72, 0x72, 0x6e, 0x6f, 0x2a, 0x30, 0x0a, 0x0b, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x13, 0x0a, 0x0f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x10, 0x00, 0x12, 0x0c,
	0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x10, 0xe8, 0x07, 0x42, 0x64, 0x0a, 0x23,
	0x69, 0x6f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x2d, 0x67,
	0x6f, 0x2e, 0x64, 0x79, 0x6e, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x65, 0x72,
	0x72, 0x6e, 0x6f, 0x42, 0x0f, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x2d, 0x67, 0x6f, 0x2f, 0x64, 0x79, 0x6e,
	0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x65, 0x72, 0x72,
	0x6e, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_errno_errno_proto_rawDescOnce sync.Once
	file_errno_errno_proto_rawDescData = file_errno_errno_proto_rawDesc
)

func file_errno_errno_proto_rawDescGZIP() []byte {
	file_errno_errno_proto_rawDescOnce.Do(func() {
		file_errno_errno_proto_rawDescData = protoimpl.X.CompressGZIP(file_errno_errno_proto_rawDescData)
	})
	return file_errno_errno_proto_rawDescData
}

var file_errno_errno_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_errno_errno_proto_goTypes = []interface{}{
	(ErrorReason)(0), // 0: errno.ErrorReason
}
var file_errno_errno_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_errno_errno_proto_init() }
func file_errno_errno_proto_init() {
	if File_errno_errno_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_errno_errno_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_errno_errno_proto_goTypes,
		DependencyIndexes: file_errno_errno_proto_depIdxs,
		EnumInfos:         file_errno_errno_proto_enumTypes,
	}.Build()
	File_errno_errno_proto = out.File
	file_errno_errno_proto_rawDesc = nil
	file_errno_errno_proto_goTypes = nil
	file_errno_errno_proto_depIdxs = nil
}
