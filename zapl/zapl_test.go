package zapl

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	l, lv := New(WithConfig(Config{Level: "debug", Format: "json"}))
	ReplaceGlobals(NewLoggerWith(l, lv))
	SetDefaultFieldFn(func(ctx context.Context) zap.Field { return zap.String("field_fn_key1", "field_fn_value1") })

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

	Debugw("Debugw: %s", "debug")
	Infow("Infow: %s", "info")
	Warnw("Warnw: %s", "warn")
	Infow("Infow: %s", "info")
	Errorw("Errorw: %s", "error")
	DPanicw("DPanicw: %s", "dPanic")

	shouPanic(t, func() {
		Panic("Panic")
	})
	shouPanic(t, func() {
		Panicf("Panicf: %s", "panic")
	})
	shouPanic(t, func() {
		Panicw("Panicw: %s", "panic")
	})

	With("aa", "bb").Debug("debug with")

	Named("another").Debug("debug named")

	Logger().Debug("desugar")

	WithContext(context.Background(),
		func(ctx context.Context) zap.Field { return zap.String("field_fn_key2", "field_fn_value2") },
	).Debug("with context")

	WithFieldFn(func(ctx context.Context) zap.Field { return zap.String("field_fn_key3", "field_fn_value3") }).
		Inject(context.Background()).
		Debug("with field fn")

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
