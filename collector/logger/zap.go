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

type ZapLog struct {
	logger *otelzap.Logger
	writer []io.Writer
}

func (z *ZapLog) WithFields(fields ...zapcore.Field) Logger {
	newLogger := z.logger.With(fields...)
	newZap := &ZapLog{
		logger: otelzap.New(newLogger),
		writer: z.writer,
	}
	return newZap
}

func (z *ZapLog) Ctx(ctx context.Context) Logger {
	newCtx := z.logger.Ctx(ctx)
	newZap := &ZapLog{
		logger: newCtx.Logger(),
		writer: z.writer,
	}
	return newZap
}

func (z *ZapLog) Debug(msg string, fields ...zapcore.Field) {
	z.logger.Debug(msg, fields...)
}

func (z *ZapLog) DebugCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.DebugContext(ctx, msg, fields...)
}

func (z *ZapLog) Info(msg string, fields ...zapcore.Field) {
	z.logger.Info(msg, fields...)
}

func (z *ZapLog) InfoCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.InfoContext(ctx, msg, fields...)
}

func (z *ZapLog) Warn(msg string, fields ...zapcore.Field) {
	z.logger.Warn(msg, fields...)
}

func (z *ZapLog) WarnCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.WarnContext(ctx, msg, fields...)
}

func (z *ZapLog) Error(msg string, fields ...zapcore.Field) {
	z.logger.Error(msg, fields...)
}

func (z *ZapLog) ErrorCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.ErrorContext(ctx, msg, fields...)
}

func (z *ZapLog) Fatal(msg string, fields ...zapcore.Field) {
	z.logger.Fatal(msg, fields...)
}

func (z *ZapLog) FatalCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.FatalContext(ctx, msg, fields...)
}

func (z *ZapLog) PanicCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.PanicContext(ctx, msg, fields...)
}

func NewZapLog(level zapcore.Level, out ...io.Writer) *ZapLog {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = CustomTimeEncoder

	syncer := NewBufferedWriteSyncer(
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

	logger := &ZapLog{
		logger: otelLogger,
		writer: out,
	}
	return logger
}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(config.Conf.Collect.Log.Prefix + " 2006-01-02 15:04:05.000"))
}
