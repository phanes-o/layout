package cache

import (
	"go-micro.dev/v4/cache"
	"phanes/cache/redis"
	"phanes/config"
	"phanes/errors"
)

var iCache cache.Cache

var NotEnabled = errors.New("cache not enabled")

func Init() func() {
	conf := config.Conf.Cache
	if !conf.Enabled {
		return func() {}
	}

	switch conf.Type {
	case "redis":
		iCache = redis.NewRedisCache(
			redis.WithDB(0),
			redis.WithAddr(conf.Addr),
			redis.WithPwd(conf.Pwd),
			redis.WithPoolSize(30),
		)
	default:
		panic("cache type not support")
	}

	return func() {}
}

func GetCache() (cache.Cache, error) {
	if config.Conf.Cache.Enabled {
		return iCache, nil
	}
	return nil, NotEnabled
}
