package logger

import "context"

var defaultLogger Logger

func InitLogger(l Logger) {
	defaultLogger = l
}

type Fields map[string]interface{}

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	WithContext(ctx context.Context) Logger
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

func WithContext(ctx context.Context) Logger {
	return defaultLogger.WithContext(ctx)
}
func WithField(key string, value interface{}) Logger {
	return defaultLogger.WithField(key, value)
}

func WithFields(fields Fields) Logger {
	return defaultLogger.WithFields(fields)
}

func Trace(args ...interface{}) {
	defaultLogger.Trace(args)
}
func Tracef(format string, args ...interface{}) {
	defaultLogger.Tracef(format, args)
}
func Debug(args ...interface{}) {
	defaultLogger.Debug(args)
}
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args)
}
func Info(args ...interface{}) {
	defaultLogger.Info(args)
}
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args)
}
func Warn(args ...interface{}) {
	defaultLogger.Warn(args)
}
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args)
}
func Error(args ...interface{}) {
	defaultLogger.Error(args)
}
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args)
}
func Panic(args ...interface{}) {
	defaultLogger.Panic(args)
}
func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args)
}
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args)
}
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args)
}
