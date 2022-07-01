package log

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	l, lv := New(WithConfig(Config{Level: "debug", Format: "json"}))
	ReplaceGlobals(NewLoggerWith(l, lv))
	SetDefaultValuer(Caller(2), func(ctx context.Context) Field { return zap.String("field_fn_key1", "field_fn_value1") })

	Debug("Debug")
	Info("Info")
	Warn("Warn")
	Info("info")
	Error("Error")
	DPanic("DPanic")

	Debugf("Debugf: %s", "debug")
	Infof("Infof: %s", "info")
	Warnf("Warnf: %s", "warn")
	Infof("Infof: %s", "info")
	Errorf("Errorf: %s", "error")
	DPanicf("DPanicf: %s", "dPanic")

	Debugw("Debugw", "Debugw", "w")
	Infow("Infow", "Infow", "w")
	Warnw("Warnw", "Warnw", "w")
	Infow("Infow", "Infow", "w")
	Errorw("Errorw", "Errorw", "w")
	DPanicw("DPanicw", "DPanicw", "w")

	shouPanic(t, func() {
		Panic("Panic")
	})
	shouPanic(t, func() {
		Panicf("Panicf: %s", "panic")
	})
	shouPanic(t, func() {
		Panicw("Panicw: %s", "panic", "w")
	})

	With(zap.String("aa", "bb")).Debug("debug with")

	Named("another").Debug("debug named")

	Logger().Debug("desugar")

	Inject(context.Background(),
		func(ctx context.Context) Field { return zap.String("field_fn_key2", "field_fn_value2") },
	).Debug("with context")

	WithValuer(func(ctx context.Context) Field { return zap.String("field_fn_key3", "field_fn_value3") }).
		Inject(context.Background()).
		Debug("with field fn")

	Logger().With(zap.Namespace("aaaa")).With(zap.String("xx", "yy")).Debug("----<>")

	_ = Sync()
}

func shouPanic(t *testing.T, f func()) {
	defer func() {
		e := recover()
		if e == nil {
			t.Errorf("should panic but not")
		}
	}()
	f()
}
