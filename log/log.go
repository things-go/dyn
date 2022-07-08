package log

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New constructs a new Log
func New(opts ...Option) (*zap.Logger, zap.AtomicLevel) {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	var options []zap.Option

	if c.AddCaller {
		// 添加显示文件名和行号,跳过封装调用层,
		options = append(options, zap.AddCaller(), zap.AddCallerSkip(c.CallerSkip))
	}
	if c.Stack {
		// 栈调用,及使能等级
		options = append(options,
			zap.AddStacktrace(zap.NewAtomicLevelAt(zap.DPanicLevel)), // 只显示栈的错误等级
		)
	}

	level, err := zap.ParseAtomicLevel(c.Level)
	if err != nil {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// 初始化core
	core := zapcore.NewCore(
		toEncoder(c, level), // 设置encoder
		toWriter(c),         // 设置输出
		level,               // 设置日志输出等级
	)
	return zap.New(core, options...), level
}

func toEncoder(c *Config, level zap.AtomicLevel) zapcore.Encoder {
	encoderConfig := c.EncoderConfig
	if encoderConfig == nil {
		encoderConfig = &zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    toEncodeLevel(c.EncodeLevel),
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		if level.Level() == zap.DebugLevel {
			encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		}
	}

	if c.Format == "console" {
		return zapcore.NewConsoleEncoder(*encoderConfig)
	}
	return zapcore.NewJSONEncoder(*encoderConfig)
}

func toEncodeLevel(l string) zapcore.LevelEncoder {
	switch l {
	case "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	case "LowercaseLevelEncoder": // 小写编码器(默认)
		fallthrough
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func toWriter(c *Config) zapcore.WriteSyncer {
	fileWriter := func() zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{ // 文件切割
			Filename:   filepath.Join(c.Path, c.Filename),
			MaxSize:    c.MaxSize,
			MaxAge:     c.MaxAge,
			MaxBackups: c.MaxBackups,
			LocalTime:  c.LocalTime,
			Compress:   c.Compress,
		})
	}
	stdoutWriter := func() zapcore.WriteSyncer {
		return zapcore.AddSync(os.Stdout)
	}
	customWriter := func(w ...zapcore.WriteSyncer) []zapcore.WriteSyncer {
		ws := make([]zapcore.WriteSyncer, 0, len(c.Writer)+len(w))

		for _, writer := range c.Writer {
			ws = append(ws, zapcore.AddSync(writer))
		}
		for _, writer := range w {
			ws = append(ws, zapcore.AddSync(writer))
		}
		return ws
	}
	switch strings.ToLower(c.Adapter) {
	case "file":
		return fileWriter()
	case "multi":
		return zapcore.NewMultiWriteSyncer(stdoutWriter(), fileWriter())
	case "file-custom":
		return zapcore.NewMultiWriteSyncer(customWriter(fileWriter())...)
	case "console-custom":
		return zapcore.NewMultiWriteSyncer(customWriter(stdoutWriter())...)
	case "multi-custom":
		return zapcore.NewMultiWriteSyncer(customWriter(stdoutWriter(), fileWriter())...)
	case "custom":
		ws := customWriter()
		if len(ws) == 0 {
			return stdoutWriter()
		}
		if len(ws) == 1 {
			return ws[0]
		}
		return zapcore.NewMultiWriteSyncer(ws...)
	default: // console
		return stdoutWriter()
	}
}
