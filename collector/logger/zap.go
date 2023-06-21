package logger

import (
	"context"
	"io"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"phanes/config"
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
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	conf := config.Conf.Collect.Log
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		NewBufferedWriteSyncer(conf.BufferSize, time.Duration(conf.Interval)*time.Second, io.MultiWriter(out...)),
		level,
	)

	logger := otelzap.New(
		zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2)),
		otelzap.WithCaller(true),
		otelzap.WithTraceIDField(true),
	)

	l := &ZapLog{
		logger: logger,
		writer: out,
	}
	return l
}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(config.Conf.Collect.Log.Prefix + " 2006-01-02 15:04:05.000"))
}
