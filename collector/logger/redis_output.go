package logger

import (
	"context"
	"github.com/go-redis/redis/v8"
	"io"
)

type redisLogWriter struct {
	key string
	kv  *redis.Client
}

func RedisOutputWriter(kv *redis.Client, key string) io.Writer {
	return &redisLogWriter{
		kv:  kv,
		key: key,
	}
}

func (w *redisLogWriter) Write(b []byte) (int, error) {
	err := w.kv.RPush(context.Background(), w.key, b).Err()
	return len(b), err
}
