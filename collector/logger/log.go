package logger

import (
	"context"

	"github.com/phanes-o/lib/otel/logger"
	"go.uber.org/zap/zapcore"
)

var defaultLogger logger.Logger

func initLogger(l logger.Logger) {
	defaultLogger = l
}

func GetLogger() logger.Logger {
	return defaultLogger
}

func WithContext(ctx context.Context) logger.Logger {
	return defaultLogger.Ctx(ctx)
}

func WithFields(fields ...zapcore.Field) logger.Logger {
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
