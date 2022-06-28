package logger

import "context"

var defaultLogger Logger

func InitLogger(l Logger) {
	defaultLogger = l
}

type Fields map[string]interface{}

type Logger interface {
	WithField(ctx context.Context, key string, value interface{}) Logger
	WithFields(ctx context.Context, fields Fields) Logger
	Trace(ctx context.Context, args ...interface{})
	Tracef(ctx context.Context, format string, args ...interface{})
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Panic(ctx context.Context, args ...interface{})
	Panicf(ctx context.Context, format string, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
}

func WithField(ctx context.Context, key string, value interface{}) Logger {
	return defaultLogger.WithField(ctx, key, value)
}

func WithFields(ctx context.Context, fields Fields) Logger {
	return defaultLogger.WithFields(ctx, fields)
}

func Trace(ctx context.Context, args ...interface{}) {
	defaultLogger.Trace(ctx, args)
}
func Tracef(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Tracef(ctx, format, args)
}
func Debug(ctx context.Context, args ...interface{}) {
	defaultLogger.Debug(ctx, args)
}
func Debugf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Debugf(ctx, format, args)
}
func Info(ctx context.Context, args ...interface{}) {
	defaultLogger.Info(ctx, args)
}
func Infof(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Infof(ctx, format, args)
}
func Warn(ctx context.Context, args ...interface{}) {
	defaultLogger.Warn(ctx, args)
}
func Warnf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Warnf(ctx, format, args)
}
func Error(ctx context.Context, args ...interface{}) {
	defaultLogger.Error(ctx, args)
}
func Errorf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Errorf(ctx, format, args)
}
func Panic(ctx context.Context, args ...interface{}) {
	defaultLogger.Panic(ctx, args)
}
func Panicf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Panicf(ctx, format, args)
}
func Fatal(ctx context.Context, args ...interface{}) {
	defaultLogger.Fatal(ctx, args)
}
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Fatalf(ctx, format, args)
}
