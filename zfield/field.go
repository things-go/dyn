package zfield

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func FromBinary(key string, vf func(context.Context) []byte) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Binary(key, vf(ctx))
	}
}

func FromBoolp(key string, vf func(context.Context) *bool) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Boolp(key, vf(ctx))
	}
}

func FromByteString(key string, vf func(context.Context) []byte) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.ByteString(key, vf(ctx))
	}
}

func FromComplex128(key string, vf func(context.Context) complex128) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Complex128(key, vf(ctx))
	}
}
func FromComplex128p(key string, vf func(context.Context) *complex128) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Complex128p(key, vf(ctx))
	}
}
func FromComplex64(key string, vf func(context.Context) complex64) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Complex64(key, vf(ctx))
	}
}
func FromComplex64p(key string, vf func(context.Context) *complex64) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Complex64p(key, vf(ctx))
	}
}

func FromFloat64(key string, vf func(context.Context) float64) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Float64(key, vf(ctx))
	}
}

func FromFloat64p(key string, vf func(context.Context) *float64) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Float64p(key, vf(ctx))
	}
}
func FromFloat32(key string, vf func(context.Context) float32) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Float32(key, vf(ctx))
	}
}
func FromFloat32p(key string, vf func(context.Context) *float32) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Float32p(key, vf(ctx))
	}
}
func FromInt(key string, vf func(context.Context) int) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int(key, vf(ctx))
	}
}
func FromIntp(key string, vf func(context.Context) *int) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Intp(key, vf(ctx))
	}
}
func FromInt64(key string, vf func(context.Context) int64) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int64(key, vf(ctx))
	}
}
func FromInt64p(key string, vf func(context.Context) *int64) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int64p(key, vf(ctx))
	}
}
func FromInt32(key string, vf func(context.Context) int32) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int32(key, vf(ctx))
	}
}
func FromInt32p(key string, vf func(context.Context) *int32) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int32p(key, vf(ctx))
	}
}

func FromInt16(key string, vf func(context.Context) int16) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int16(key, vf(ctx))
	}
}
func FromInt16p(key string, vf func(context.Context) *int16) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int16p(key, vf(ctx))
	}
}
func FromInt8(key string, vf func(context.Context) int8) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int8(key, vf(ctx))
	}
}
func FromInt8p(key string, vf func(context.Context) *int8) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int8p(key, vf(ctx))
	}
}

func FromUint(key string, vf func(context.Context) int) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int(key, vf(ctx))
	}
}
func FromUintp(key string, vf func(context.Context) *int) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Intp(key, vf(ctx))
	}
}
func FromUint64(key string, vf func(context.Context) int64) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int64(key, vf(ctx))
	}
}
func FromUint64p(key string, vf func(context.Context) *int64) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int64p(key, vf(ctx))
	}
}
func FromUint32(key string, vf func(context.Context) int32) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int32(key, vf(ctx))
	}
}
func FromUint32p(key string, vf func(context.Context) *int32) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int32p(key, vf(ctx))
	}
}

func FromUint16(key string, vf func(context.Context) int16) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int16(key, vf(ctx))
	}
}
func FromUint16p(key string, vf func(context.Context) *int16) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int16p(key, vf(ctx))
	}
}
func FromUint8(key string, vf func(context.Context) int8) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int8(key, vf(ctx))
	}
}
func FromUint8p(key string, vf func(context.Context) *int8) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Int8p(key, vf(ctx))
	}
}

func FromString(key string, vf func(context.Context) string) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.String(key, vf(ctx))
	}
}
func FromStringp(key string, vf func(context.Context) *string) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Stringp(key, vf(ctx))
	}
}

func FromUintptr(key string, vf func(context.Context) uintptr) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Uintptr(key, vf(ctx))
	}
}
func FromUintptrp(key string, vf func(context.Context) *uintptr) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Uintptrp(key, vf(ctx))
	}
}
func FromReflect(key string, vf func(context.Context) interface{}) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Reflect(key, vf(ctx))
	}
}
func FromStringer(key string, vf func(context.Context) fmt.Stringer) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Stringer(key, vf(ctx))
	}
}
func FromTime(key string, vf func(context.Context) time.Time) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Time(key, vf(ctx))
	}
}
func FromTimep(key string, vf func(context.Context) *time.Time) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Timep(key, vf(ctx))
	}
}
func FromDuration(key string, vf func(context.Context) time.Duration) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Duration(key, vf(ctx))
	}
}
func FromDurationp(key string, vf func(context.Context) *time.Duration) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Durationp(key, vf(ctx))
	}
}
func FromAny(key string, vf func(context.Context) interface{}) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		return zap.Any(key, vf(ctx))
	}
}

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

func Uint(key string, vf func(context.Context) int) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int(key, vf(c.Request.Context()))
	}
}
func Uintp(key string, vf func(context.Context) *int) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Intp(key, vf(c.Request.Context()))
	}
}
func Uint64(key string, vf func(context.Context) int64) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int64(key, vf(c.Request.Context()))
	}
}
func Uint64p(key string, vf func(context.Context) *int64) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int64p(key, vf(c.Request.Context()))
	}
}
func Uint32(key string, vf func(context.Context) int32) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int32(key, vf(c.Request.Context()))
	}
}
func Uint32p(key string, vf func(context.Context) *int32) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int32p(key, vf(c.Request.Context()))
	}
}

func Uint16(key string, vf func(context.Context) int16) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int16(key, vf(c.Request.Context()))
	}
}
func Uint16p(key string, vf func(context.Context) *int16) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int16p(key, vf(c.Request.Context()))
	}
}
func Uint8(key string, vf func(context.Context) int8) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int8(key, vf(c.Request.Context()))
	}
}
func Uint8p(key string, vf func(context.Context) *int8) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		return zap.Int8p(key, vf(c.Request.Context()))
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
