package logger

import (
	"context"

	"go.uber.org/zap/zapcore"
)

var defaultLogger Logger

func InitLogger(l Logger) {
	defaultLogger = l
}

type Logger interface {
	WithFields(fields ...zapcore.Field) Logger
	Ctx(ctx context.Context) Logger
	Debug(msg string, fields ...zapcore.Field)
	DebugCtx(ctx context.Context, msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	InfoCtx(ctx context.Context, msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	WarnCtx(ctx context.Context, msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	ErrorCtx(ctx context.Context, msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	FatalCtx(ctx context.Context, msg string, fields ...zapcore.Field)
	PanicCtx(ctx context.Context, msg string, fields ...zapcore.Field)
}

func WithContext(ctx context.Context) Logger {
	return defaultLogger.Ctx(ctx)
}

func WithFields(fields ...zapcore.Field) Logger {
	return defaultLogger.WithFields(fields...)
}

func Debug(msg string, fields ...zapcore.Field) {
	defaultLogger.Debug(msg, fields...)
}

func DebugCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	defaultLogger.DebugCtx(ctx, msg, fields...)
}

func Info(msg string, fields ...zapcore.Field) {
	defaultLogger.Info(msg, fields...)
}

func InfoCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	defaultLogger.InfoCtx(ctx, msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	defaultLogger.Warn(msg, fields...)
}

func WarnCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	defaultLogger.WarnCtx(ctx, msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	defaultLogger.Error(msg, fields...)
}
func ErrorCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	defaultLogger.ErrorCtx(ctx, msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	defaultLogger.Fatal(msg, fields...)
}

func FatalCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	defaultLogger.FatalCtx(ctx, msg, fields...)
}

func PanicCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	defaultLogger.PanicCtx(ctx, msg, fields...)
}
