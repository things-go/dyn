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

/**************************************  help  ****************************************************/

func App(v string) func(c *gin.Context) zap.Field {
	field := zap.String("app", v)
	return func(c *gin.Context) zap.Field { return field }
}
func Component(v string) func(c *gin.Context) zap.Field {
	field := zap.String("component", v)
	return func(c *gin.Context) zap.Field { return field }
}
func Module(v string) func(c *gin.Context) zap.Field {
	field := zap.String("module", v)
	return func(c *gin.Context) zap.Field { return field }
}
func Unit(v string) func(c *gin.Context) zap.Field {
	field := zap.String("unit", v)
	return func(c *gin.Context) zap.Field { return field }
}
func Kind(v string) func(c *gin.Context) zap.Field {
	field := zap.String("kind", v)
	return func(c *gin.Context) zap.Field { return field }
}
func TraceId(f func(c context.Context) string) func(c *gin.Context) zap.Field {
	return String("traceId", f)
}
func RequestId(f func(c context.Context) string) func(c *gin.Context) zap.Field {
	return String("requestId", f)
}
