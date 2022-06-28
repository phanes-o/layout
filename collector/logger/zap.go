package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

type zopLog struct {
	ctx    context.Context
	logger *zap.SugaredLogger
	writer []io.Writer
}

func ZapLog(level zapcore.Level, out ...io.Writer) Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	w := zapcore.AddSync(io.MultiWriter(out...))
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), w, level)
	z := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
	l := &zopLog{
		logger: z.Sugar(),
		writer: out,
	}
	return l
}

func (z *zopLog) WithContext(ctx context.Context) Logger {
	newZap := &zopLog{
		ctx:    ctx,
		logger: z.logger,
		writer: z.writer,
	}
	return newZap
}

func (z *zopLog) WithField(key string, value interface{}) Logger {
	field := zap.Any(key, value)
	newLogger := z.logger.With(field)
	newZop := &zopLog{
		logger: newLogger,
		ctx:    z.ctx,
		writer: z.writer,
	}
	return newZop
}

func (z *zopLog) WithFields(fields Fields) Logger {
	zapFields := make([]interface{}, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	newLogger := z.logger.With(zapFields...)

	newZop := &zopLog{
		logger: newLogger,
		ctx:    z.ctx,
		writer: z.writer,
	}
	return newZop
}

func (z *zopLog) Trace(args ...interface{}) {
	z.logger.Info(args)
}

func (z *zopLog) Tracef(format string, args ...interface{}) {
	z.logger.Info(args)
}

func (z *zopLog) Debug(args ...interface{}) {
	z.logger.Debug(args)
}

func (z *zopLog) Debugf(format string, args ...interface{}) {
	z.logger.Debugf(format, args)
}

func (z *zopLog) Info(args ...interface{}) {
	z.logger.Info(args)
}

func (z *zopLog) Infof(format string, args ...interface{}) {
	z.logger.Infof(format, args)
}

func (z *zopLog) Warn(args ...interface{}) {
	z.logger.Warn(args)
}

func (z *zopLog) Warnf(format string, args ...interface{}) {
	z.logger.Warnf(format, args)
}

func (z *zopLog) Error(args ...interface{}) {
	z.logger.Error(args)
}

func (z *zopLog) Errorf(format string, args ...interface{}) {
	z.logger.Errorf(format, args)
}

func (z *zopLog) Panic(args ...interface{}) {
	z.logger.Panic(args)
}

func (z *zopLog) Panicf(format string, args ...interface{}) {
	z.logger.Panicf(format, args)
}

func (z *zopLog) Fatal(args ...interface{}) {
	z.logger.Fatal(args)
}

func (z *zopLog) Fatalf(format string, args ...interface{}) {
	z.logger.Fatalf(format, args)
}
