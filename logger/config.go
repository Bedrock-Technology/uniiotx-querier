package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

const (
	DefaultCallerSkip = 1
)

var (
	// For zap Logger

	DevMode        bool
	ConsoleEncoder bool
	Stacktrace     bool
	TopFrames      = 5

	// For lumberJack log rotation

	RollingFilename   = ""
	RollingMaxSize    = 500 // megabytes
	RollingMaxAge     = 30  // days
	RollingMaxBackups = 3
	RollingLocalTime  = false
	RollingCompress   = false
)

// NewZapLogger creates new a zap logger with given caller skip and some opinioned configs.
func NewZapLogger(callerSkip int) *zap.Logger {
	var logger *zap.Logger

	// Set Config
	var logConfig zap.Config
	logConfig = zap.NewProductionConfig()
	if DevMode {
		logConfig = zap.NewDevelopmentConfig()
	}

	logConfig.EncoderConfig.TimeKey = "ts"
	logConfig.EncoderConfig.LevelKey = "level"
	logConfig.EncoderConfig.NameKey = "logger"
	logConfig.EncoderConfig.CallerKey = "caller"
	logConfig.EncoderConfig.MessageKey = "msg"
	logConfig.EncoderConfig.StacktraceKey = "stacktrace"
	logConfig.DisableStacktrace = true
	logConfig.EncoderConfig.EncodeTime = iso8601UTCTimeEncoder
	if ConsoleEncoder {
		logConfig.Encoding = "console"
		logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		logConfig.Encoding = "json"
	}

	// Customize options
	var opts []zap.Option
	if RollingFilename != "" {
		// Customize zap Core
		coreW := zapcore.AddSync(&lumberjack.Logger{
			Filename:   RollingFilename,
			MaxSize:    RollingMaxSize,
			MaxAge:     RollingMaxAge,
			MaxBackups: RollingMaxBackups,
			LocalTime:  RollingLocalTime,
			Compress:   RollingCompress,
		})
		wrapCoreOpt := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			var encoder zapcore.Encoder
			if ConsoleEncoder {
				encoder = zapcore.NewConsoleEncoder(logConfig.EncoderConfig)
			} else {
				encoder = zapcore.NewJSONEncoder(logConfig.EncoderConfig)
			}
			newCore := zapcore.NewCore(
				encoder,
				coreW,
				logConfig.Level,
			)
			return newCore
		})
		opts = append(opts, wrapCoreOpt)
	}

	opts = append(opts, zap.AddCallerSkip(callerSkip))

	logger, _ = logConfig.Build(opts...)

	return logger
}

// A UTC variation of ZapCore.ISO8601TimeEncoder with nanosecond precision
func iso8601UTCTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z"))
}
