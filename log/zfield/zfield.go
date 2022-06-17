package zfield

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Binary(key string, vf func(context.Context) []byte) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Binary(key, vf(c.Request.Context()))
	}
}

func Boolp(key string, vf func(context.Context) *bool) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Boolp(key, vf(c.Request.Context()))
	}
}

func ByteString(key string, vf func(context.Context) []byte) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.ByteString(key, vf(c.Request.Context()))
	}
}

func Complex128(key string, vf func(context.Context) complex128) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Complex128(key, vf(c.Request.Context()))
	}
}
func Complex128p(key string, vf func(context.Context) *complex128) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {

		return zap.Complex128p(key, vf(c.Request.Context()))
	}
}
func Complex64(key string, vf func(context.Context) complex64) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Complex64(key, vf(c.Request.Context()))
	}
}
func Complex64p(key string, vf func(context.Context) *complex64) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Complex64p(key, vf(c.Request.Context()))
	}
}

func Float64(key string, vf func(context.Context) float64) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Float64(key, vf(c.Request.Context()))
	}
}

func Float64p(key string, vf func(context.Context) *float64) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Float64p(key, vf(c.Request.Context()))
	}
}
func Float32(key string, vf func(context.Context) float32) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Float32(key, vf(c.Request.Context()))
	}
}
func Float32p(key string, vf func(context.Context) *float32) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Float32p(key, vf(c.Request.Context()))
	}
}
func Int(key string, vf func(context.Context) int) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int(key, vf(c.Request.Context()))
	}
}
func Intp(key string, vf func(context.Context) *int) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Intp(key, vf(c.Request.Context()))
	}
}
func Int64(key string, vf func(context.Context) int64) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int64(key, vf(c.Request.Context()))
	}
}
func Int64p(key string, vf func(context.Context) *int64) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int64p(key, vf(c.Request.Context()))
	}
}
func Int32(key string, vf func(context.Context) int32) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int32(key, vf(c.Request.Context()))
	}
}
func Int32p(key string, vf func(context.Context) *int32) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int32p(key, vf(c.Request.Context()))
	}
}

func Int16(key string, vf func(context.Context) int16) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int16(key, vf(c.Request.Context()))
	}
}
func Int16p(key string, vf func(context.Context) *int16) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int16p(key, vf(c.Request.Context()))
	}
}
func Int8(key string, vf func(context.Context) int8) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int8(key, vf(c.Request.Context()))
	}
}
func Int8p(key string, vf func(context.Context) *int8) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int8p(key, vf(c.Request.Context()))
	}
}

func Uint(key string, vf func(context.Context) uint) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Uint(key, vf(c.Request.Context()))
	}
}
func Uintp(key string, vf func(context.Context) *uint) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Uintp(key, vf(c.Request.Context()))
	}
}
func Uint64(key string, vf func(context.Context) uint64) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Uint64(key, vf(c.Request.Context()))
	}
}
func Uint64p(key string, vf func(context.Context) *uint64) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Uint64p(key, vf(c.Request.Context()))
	}
}
func Uint32(key string, vf func(context.Context) uint32) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Uint32(key, vf(c.Request.Context()))
	}
}
func Uint32p(key string, vf func(context.Context) *uint32) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Uint32p(key, vf(c.Request.Context()))
	}
}

func Uint16(key string, vf func(context.Context) uint16) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Uint16(key, vf(c.Request.Context()))
	}
}
func Uint16p(key string, vf func(context.Context) *uint16) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Uint16p(key, vf(c.Request.Context()))
	}
}
func Uint8(key string, vf func(context.Context) uint8) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Uint8(key, vf(c.Request.Context()))
	}
}
func Uint8p(key string, vf func(context.Context) *uint8) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Uint8p(key, vf(c.Request.Context()))
	}
}

func String(key string, vf func(context.Context) string) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.String(key, vf(c.Request.Context()))
	}
}
func Stringp(key string, vf func(context.Context) *string) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Stringp(key, vf(c.Request.Context()))
	}
}

func Uintptr(key string, vf func(context.Context) uintptr) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Uintptr(key, vf(c.Request.Context()))
	}
}
func Uintptrp(key string, vf func(context.Context) *uintptr) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Uintptrp(key, vf(c.Request.Context()))
	}
}
func Reflect(key string, vf func(context.Context) interface{}) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Reflect(key, vf(c.Request.Context()))
	}
}
func Stringer(key string, vf func(context.Context) fmt.Stringer) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Stringer(key, vf(c.Request.Context()))
	}
}
func Time(key string, vf func(context.Context) time.Time) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Time(key, vf(c.Request.Context()))
	}
}
func Timep(key string, vf func(context.Context) *time.Time) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Timep(key, vf(c.Request.Context()))
	}
}
func Duration(key string, vf func(context.Context) time.Duration) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Duration(key, vf(c.Request.Context()))
	}
}
func Durationp(key string, vf func(context.Context) *time.Duration) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Durationp(key, vf(c.Request.Context()))
	}
}
func Any(key string, vf func(context.Context) interface{}) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Any(key, vf(c.Request.Context()))
	}
}

func ImmutBinary(key string, v []byte) func(*gin.Context) zap.Field {
	field := zap.Binary(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}

func ImmutBoolp(key string, v *bool) func(*gin.Context) zap.Field {
	field := zap.Boolp(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}

func ImmutByteString(key string, v []byte) func(*gin.Context) zap.Field {
	field := zap.ByteString(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}

func ImmutComplex128(key string, v complex128) func(*gin.Context) zap.Field {
	field := zap.Complex128(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutComplex128p(key string, v *complex128) func(*gin.Context) zap.Field {
	field := zap.Complex128p(key, v)
	return func(c *gin.Context) zap.Field {

		return field
	}
}
func ImmutComplex64(key string, v complex64) func(*gin.Context) zap.Field {
	field := zap.Complex64(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutComplex64p(key string, v *complex64) func(*gin.Context) zap.Field {
	field := zap.Complex64p(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}

func ImmutFloat64(key string, v float64) func(*gin.Context) zap.Field {
	field := zap.Float64(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}

func ImmutFloat64p(key string, v *float64) func(*gin.Context) zap.Field {
	field := zap.Float64p(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutFloat32(key string, v float32) func(*gin.Context) zap.Field {
	field := zap.Float32(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutFloat32p(key string, v *float32) func(*gin.Context) zap.Field {
	field := zap.Float32p(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutInt(key string, v int) func(*gin.Context) zap.Field {
	field := zap.Int(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutIntp(key string, v *int) func(*gin.Context) zap.Field {
	field := zap.Intp(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutInt64(key string, v int64) func(*gin.Context) zap.Field {
	field := zap.Int64(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutInt64p(key string, v *int64) func(*gin.Context) zap.Field {
	field := zap.Int64p(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutInt32(key string, v int32) func(*gin.Context) zap.Field {
	field := zap.Int32(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutInt32p(key string, v *int32) func(*gin.Context) zap.Field {
	field := zap.Int32p(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}

func ImmutInt16(key string, v int16) func(*gin.Context) zap.Field {
	field := zap.Int16(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutInt16p(key string, v *int16) func(*gin.Context) zap.Field {
	field := zap.Int16p(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutInt8(key string, v int8) func(*gin.Context) zap.Field {
	field := zap.Int8(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutInt8p(key string, v *int8) func(*gin.Context) zap.Field {
	field := zap.Int8p(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}

func ImmutUint(key string, v uint) func(*gin.Context) zap.Field {
	field := zap.Uint(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutUintp(key string, v *uint) func(*gin.Context) zap.Field {
	field := zap.Uintp(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutUint64(key string, v uint64) func(*gin.Context) zap.Field {
	field := zap.Uint64(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutUint64p(key string, v *uint64) func(*gin.Context) zap.Field {
	field := zap.Uint64p(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutUint32(key string, v uint32) func(*gin.Context) zap.Field {
	field := zap.Uint32(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutUint32p(key string, v *uint32) func(*gin.Context) zap.Field {
	field := zap.Uint32p(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}

func ImmutUint16(key string, v uint16) func(*gin.Context) zap.Field {
	field := zap.Uint16(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutUint16p(key string, v *uint16) func(*gin.Context) zap.Field {
	field := zap.Uint16p(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutUint8(key string, v uint8) func(*gin.Context) zap.Field {
	field := zap.Uint8(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutUint8p(key string, v *uint8) func(*gin.Context) zap.Field {
	field := zap.Uint8p(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}

func ImmutString(key string, v string) func(*gin.Context) zap.Field {
	field := zap.String(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutStringp(key string, v *string) func(*gin.Context) zap.Field {
	field := zap.Stringp(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}

func ImmutUintptr(key string, v uintptr) func(*gin.Context) zap.Field {
	field := zap.Uintptr(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutUintptrp(key string, v *uintptr) func(*gin.Context) zap.Field {
	field := zap.Uintptrp(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutReflect(key string, v interface{}) func(*gin.Context) zap.Field {
	field := zap.Reflect(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutStringer(key string, v fmt.Stringer) func(*gin.Context) zap.Field {
	field := zap.Stringer(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutTime(key string, v time.Time) func(*gin.Context) zap.Field {
	field := zap.Time(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutTimep(key string, v *time.Time) func(*gin.Context) zap.Field {
	field := zap.Timep(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutDuration(key string, v time.Duration) func(*gin.Context) zap.Field {
	field := zap.Duration(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutDurationp(key string, v *time.Duration) func(*gin.Context) zap.Field {
	field := zap.Durationp(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}
func ImmutAny(key string, v interface{}) func(*gin.Context) zap.Field {
	field := zap.Any(key, v)
	return func(c *gin.Context) zap.Field {
		return field
	}
}

/**************************************  help  ****************************************************/

func App(v string) func(c *gin.Context) zap.Field {
	return ImmutString("app", v)
}
func Component(v string) func(c *gin.Context) zap.Field {
	return ImmutString("component", v)
}
func Module(v string) func(c *gin.Context) zap.Field {
	return ImmutString("module", v)
}
func Unit(v string) func(c *gin.Context) zap.Field {
	return ImmutString("unit", v)
}
func Kind(v string) func(c *gin.Context) zap.Field {
	return ImmutString("kind", v)
}
func TraceId(f func(c context.Context) string) func(c *gin.Context) zap.Field {
	return String("traceId", f)
}
func RequestId(f func(c context.Context) string) func(c *gin.Context) zap.Field {
	return String("requestId", f)
}
