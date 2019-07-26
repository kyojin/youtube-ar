package worker

import (
	"context"

	"github.com/yansal/youtube-ar/broker"
	"golang.org/x/sync/errgroup"
)

// Worker is the worker interface.
type Worker interface {
	Listen(context.Context) error
}

// New returns a new worker.
func New(b broker.Broker, h map[string]broker.Handler) Worker {
	return &worker{broker: b, handlers: h}
}

type worker struct {
	broker   broker.Broker
	handlers map[string]broker.Handler
}

func (w *worker) Listen(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	for queue, handler := range w.handlers {
		queue := queue
		handler := handler
		g.Go(func() error {
			for {
				if err := w.broker.Receive(ctx, queue, handler); err != nil {
					return err
				}
			}
		})
	}
	return g.Wait()
}
