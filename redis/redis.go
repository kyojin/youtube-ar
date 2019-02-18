package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

func New(opts ...Option) (*redis.Client, error) {
	var o Options
	for i := range opts {
		opts[i](&o)
	}

	redisopts, err := redis.ParseURL(o.url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	redisopts.PoolSize = o.poolsize

	redis := redis.NewClient(redisopts)
	return redis, errors.WithStack(redis.Ping().Err())
}

type Option func(*Options)
type Options struct {
	url      string
	poolsize int
}

func URL(url string) Option        { return func(o *Options) { o.url = url } }
func PoolSize(poolsize int) Option { return func(o *Options) { o.poolsize = poolsize } }
