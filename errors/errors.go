//go:generate protoc --go_out=paths=source_relative:. errors.proto
package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	// ErrorsProtoPackageIsVersion3 this constant should not be referenced by any other code.
	ErrorsProtoPackageIsVersion3 = true
)

func (x *Error) Error() string {
	b, _ := json.Marshal(x)
	return *(*string)(unsafe.Pointer(&b))
}

// GRPCStatus returns the Status represented by se.
func (x *Error) GRPCStatus() *status.Status {
	s, _ := status.New(ToGRPCCode(int(x.Code)), x.Detail).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   x.Message,
			Metadata: x.Metadata,
		})
	return s
}

// Is matches each error in the chain with the target value.
func (x *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Message == x.Message
	}
	return false
}

func (x *Error) WithMetadata(md map[string]string) *Error {
	err := proto.Clone(x).(*Error)
	err.Metadata = md
	return err
}

func New(code int, message, detail string) *Error {
	return &Error{
		Code:    int32(code),
		Message: message,
		Detail:  detail,
	}
}

func Newf(code int, message, format string, a ...any) *Error {
	return New(code, message, fmt.Sprintf(format, a...))
}

func Errorf(code int, message, format string, a ...any) error {
	return New(code, message, fmt.Sprintf(format, a...))
}

// Code returns the http code for a error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return 200 //nolint:gomnd
	}
	return int(FromError(err).Code)
}

func Message(err error) string {
	if err == nil {
		return ""
	}
	return FromError(err).Message
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	if e := new(Error); errors.As(err, &e) {
		return e
	}

	if se, ok := status.FromError(err); ok {
		if se.Code() == codes.Unknown {
			return Parse(se.Message())
		}
		return New(FromGRPCCode(se.Code()), se.Message(), "")
	}

	return Parse(err.Error())
}

func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Code = 500
		e.Detail = err
	}
	if e.Code == 0 {
		e.Code = 500
	}
	return e
}
