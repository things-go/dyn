package zfield

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func FromBinary(key string, vf func(context.Context) (v []byte, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Binary(key, v)
	}
}

func FromBoolp(key string, vf func(context.Context) (v *bool, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Boolp(key, v)
	}
}

func FromByteString(key string, vf func(context.Context) (v []byte, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.ByteString(key, v)
	}
}

func FromComplex128(key string, vf func(context.Context) (v complex128, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Complex128(key, v)
	}
}
func FromComplex128p(key string, vf func(context.Context) (v *complex128, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Complex128p(key, v)
	}
}
func FromComplex64(key string, vf func(context.Context) (v complex64, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Complex64(key, v)
	}
}
func FromComplex64p(key string, vf func(context.Context) (v *complex64, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Complex64p(key, v)
	}
}

func FromFloat64(key string, vf func(context.Context) (v float64, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Float64(key, v)
	}
}

func FromFloat64p(key string, vf func(context.Context) (v *float64, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Float64p(key, v)
	}
}
func FromFloat32(key string, vf func(context.Context) (v float32, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Float32(key, v)
	}
}
func FromFloat32p(key string, vf func(context.Context) (v *float32, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Float32p(key, v)
	}
}
func FromInt(key string, vf func(context.Context) (v int, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int(key, v)
	}
}
func FromIntp(key string, vf func(context.Context) (v *int, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Intp(key, v)
	}
}
func FromInt64(key string, vf func(context.Context) (v int64, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int64(key, v)
	}
}
func FromInt64p(key string, vf func(context.Context) (v *int64, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int64p(key, v)
	}
}
func FromInt32(key string, vf func(context.Context) (v int32, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int32(key, v)
	}
}
func FromInt32p(key string, vf func(context.Context) (v *int32, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int32p(key, v)
	}
}

func FromInt16(key string, vf func(context.Context) (v int16, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int16(key, v)
	}
}
func FromInt16p(key string, vf func(context.Context) (v *int16, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int16p(key, v)
	}
}
func FromInt8(key string, vf func(context.Context) (v int8, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int8(key, v)
	}
}
func FromInt8p(key string, vf func(context.Context) (v *int8, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int8p(key, v)
	}
}

func FromUint(key string, vf func(context.Context) (v int, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int(key, v)
	}
}
func FromUintp(key string, vf func(context.Context) (v *int, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Intp(key, v)
	}
}
func FromUint64(key string, vf func(context.Context) (v int64, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int64(key, v)
	}
}
func FromUint64p(key string, vf func(context.Context) (v *int64, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int64p(key, v)
	}
}
func FromUint32(key string, vf func(context.Context) (v int32, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int32(key, v)
	}
}
func FromUint32p(key string, vf func(context.Context) (v *int32, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int32p(key, v)
	}
}

func FromUint16(key string, vf func(context.Context) (v int16, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int16(key, v)
	}
}
func FromUint16p(key string, vf func(context.Context) (v *int16, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int16p(key, v)
	}
}
func FromUint8(key string, vf func(context.Context) (v int8, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int8(key, v)
	}
}
func FromUint8p(key string, vf func(context.Context) (v *int8, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Int8p(key, v)
	}
}

func FromString(key string, vf func(context.Context) (v string, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.String(key, v)
	}
}
func FromStringp(key string, vf func(context.Context) (v *string, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Stringp(key, v)
	}
}

func FromUintptr(key string, vf func(context.Context) (v uintptr, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Uintptr(key, v)
	}
}
func FromUintptrp(key string, vf func(context.Context) (v *uintptr, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Uintptrp(key, v)
	}
}
func FromReflect(key string, vf func(context.Context) (v interface{}, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Reflect(key, v)
	}
}
func FromStringer(key string, vf func(context.Context) (v fmt.Stringer, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Stringer(key, v)
	}
}
func FromTime(key string, vf func(context.Context) (v time.Time, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Time(key, v)
	}
}
func FromTimep(key string, vf func(context.Context) (v *time.Time, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Timep(key, v)
	}
}
func FromDuration(key string, vf func(context.Context) (v time.Duration, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Duration(key, v)
	}
}
func FromDurationp(key string, vf func(context.Context) (v *time.Duration, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Durationp(key, v)
	}
}
func FromAny(key string, vf func(context.Context) (v interface{}, skip bool)) func(context.Context) zap.Field {
	return func(ctx context.Context) zap.Field {
		v, skip := vf(ctx)
		if skip {
			return zap.Skip()
		}
		return zap.Any(key, v)
	}
}

func Binary(key string, vf func(context.Context) (v []byte, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Binary(key, v)
	}
}

func Boolp(key string, vf func(context.Context) (v *bool, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Boolp(key, v)
	}
}

func ByteString(key string, vf func(context.Context) (v []byte, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.ByteString(key, v)
	}
}

func Complex128(key string, vf func(context.Context) (v complex128, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Complex128(key, v)
	}
}
func Complex128p(key string, vf func(context.Context) (v *complex128, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Complex128p(key, v)
	}
}
func Complex64(key string, vf func(context.Context) (v complex64, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Complex64(key, v)
	}
}
func Complex64p(key string, vf func(context.Context) (v *complex64, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Complex64p(key, v)
	}
}

func Float64(key string, vf func(context.Context) (v float64, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Float64(key, v)
	}
}

func Float64p(key string, vf func(context.Context) (v *float64, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Float64p(key, v)
	}
}
func Float32(key string, vf func(context.Context) (v float32, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Float32(key, v)
	}
}
func Float32p(key string, vf func(context.Context) (v *float32, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Float32p(key, v)
	}
}
func Int(key string, vf func(context.Context) (v int, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int(key, v)
	}
}
func Intp(key string, vf func(context.Context) (v *int, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Intp(key, v)
	}
}
func Int64(key string, vf func(context.Context) (v int64, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int64(key, v)
	}
}
func Int64p(key string, vf func(context.Context) (v *int64, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int64p(key, v)
	}
}
func Int32(key string, vf func(context.Context) (v int32, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int32(key, v)
	}
}
func Int32p(key string, vf func(context.Context) (v *int32, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int32p(key, v)
	}
}

func Int16(key string, vf func(context.Context) (v int16, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int16(key, v)
	}
}
func Int16p(key string, vf func(context.Context) (v *int16, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int16p(key, v)
	}
}
func Int8(key string, vf func(context.Context) (v int8, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int8(key, v)
	}
}
func Int8p(key string, vf func(context.Context) (v *int8, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int8p(key, v)
	}
}

func Uint(key string, vf func(context.Context) (v int, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int(key, v)
	}
}
func Uintp(key string, vf func(context.Context) (v *int, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Intp(key, v)
	}
}
func Uint64(key string, vf func(context.Context) (v int64, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int64(key, v)
	}
}
func Uint64p(key string, vf func(context.Context) (v *int64, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int64p(key, v)
	}
}
func Uint32(key string, vf func(context.Context) (v int32, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int32(key, v)
	}
}
func Uint32p(key string, vf func(context.Context) (v *int32, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int32p(key, v)
	}
}

func Uint16(key string, vf func(context.Context) (v int16, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int16(key, v)
	}
}
func Uint16p(key string, vf func(context.Context) (v *int16, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int16p(key, v)
	}
}
func Uint8(key string, vf func(context.Context) (v int8, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int8(key, v)
	}
}
func Uint8p(key string, vf func(context.Context) (v *int8, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Int8p(key, v)
	}
}

func String(key string, vf func(context.Context) (v string, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.String(key, v)
	}
}
func Stringp(key string, vf func(context.Context) (v *string, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Stringp(key, v)
	}
}

func Uintptr(key string, vf func(context.Context) (v uintptr, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Uintptr(key, v)
	}
}
func Uintptrp(key string, vf func(context.Context) (v *uintptr, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Uintptrp(key, v)
	}
}
func Reflect(key string, vf func(context.Context) (v interface{}, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Reflect(key, v)
	}
}
func Stringer(key string, vf func(context.Context) (v fmt.Stringer, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Stringer(key, v)
	}
}
func Time(key string, vf func(context.Context) (v time.Time, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Time(key, v)
	}
}
func Timep(key string, vf func(context.Context) (v *time.Time, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Timep(key, v)
	}
}
func Duration(key string, vf func(context.Context) (v time.Duration, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Duration(key, v)
	}
}
func Durationp(key string, vf func(context.Context) (v *time.Duration, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Durationp(key, v)
	}
}
func Any(key string, vf func(context.Context) (v interface{}, skip bool)) func(*gin.Context) zap.Field {
	return func(c *gin.Context) zap.Field {
		v, skip := vf(c.Request.Context())
		if skip {
			return zap.Skip()
		}
		return zap.Any(key, v)
	}
}
