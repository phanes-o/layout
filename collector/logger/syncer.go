package logger

import (
	"io"
	"time"

	"go.uber.org/zap/zapcore"
)

func newBufferedWriteSyncer(size int, duration time.Duration, w io.Writer) zapcore.WriteSyncer {
	bufferedSyncer := &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(w),
		Size:          size,
		FlushInterval: duration,
	}
	return bufferedSyncer
}
