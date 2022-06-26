package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"phanes/config"
)

func Init() func() {
	var (
		l       = config.Conf.Log
		writers = make([]io.Writer, 0, 0)
	)
	if l.FileName == "" && l.RedisKey == "" {
		panic("no log storage target")
	}

	if l.FileName != "" {
		writers = append(writers, FileOutputWriter("./logs", l.FileName, 50, 3))
	}
	if l.RedisKey != "" {
		writers = append(writers, RedisOutputWriter(config.KV, l.RedisKey))
	}

	log := NewLogrusAdapt(logrus.DebugLevel, true, writers...)
	InitGlobal(log)

	return func() {}
}
