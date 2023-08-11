package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"phanes/config"
	"phanes/model"
)

type zapLog struct {
	logger *otelzap.Logger
	writer []io.Writer
	prefix string
}

func (z *zapLog) WithFields(fields ...zapcore.Field) Logger {
	newLogger := z.logger.With(fields...)
	newZap := &zapLog{
		logger: otelzap.New(newLogger),
		writer: z.writer,
	}
	return newZap
}

func (z *zapLog) Ctx(ctx context.Context) Logger {
	newCtx := z.logger.Ctx(ctx)
	newZap := &zapLog{
		logger: newCtx.Logger(),
		writer: z.writer,
	}
	return newZap
}

func (z *zapLog) Debug(msg string, fields ...zapcore.Field) {
	z.logger.Debug(z.withPrefix(msg), fields...)
}

func (z *zapLog) DebugCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.DebugContext(ctx, z.withPrefix(msg), fields...)
}

func (z *zapLog) Info(msg string, fields ...zapcore.Field) {
	z.logger.Info(z.withPrefix(msg), fields...)
}

func (z *zapLog) InfoCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.InfoContext(ctx, z.withPrefix(msg), fields...)
}

func (z *zapLog) Warn(msg string, fields ...zapcore.Field) {
	z.logger.Warn(z.withPrefix(msg), fields...)
}

func (z *zapLog) WarnCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.WarnContext(ctx, z.withPrefix(msg), fields...)
}

func (z *zapLog) Error(msg string, fields ...zapcore.Field) {
	z.logger.Error(z.withPrefix(msg), fields...)
}

func (z *zapLog) ErrorCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.ErrorContext(ctx, z.withPrefix(msg), fields...)
}

func (z *zapLog) Fatal(msg string, fields ...zapcore.Field) {
	z.logger.Fatal(z.withPrefix(msg), fields...)
}

func (z *zapLog) FatalCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.FatalContext(ctx, z.withPrefix(msg), fields...)
}

func (z *zapLog) PanicCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.PanicContext(ctx, z.withPrefix(msg), fields...)
}

func (z *zapLog) withPrefix(msg string) string {
	return z.prefix + " " + msg
}

func newZapLog(level zapcore.Level, out ...io.Writer) *zapLog {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder

	syncer := newBufferedWriteSyncer(
		config.Conf.Collect.Log.BufferSize,
		time.Duration(config.Conf.Collect.Log.Interval)*time.Second,
		io.MultiWriter(out...),
	)

	outputCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), syncer, level)
	var core zapcore.Core
	switch config.Conf.Env {
	case model.EnvProd:
		core = zapcore.NewTee(outputCore)
	case model.EnvDev:
		consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), os.Stdout, level)
		core = zapcore.NewTee(consoleCore, outputCore)
	}

	otelLogger := otelzap.New(
		zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2)),
		otelzap.WithCaller(true),
		otelzap.WithTraceIDField(true),
	)

	logger := &zapLog{
		logger: otelLogger,
		writer: out,
	}

	if config.Conf.Collect.Log.Prefix != "" {
		logger.prefix = config.Conf.Collect.Log.Prefix
	}

	return logger
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
