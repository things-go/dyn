package errorx

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

var _ error = (*Error)(nil)

// Error define the error type
type Error struct {
	code     int32
	message  string
	cause    error
	metadata map[string]string
}

// Error implement `Error() string` interface.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

// Code get the code
func (e *Error) Code() int32 {
	if e == nil {
		return 0
	}
	return e.code
}

// Message get the message
func (e *Error) Message() string {
	if e == nil {
		return ""
	}
	return e.message
}

// Metadata get metadata
func (e *Error) Metadata() map[string]string {
	if e == nil {
		return nil
	}
	return e.metadata
}

// Unwrap implement `Unwrap() error` interface.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

type Option func(*Error)

// New new Error
func New(code int32, message string, opts ...Option) *Error {
	e := &Error{code: code, message: message}
	return e.TakeOption(opts...)
}

// Newf new Error
func Newf(code int32, format string, args ...any) *Error {
	return &Error{code: code, message: fmt.Sprintf(format, args...)}
}

// TakeOption custom options
func (e *Error) TakeOption(opts ...Option) *Error {
	if e == nil {
		return nil
	}
	for _, f := range opts {
		f(e)
	}
	return e
}

// WithMessage modifies the message
func (e *Error) WithMessage(s string) *Error {
	return e.TakeOption(WithMessage(s))
}

// WithMessagef modifies the message
func (e *Error) WithMessagef(format string, args ...any) *Error {
	return e.TakeOption(WithMessagef(format, args...))
}

// WithCause set cause error
func (e *Error) WithCause(err error) *Error {
	return e.TakeOption(WithCause(err))
}

// WithError inner `errors.New`
func (e *Error) WithError(text string) *Error {
	return e.TakeOption(WithError(text))
}

// WithErrorf inner `fmt.Errorf`
func (e *Error) WithErrorf(format string, args ...any) *Error {
	return e.TakeOption(WithErrorf(format, args...))
}

// WithMetadata add metadata to the error
func (e *Error) WithMetadata(k, v string) *Error {
	return e.TakeOption(WithMetadata(k, v))
}

// WithMessage modifies the message
func WithMessage(s string) Option {
	return func(e *Error) {
		e.message = s
	}
}

// WithMessagef modifies the message
func WithMessagef(format string, args ...any) Option {
	return func(e *Error) {
		e.message = fmt.Sprintf(format, args...)
	}
}

// WithCause set cause error
func WithCause(err error) Option {
	return func(e *Error) {
		e.cause = err
	}
}

// WithError inner `errors.New`
func WithError(text string) Option {
	return WithCause(errors.New(text))
}

// WithErrorf inner `fmt.Errorf`
func WithErrorf(format string, args ...any) Option {
	return WithCause(fmt.Errorf(format, args...))
}

// WithMetadata add metadata to the error
func WithMetadata(k, v string) Option {
	return func(e *Error) {
		if k != "" && v != "" {
			if e.metadata == nil {
				e.metadata = make(map[string]string)
			}
			e.metadata[k] = v
		}
	}
}

// Parse parser error is `Error`, if not `Error`, new `Error` with code 500 and warp the err.
// err == nil: return nil
// err is not Error: return NewInternalServer
// err is Error(as te):
//
//	te == nil:  return nil
//	te != nil:  return te
func Parse(err error) *Error {
	if err == nil {
		return nil
	}
	if te := new(Error); errors.As(err, &te) {
		return te
	}
	return NewInternalServer(WithCause(err))
}

// EqualCode return true if error underlying code equal target code.
// err == nil: code = 200
// err is not Error: code = 500
// err is Error(as te):
//
//	te == nil:  code = 200
//	te != nil:  te.code
func EqualCode(err error, targetCode int32) bool {
	if err == nil {
		return http.StatusOK == targetCode
	}
	if te := new(Error); errors.As(err, &te) {
		if te == nil {
			return http.StatusOK == targetCode
		} else {
			return te != nil && te.code == targetCode
		}
	}
	return http.StatusInternalServerError == targetCode
}

// GRPCStatus returns the Status represented by se.
func (x *Error) GRPCStatus() *status.Status {
	s, _ := status.New(ToGRPCCode(int(x.code)), x.message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   x.Error(),
			Metadata: x.metadata,
		})
	return s
}
