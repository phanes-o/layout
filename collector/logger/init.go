package logger

import (
	"io"

	"github.com/phanes-o/lib/otel/logger"
	"phanes/config"
)

func Init() func() {
	var (
		l       = config.Conf.Collect.Log
		writers = make([]io.Writer, 0, 0)
	)
	if l.FileName == "" {
		panic("no logger storage target")
	}

	if l.FileName != "" {
		writers = append(writers, logger.NewFileWriter("./logs", l.FileName, 500, 3))
	}
	//writers = append(writers, os.Stderr)

	log := logger.NewZapLog(
		logger.WithLevel(l.Level),
		logger.WithWriters(writers...),
	)
	// set your logger level here
	initLogger(log)

	return func() {
		log.Sync()
	}
}
