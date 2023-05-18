package logger

import (
	"io"
	"os"

	"go.uber.org/zap/zapcore"
	"phanes/config"
)

func Init() func() {
	var (
		l       = config.Conf.Collect.Log
		writers = make([]io.Writer, 0, 0)
	)
	if l.FileName == "" {
		panic("no log storage target")
	}

	if l.FileName != "" {
		writers = append(writers, FileOutputWriter("./logs", l.FileName, 50, 3))
	}
	writers = append(writers, os.Stderr)

	// set your log level here
	logger := ZapLog(zapcore.Level(l.LogLevel), writers...)
	InitLogger(logger)

	return func() {}
}
