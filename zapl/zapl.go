package zapl

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

	if c.Stack {
		// 添加显示文件名和行号,跳过封装调用层,栈调用,及使能等级
		stackLevel := zap.NewAtomicLevel()
		stackLevel.SetLevel(zap.WarnLevel) // 只显示栈的错误等级
		options = append(options,
			zap.AddCaller(),
			zap.AddCallerSkip(0),
			zap.AddStacktrace(stackLevel),
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
	switch strings.ToLower(c.Adapter) {
	case "file":
		return zapcore.AddSync(&lumberjack.Logger{ // 文件切割
			Filename:   filepath.Join(c.Path, c.Filename),
			MaxSize:    c.MaxSize,
			MaxAge:     c.MaxAge,
			MaxBackups: c.MaxBackups,
			LocalTime:  c.LocalTime,
			Compress:   c.Compress,
		})
	case "multi":
		return zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(&lumberjack.Logger{ // 文件切割
				Filename:   filepath.Join(c.Path, c.Filename),
				MaxSize:    c.MaxSize,
				MaxAge:     c.MaxAge,
				MaxBackups: c.MaxBackups,
				LocalTime:  c.LocalTime,
				Compress:   c.Compress,
			}))
	case "custom":
		ws := make([]zapcore.WriteSyncer, 0, len(c.Writer))

		for _, writer := range c.Writer {
			ws = append(ws, zapcore.AddSync(writer))
		}
		if len(ws) == 0 {
			return zapcore.AddSync(os.Stdout)
		}
		if len(ws) == 1 {
			return ws[0]
		}
		return zapcore.NewMultiWriteSyncer(ws...)
	case "console":
		fallthrough
	default:
		return zapcore.AddSync(os.Stdout)
	}
}
