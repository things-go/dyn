package zapl

import (
	"testing"
)

func TestNew(t *testing.T) {
	l := New(WithConfig(Config{Level: "debug", Format: "json"}))
	ReplaceGlobals(l.Sugar().With("test", "test"))

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

	Desugar().Debug("desugar")

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
