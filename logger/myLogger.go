package logger

import (
	"fmt"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var _ common.Logger = (*MyLogger)(nil)

type MyLogger struct {
	*zap.Logger
}

// DefaultProduction returns a reasonable production logger.
// Logging is enabled at InfoLevel and above.
// It uses a JSON encoder, writes to standard error, and enables sampling.
// Stacktraces are automatically included on logs of ErrorLevel and above.
func DefaultProduction() *MyLogger {
	DevMode = false
	ConsoleEncoder = false
	Stacktrace = true
	return &MyLogger{Logger: NewZapLogger(DefaultCallerSkip)}
}

// DefaultDevelopment returns a reasonable development logger.
// Logging is enabled at DebugLevel and above.
// It enables development mode (which makes DPanicLevel logs panic), uses a console encoder,
// writes to standard error, and disables sampling.
// Stacktraces are automatically included on logs of WarnLevel and above.
func DefaultDevelopment() *MyLogger {
	DevMode = true
	ConsoleEncoder = true
	Stacktrace = true
	return &MyLogger{Logger: NewZapLogger(DefaultCallerSkip)}
}

func New(zapLogger *zap.Logger) *MyLogger {
	return &MyLogger{Logger: zapLogger}
}

func (l *MyLogger) AddCallerSkip(skip int) common.Logger {
	clonedZapLogger := l.Logger.WithOptions(zap.AddCallerSkip(skip))
	return New(clonedZapLogger)
}

func (l *MyLogger) Debug(msg string, fields ...any) {
	l.Logger.Debug(msg, l.zapFields(fields...)...)
}

func (l *MyLogger) Info(msg string, fields ...any) {
	l.Logger.Info(msg, l.zapFields(fields...)...)
}

func (l *MyLogger) Warn(msg string, fields ...any) {
	l.Logger.Warn(msg, l.zapFields(fields...)...)
}

// Error logs error msg without error verbose
func (l *MyLogger) Error(msg string, err error, fields ...any) {
	l.Logger.Error(msg, l.appendErrorFields(err, l.zapFields(fields...))...)
}

func (l *MyLogger) Fatal(msg string, err error, fields ...any) {
	l.Logger.Fatal(msg, l.appendErrorFields(err, l.zapFields(fields...))...)
}

func (l *MyLogger) Sync() {
	if err := l.Logger.Sync(); err != nil {
		l.Logger.Error("failed to flush buffered log entries", l.appendErrorFields(err, nil)...)
		return
	}
	l.Logger.Info("buffered log entries are flushed")
}

// zapFields accepts K-V(string-any) pair(s) and return zap.Field array
func (l *MyLogger) zapFields(fields ...any) []zap.Field {
	zapFields := make([]zap.Field, 0)
	for i := 0; i < len(fields)-1; i = i + 2 {
		zapFields = append(zapFields, zap.Any(fields[i].(string), fields[i+1]))
	}
	return zapFields
}

func (l *MyLogger) appendErrorFields(err error, zapFields []zap.Field) []zap.Field {
	if err != nil {
		zapFields = append(zapFields, zap.String("error", err.Error()))
		zapFields = append(zapFields, zap.String("cause", errors.Cause(err).Error()))
		if Stacktrace {
			err, ok := err.(stackTracer)
			if ok {
				st := err.StackTrace()
				to := TopFrames
				if len(st) < TopFrames {
					to = len(st)
				}
				_st := fmt.Sprintf("%+v", st[0:to])
				zapFields = append(zapFields, zap.String("stacktrace", _st))
			}
		}
	}
	return zapFields
}

// stackTracer is for retrieving the stack trace of an error or wrapper.
// https://pkg.go.dev/github.com/pkg/errors
type stackTracer interface {
	StackTrace() errors.StackTrace
}
