package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger = NewLoggerWith(zap.NewNop(), zap.NewAtomicLevel())

// ReplaceGlobals replaces the global Log,
func ReplaceGlobals(logger *Log) { defaultLogger = logger }

// SetLevelWithText alters the logging level.
// ParseAtomicLevel set the logging level based on a lowercase or all-caps ASCII
// representation of the log level.
// If the provided ASCII representation is
// invalid an error is returned.
func SetLevelWithText(text string) error { return defaultLogger.SetLevelWithText(text) }

// SetLevel alters the logging level.
func SetLevel(lv zapcore.Level) *Log { return defaultLogger.SetLevel(lv) }

// Level returns the minimum enabled log level.
func Level() zapcore.Level { return defaultLogger.Level() }

// SetDefaultFieldFn set default field function, which hold always until you call Inject.
// Deprecated: use SetDefaultValuer instead
func SetDefaultFieldFn(vs ...Valuer) *Log {
	return defaultLogger.SetDefaultValuer(vs...)
}

// SetDefaultValuer set default field function, which hold always until you call Inject.
func SetDefaultValuer(vs ...Valuer) *Log {
	return defaultLogger.SetDefaultValuer(vs...)
}

// WithFieldFn with field function, until you call Inject.
// Deprecated: use WithValuer instead
func WithFieldFn(vs ...Valuer) *Log {
	return defaultLogger.WithValuer(vs...)
}

// WithValuer with field function, until you call Inject
func WithValuer(vs ...Valuer) *Log {
	return defaultLogger.WithValuer(vs...)
}

// Inject return log with inject fn(ctx) field from context, which your set before.
func Inject(ctx context.Context, fs ...func(context.Context) zap.Field) *Log {
	return defaultLogger.Inject(ctx, fs...)
}

// Sugar wraps the Logger to provide a more ergonomic, but slightly slower,
// API. Sugaring a Logger is quite inexpensive, so it's reasonable for a
// single application to use both Loggers and SugaredLoggers, converting
// between them on the boundaries of performance-sensitive code.
func Sugar() *zap.SugaredLogger { return defaultLogger.Sugar() }

// Logger return internal logger
func Logger() *zap.Logger { return defaultLogger.Logger() }

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
func With(args ...interface{}) *zap.SugaredLogger { return defaultLogger.Sugar().With(args...) }

// Named adds a sub-scope to the logger's name. See Log.Named for details.
func Named(name string) *zap.SugaredLogger { return defaultLogger.Named(name) }

// Sync flushes any buffered log entries.
func Sync() error { return defaultLogger.Sync() }

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) { defaultLogger.Debug(args...) }

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) { defaultLogger.Info(args...) }

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) { defaultLogger.Warn(args...) }

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) { defaultLogger.Error(args...) }

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanic(args ...interface{}) { defaultLogger.DPanic(args...) }

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) { defaultLogger.Panic(args...) }

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) { defaultLogger.Fatal(args...) }

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) { defaultLogger.Debugf(template, args...) }

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) { defaultLogger.Infof(template, args...) }

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) { defaultLogger.Warnf(template, args...) }

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) { defaultLogger.Errorf(template, args...) }

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanicf(template string, args ...interface{}) { defaultLogger.DPanicf(template, args...) }

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) { defaultLogger.Panicf(template, args...) }

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) { defaultLogger.Fatalf(template, args...) }

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func Debugw(msg string, keysAndValues ...interface{}) { defaultLogger.Debugw(msg, keysAndValues...) }

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) { defaultLogger.Infow(msg, keysAndValues...) }

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) { defaultLogger.Warnw(msg, keysAndValues...) }

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) { defaultLogger.Errorw(msg, keysAndValues...) }

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func DPanicw(msg string, keysAndValues ...interface{}) { defaultLogger.DPanicw(msg, keysAndValues...) }

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) { defaultLogger.Panicw(msg, keysAndValues...) }

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) { defaultLogger.Fatalw(msg, keysAndValues...) }
