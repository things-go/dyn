// Code generated by protoc-gen-dyn-errno. DO NOT EDIT.
// versions:
//   - protoc-gen-dyn-errno v1.0.0
//   - protoc                 v5.27.0
// source: errno/errno.proto

package errno

import (
	errors "github.com/things-go/dyn/errorx"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = errors.New

// ErrInternalServer 500: 服务器错误
func ErrInternalServer(opts ...errors.Option) *errors.Error {
	return errors.New(500, "服务器错误", opts...)
}

// ErrTimeout 1000: 操作超时
func ErrTimeout(opts ...errors.Option) *errors.Error {
	return errors.New(1000, "操作超时", opts...)
}
