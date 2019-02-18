package main

import (
	"context"

	"github.com/yansal/youtube-ar/broker"
)

type worker struct {
	broker  broker.Broker
	manager *manager
}

func (w *worker) run(ctx context.Context) error {
	return w.broker.Handle(ctx, w.manager.youtubedl)
}
