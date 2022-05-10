package zapl

import (
	"io"

	"go.uber.org/zap/zapcore"
)

// Config 日志配置
type Config struct {
	// Level 日志等级, debug,info,warn,error,dpanic,panic,fatal, 默认warn
	Level string `yaml:"level" json:"level"`
	// Format: 编码格式: json,console 默认json
	Format string `yaml:"format" json:"format"`
	// 编码器类型, 默认: LowercaseLevelEncoder
	// LowercaseLevelEncoder: 小写编码器
	// LowercaseColorLevelEncoder: 小写编码器带颜色
	// CapitalLevelEncoder: 大写编码器
	// CapitalColorLevelEncoder: 大写编码器带颜色
	EncodeLevel string `yaml:"encodeLevel" json:"encodeLevel"`
	// Adapter 输出适配器, file,console,multi,custom, 默认 console
	Adapter string `yaml:"adapter" json:"adapter"`
	// Stack 是否使能栈调试输出, 默认false
	Stack bool `yaml:"stack" json:"stack"`
	// Path 日志保存路径, 默认 empty, 即当前路径
	Path string `yaml:"path" json:"path"`
	// Writer 输出
	// 当 adapter=custom使用,如果为writer为空,将使用os.Stdout
	Writer []io.Writer `yaml:"-" json:"-"`
	// EncoderConfig 如果配置该项,则 EncodeLevel 将被覆盖
	EncoderConfig *zapcore.EncoderConfig `yaml:"-" json:"-"`

	// see https://github.com/natefinch/lumberjack
	// lumberjack.Logger
	// Filename 空字符使用默认, 默认<processname>-lumberjack.log
	Filename string `yaml:"filename" json:"filename"`
	// MaxSize 每个日志文件最大尺寸(MB), 默认100MB
	MaxSize int `yaml:"maxSize" json:"maxSize"`
	// MaxAge 日志文件保存天数, 默认0 不删除
	MaxAge int `yaml:"maxAge" json:"maxAge"`
	// MaxBackups 日志文件保存备份数, 默认0 都保存
	MaxBackups int `yaml:"maxBackups" json:"maxBackups"`
	// LocalTime 是否格式化时间戳, 默认UTC时间
	LocalTime bool `yaml:"localTime" json:"localTime"`
	// Compress 是否使用gzip压缩文件, 采用默认不压缩
	Compress bool `yaml:"compress" json:"compress"`
}

// Option An Option configures a Logger.
type Option func(c *Config)

// WithConfig with config
func WithConfig(cfg Config) Option {
	return func(c *Config) { *c = cfg }
}

// WithLevel with level
// debug,info,warn,error,dpanic,panic,fatal
// 默认warn
func WithLevel(level string) Option {
	return func(c *Config) { c.Level = level }
}

// WithFormat with format
// json,console
// 默认json
func WithFormat(format string) Option {
	return func(c *Config) { c.Format = format }
}

// WithEncodeLevel with EncodeLevel
// LowercaseLevelEncoder: 小写编码器
// LowercaseColorLevelEncoder: 小写编码器带颜色
// CapitalLevelEncoder: 大写编码器
// CapitalColorLevelEncoder: 大写编码器带颜色
// 默认: LowercaseLevelEncoder
func WithEncodeLevel(encodeLevel string) Option {
	return func(c *Config) { c.EncodeLevel = encodeLevel }
}

// EncoderConfig 如果配置该项,则 EncodeLevel 将被覆盖
func WithEncoderConfig(encoderConfig *zapcore.EncoderConfig) Option {
	return func(c *Config) { c.EncoderConfig = encoderConfig }
}

// WithAdapter with adapter
// file,console,multi,custom
// writer: 当 adapter=custom使用,如果为writer为空,将使用os.Stdout
// 默认 console
func WithAdapter(adapter string, writer ...io.Writer) Option {
	return func(c *Config) {
		c.Adapter = adapter
		c.Writer = writer
	}
}

// WithStack with stack
// Stack 是否使能栈调试输出, 默认false
func WithStack(stack bool) Option {
	return func(c *Config) { c.Stack = stack }
}

// WithPath with path
// 日志保存路径, 默认 empty, 即当前路径
func WithPath(path string) Option {
	return func(c *Config) { c.Path = path }
}

/******************************** lumberjack.Logger **************************************/

// WithFilename with filename
// 空字符使用默认, 默认<processname>-lumberjack.log
func WithFilename(filename string) Option {
	return func(c *Config) { c.Filename = filename }
}

// WithMaxSize with max size
// 每个日志文件最大尺寸(MB), 默认100MB
func WithMaxSize(maxSize int) Option {
	return func(c *Config) { c.MaxSize = maxSize }
}

// WithMaxAge with max age
// 日志文件保存天数, 默认0 不删除
func WithMaxAge(maxAge int) Option {
	return func(c *Config) { c.MaxAge = maxAge }
}

// WithMaxBackups with max backup
// 日志文件保存备份数, 默认0 都保存
func WithMaxBackups(maxBackups int) Option {
	return func(c *Config) { c.MaxBackups = maxBackups }
}

// WithEnableLocalTime with local time
// 是否格式化时间戳, 默认UTC时间
func WithEnableLocalTime() Option {
	return func(c *Config) { c.LocalTime = true }
}

// WithEnableCompress with compress
// 是否使用gzip压缩文件, 采用默认不压缩
func WithEnableCompress() Option {
	return func(c *Config) { c.Compress = true }
}
