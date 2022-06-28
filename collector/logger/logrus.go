package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type logrusLog struct {
	path   string
	logger *logrus.Logger
	writer io.Writer
}

func NewLogger(level logrus.Level, stdout bool, writes ...io.Writer) Logger {
	l := &logrusLog{}
	writer := io.MultiWriter(writes...)
	l.writer = writer

	l.logger = logrus.New()
	l.logger.Level = level
	l.logger.SetReportCaller(true)

	if stdout {
		l.logger.Out = io.MultiWriter(os.Stdout, l.writer)
	} else {
		l.logger.Out = l.writer
	}

	l.logger.SetFormatter(&logrus.TextFormatter{ForceColors: false, DisableColors: false, FieldMap: logrus.FieldMap{
		logrus.FieldKeyTime:  "time",
		logrus.FieldKeyFunc:  "func",
		logrus.FieldKeyFile:  "file",
		logrus.FieldKeyLevel: "level",
		logrus.FieldKeyMsg:   "message"}})
	return l
}

func (l *logrusLog) WithField(ctx context.Context, key string, value interface{}) Logger {
	l.logger.WithField(key, value)
	return l
}

func (l *logrusLog) WithFields(ctx context.Context, fields Fields) Logger {
	lf := logrus.Fields{}
	for k, v := range fields {
		lf[k] = v
	}
	l.logger.WithFields(lf)
	return l
}

func (l *logrusLog) Trace(ctx context.Context, args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *logrusLog) Tracef(ctx context.Context, format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

func (l *logrusLog) Debug(ctx context.Context, args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *logrusLog) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logrusLog) Info(ctx context.Context, args ...interface{}) {
	l.logger.Info(args...)
}

func (l *logrusLog) Infof(ctx context.Context, format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logrusLog) Warn(ctx context.Context, args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *logrusLog) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logrusLog) Error(ctx context.Context, args ...interface{}) {
	l.logger.Error(args...)
}

func (l *logrusLog) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *logrusLog) Panic(ctx context.Context, args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *logrusLog) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l *logrusLog) Fatal(ctx context.Context, args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *logrusLog) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
