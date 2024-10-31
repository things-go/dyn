// Code generated by errno-gen. DO NOT EDIT.
// version: v0.0.1
package errno

import (
	errors "github.com/things-go/dyn/errorx"
)

// ErrTimeout 1000: 操作超时
func ErrTimeout(opts ...errors.Option) error {
	return errors.New(int32(Timeout), Timeout.String(), opts...)
}

// ErrUserNotExist 1001: 用户不存在
func ErrUserNotExist(opts ...errors.Option) error {
	return errors.New(int32(UserNotExist), UserNotExist.String(), opts...)
}