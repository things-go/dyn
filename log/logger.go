package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Valuer is returns a log value.
type Valuer func(ctx context.Context) zap.Field

// Log wrap zap logger
type Log struct {
	log   *zap.Logger
	level zap.AtomicLevel
	fn    []Valuer
}

// NewLoggerWith new logger with zap logger and atomic level
func NewLoggerWith(logger *zap.Logger, lv zap.AtomicLevel) *Log {
	return &Log{
		logger,
		lv,
		nil,
	}
}

// NewLogger new logger
func NewLogger(opts ...Option) *Log { return NewLoggerWith(New(opts...)) }

// SetLevelWithText alters the logging level.
// ParseAtomicLevel set the logging level based on a lowercase or all-caps ASCII
// representation of the log level.
// If the provided ASCII representation is
// invalid an error is returned.
// see zapcore.Level
func (l *Log) SetLevelWithText(text string) error {
	lv, err := zapcore.ParseLevel(text)
	if err != nil {
		return err
	}
	l.level.SetLevel(lv)
	return nil
}

// SetLevel alters the logging level.
func (l *Log) SetLevel(lv zapcore.Level) *Log {
	l.level.SetLevel(lv)
	return l
}

// SetDefaultValuer set default Valuer function, which hold always until you call Inject.
func (l *Log) SetDefaultValuer(fs ...Valuer) *Log {
	fn := make([]Valuer, 0, len(fs)+len(l.fn))
	fn = append(fn, l.fn...)
	fn = append(fn, fs...)
	l.fn = fn
	return l
}

// Level returns the minimum enabled log level.
func (l *Log) Level() zapcore.Level { return l.level.Level() }

// Sugar wraps the Logger to provide a more ergonomic, but slightly slower,
// API. Sugaring a Logger is quite inexpensive, so it's reasonable for a
// single application to use both Loggers and SugaredLoggers, converting
// between them on the boundaries of performance-sensitive code.
func (l *Log) Sugar() *zap.SugaredLogger { return l.log.Sugar() }

// Logger return internal logger
func (l *Log) Logger() *zap.Logger { return l.log }

// WithValuer with Valuer function, until you call Inject
func (l *Log) WithValuer(fs ...Valuer) *Log {
	fn := make([]Valuer, 0, len(fs)+len(l.fn))
	fn = append(fn, l.fn...)
	fn = append(fn, fs...)
	return &Log{
		l.log,
		l.level,
		fn,
	}
}

// Inject return log with inject fn(ctx) field from context, which your set before.
func (l *Log) Inject(ctx context.Context, fs ...func(context.Context) zap.Field) *Log {
	fields := make([]zap.Field, 0, len(l.fn)+len(fs))
	for _, f := range l.fn {
		fields = append(fields, f(ctx))
	}
	for _, f := range fs {
		fields = append(fields, f(ctx))
	}
	return &Log{
		l.log.With(fields...),
		l.level,
		nil,
	}
}

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
func (l *Log) With(args ...any) *zap.SugaredLogger {
	return l.log.Sugar().With(args...)
}

// Named adds a sub-scope to the logger's name. See Log.Named for details.
func (l *Log) Named(name string) *zap.SugaredLogger {
	return l.log.Sugar().Named(name)
}

// Sync flushes any buffered log entries.
func (l *Log) Sync() error {
	return l.log.Sugar().Sync()
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *Log) Debug(args ...any) {
	l.log.Sugar().Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func (l *Log) Info(args ...any) {
	l.log.Sugar().Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *Log) Warn(args ...any) {
	l.log.Sugar().Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *Log) Error(args ...any) {
	l.log.Sugar().Error(args...)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *Log) DPanic(args ...any) {
	l.log.Sugar().DPanic(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *Log) Panic(args ...any) {
	l.log.Sugar().Panic(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *Log) Fatal(args ...any) {
	l.log.Sugar().Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l *Log) Debugf(template string, args ...any) {
	l.log.Sugar().Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *Log) Infof(template string, args ...any) {
	l.log.Sugar().Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *Log) Warnf(template string, args ...any) {
	l.log.Sugar().Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *Log) Errorf(template string, args ...any) {
	l.log.Sugar().Errorf(template, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *Log) DPanicf(template string, args ...any) {
	l.log.Sugar().DPanicf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *Log) Panicf(template string, args ...any) {
	l.log.Sugar().Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *Log) Fatalf(template string, args ...any) {
	l.log.Sugar().Fatalf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func (l *Log) Debugw(msg string, keysAndValues ...any) {
	l.log.Sugar().Debugw(msg, keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Log) Infow(msg string, keysAndValues ...any) {
	l.log.Sugar().Infow(msg, keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Log) Warnw(msg string, keysAndValues ...any) {
	l.log.Sugar().Warnw(msg, keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Log) Errorw(msg string, keysAndValues ...any) {
	l.log.Sugar().Errorw(msg, keysAndValues...)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func (l *Log) DPanicw(msg string, keysAndValues ...any) {
	l.log.Sugar().DPanicw(msg, keysAndValues...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func (l *Log) Panicw(msg string, keysAndValues ...any) {
	l.log.Sugar().Panicw(msg, keysAndValues...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func (l *Log) Fatalw(msg string, keysAndValues ...any) {
	l.log.Sugar().Fatalw(msg, keysAndValues...)
}
