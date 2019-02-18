package pubsub

import (
	"context"

	goredis "github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/yansal/youtube-ar/redis"
)

type PubSub interface {
	Publish(message string) error
}

func New(ctx context.Context, opts ...redis.Option) (PubSub, error) {
	client, err := redis.New(opts...)
	return &pubsub{client: client}, err
}

type pubsub struct{ client *goredis.Client }

func (p *pubsub) Publish(message string) error {
	return errors.WithStack(p.client.Publish("pubsub:youtube-ar", message).Err())
}
