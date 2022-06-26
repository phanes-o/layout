package logger

type Fields map[string]interface{}

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
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

var std Logger // standard output

func WithField(key string, value interface{}) Logger {
	return std.WithField(key, value)
}
func WithFields(fields Fields) Logger {
	return std.WithFields(fields)
}
func Trace(args ...interface{}) {
	std.Trace(args...)
}
func Tracef(format string, args ...interface{}) {
	std.Tracef(format, args...)
}
func Debug(args ...interface{}) {
	std.Debug(args...)
}
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}
func Info(args ...interface{}) {
	std.Info(args...)
}
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}
func Warn(args ...interface{}) {
	std.Warn(args...)
}
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}
func Error(args ...interface{}) {
	std.Error(args...)
}
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}
func Panic(args ...interface{}) {
	std.Panic(args...)
}
func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}
func Fatal(args ...interface{}) {
	std.Fatal(args...)
}
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

func InitGlobal(s Logger) {
	std = s
}
