package broker

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/yansal/q"
	"github.com/yansal/youtube-ar/model"
	"github.com/yansal/youtube-ar/redis"
)

type Broker interface {
	Send(context.Context, *model.URL) error
	Handle(context.Context, Handler) error
}

type Handler func(context.Context, *model.URL) error

func New(ctx context.Context, opts ...redis.Option) (Broker, error) {
	redis, err := redis.New(opts...)
	return &broker{q: q.New(redis)}, err
}

type broker struct{ q q.Q }

const queue = "youtube-ar"

func (b *broker) Send(ctx context.Context, m *model.URL) error {
	j, _ := json.Marshal(m)
	return b.q.Send(ctx, queue, string(j))
}

func (b *broker) Handle(ctx context.Context, h Handler) error {
	return b.q.Receive(ctx, queue, func(ctx context.Context, payload string) error {
		var m model.URL
		if err := json.Unmarshal([]byte(payload), &m); err != nil {
			return errors.WithStack(err)
		}
		return h(ctx, &m)
	})
}
