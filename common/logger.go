package common

// Logger customized structured logger
type Logger interface {
	AddCallerSkip(skip int) Logger
	Sync()

	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, err error, fields ...any)
	Fatal(msg string, err error, fields ...any)
}
