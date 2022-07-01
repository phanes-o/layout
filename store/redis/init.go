package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	log "phanes/collector/logger"
)

var KV *redis.Client

func Init(connectAddr, pwd string) func() {
	client := redis.NewClient(&redis.Options{
		Addr:       connectAddr,
		Password:   pwd,
		DB:         0,
		PoolSize:   30,
		MaxRetries: 5,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	KV = client

	return func() {
		if err := client.Close(); err != nil {
			log.Error(err)
		}
	}
}
