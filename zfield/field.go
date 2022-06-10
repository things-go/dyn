package zfield

import (
	"context"
	"fmt"
	"time"

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

/**************************************  help  ****************************************************/

func FromApp(v string) func(ctx context.Context) zap.Field {
	field := zap.String("app", v)
	return func(ctx context.Context) zap.Field { return field }
}
func FromComponent(v string) func(ctx context.Context) zap.Field {
	field := zap.String("component", v)
	return func(ctx context.Context) zap.Field { return field }
}
func FromModule(v string) func(ctx context.Context) zap.Field {
	field := zap.String("module", v)
	return func(ctx context.Context) zap.Field { return field }
}
func FromUnit(v string) func(ctx context.Context) zap.Field {
	field := zap.String("unit", v)
	return func(ctx context.Context) zap.Field { return field }
}
func FromKind(v string) func(ctx context.Context) zap.Field {
	field := zap.String("kind", v)
	return func(ctx context.Context) zap.Field { return field }
}
func FromTraceId(f func(c context.Context) string) func(ctx context.Context) zap.Field {
	return FromString("traceId", f)
}
func FromRequestId(f func(c context.Context) string) func(ctx context.Context) zap.Field {
	return FromString("requestId", f)
}
