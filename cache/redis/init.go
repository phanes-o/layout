package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/phanes-o/utils"
	"go-micro.dev/v4/cache"
)

type redisCache struct {
	ctx  context.Context
	Kv   *redis.Client
	opts Options
}

func (r *redisCache) Context(ctx context.Context) cache.Cache {
	r.ctx = ctx
	return r
}

func (r *redisCache) Get(key string) (interface{}, time.Time, error) {
	val, err := r.Kv.Get(r.ctx, key).Bytes()
	if err != nil && err == redis.Nil {
		return nil, time.Time{}, cache.ErrKeyNotFound
	} else if err != nil {
		return nil, time.Time{}, err
	}

	dur, err := r.Kv.TTL(r.ctx, key).Result()
	if err != nil {
		return nil, time.Time{}, err
	}
	if dur == -1 {
		return val, time.Unix(1<<63-1, 0), nil
	}
	if dur == -2 {
		return val, time.Time{}, cache.ErrItemExpired
	}

	return val, time.Now().Add(dur), nil
}

func (r *redisCache) Put(key string, val interface{}, d time.Duration) error {
	return r.Kv.Set(r.ctx, key, val, d).Err()
}

func (r *redisCache) Delete(key string) error {
	return r.Kv.Del(r.ctx, key).Err()
}

type Options struct {
	addr     string
	pwd      string
	poolSize int
	db       int
}

type Option interface {
	Apply(*Options)
}

type optionFunc func(*Options)

func (f optionFunc) Apply(o *Options) {
	f(o)
}

func WithAddr(addr string) Option {
	return optionFunc(func(o *Options) {
		o.addr = addr
	})
}

func WithPwd(pwd string) Option {
	return optionFunc(func(o *Options) {
		o.pwd = pwd
	})
}

func WithPoolSize(poolSize int) Option {
	return optionFunc(func(o *Options) {
		o.poolSize = poolSize
	})
}

func WithDB(db int) Option {
	return optionFunc(func(o *Options) {
		o.db = db
	})
}

func NewRedisCache(opts ...Option) cache.Cache {
	options := &Options{
		db:       0,
		addr:     "127.0.0.1:6379",
		poolSize: 30,
	}

	for _, o := range opts {
		o.Apply(options)
	}

	kv := redis.NewClient(&redis.Options{
		Addr:     options.addr,
		Password: options.pwd,
		DB:       options.db,
		PoolSize: options.poolSize,
		MaxRetries: 5,
	})

	utils.Throw(kv.Ping(context.Background()).Err())

	return &redisCache{
		Kv:   kv,
		opts: *options,
	}
}
