package zapl

import (
	"go.uber.org/zap"
)

var defaultSugar = zap.NewNop().Sugar()

// ReplaceGlobals replaces the global SugaredLogger,
func ReplaceGlobals(s *zap.SugaredLogger) { defaultSugar = s }

// Desugar unwraps a SugaredLogger, exposing the original Logger. Desugaring
// is quite inexpensive, so it's reasonable for a single application to use
// both Loggers and SugaredLoggers, converting between them on the boundaries
// of performance-sensitive code.
func Desugar() *zap.Logger { return defaultSugar.Desugar() }

// With adds a variadic number of fields to the logging context. It accepts a
// mix of strongly-typed Field objects and loosely-typed key-value pairs. When
// processing pairs, the first element of the pair is used as the field key
// and the second as the field value.
//
// For example,
//   sugaredLogger.With(
//     "hello", "world",
//     "failure", errors.New("oh no"),
//     Stack(),
//     "count", 42,
//     "user", User{Name: "alice"},
//  )
// is the equivalent of
//   unsugared.With(
//     String("hello", "world"),
//     String("failure", "oh no"),
//     Stack(),
//     Int("count", 42),
//     Object("user", User{Name: "alice"}),
//   )
//
// Note that the keys in key-value pairs should be strings. In development,
// passing a non-string key panics. In production, the logger is more
// forgiving: a separate error is logged, but the key-value pair is skipped
// and execution continues. Passing an orphaned key triggers similar behavior:
// panics in development and errors in production.
func With(args ...interface{}) *zap.SugaredLogger { return defaultSugar.With(args...) }

// Named adds a sub-scope to the logger's name. See Logger.Named for details.
func Named(name string) *zap.SugaredLogger { return defaultSugar.Named(name) }

// Sync flushes any buffered log entries.
func Sync() error { return defaultSugar.Sync() }

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) { defaultSugar.Debug(args...) }

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) { defaultSugar.Info(args...) }

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) { defaultSugar.Warn(args...) }

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) { defaultSugar.Error(args...) }

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanic(args ...interface{}) { defaultSugar.DPanic(args...) }

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) { defaultSugar.Panic(args...) }

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) { defaultSugar.Fatal(args...) }

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) { defaultSugar.Debugf(template, args...) }

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) { defaultSugar.Infof(template, args...) }

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) { defaultSugar.Warnf(template, args...) }

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) { defaultSugar.Errorf(template, args...) }

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanicf(template string, args ...interface{}) { defaultSugar.DPanicf(template, args...) }

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) { defaultSugar.Panicf(template, args...) }

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) { defaultSugar.Fatalf(template, args...) }

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func Debugw(msg string, keysAndValues ...interface{}) { defaultSugar.Debugw(msg, keysAndValues...) }

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) { defaultSugar.Infow(msg, keysAndValues...) }

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) { defaultSugar.Warnw(msg, keysAndValues...) }

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) { defaultSugar.Errorw(msg, keysAndValues...) }

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func DPanicw(msg string, keysAndValues ...interface{}) { defaultSugar.DPanicw(msg, keysAndValues...) }

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) { defaultSugar.Panicw(msg, keysAndValues...) }

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) { defaultSugar.Fatalw(msg, keysAndValues...) }
