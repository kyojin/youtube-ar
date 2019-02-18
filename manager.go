package main

import (
	"context"
	"log"

	"github.com/yansal/youtube-ar/broker"
	"github.com/yansal/youtube-ar/model"
	"github.com/yansal/youtube-ar/process"
	"github.com/yansal/youtube-ar/pubsub"
	"github.com/yansal/youtube-ar/store"
)

type manager struct {
	store  store.Store
	broker broker.Broker
	pubsub pubsub.PubSub
}

func (m *manager) list(ctx context.Context) ([]model.URL, error) {
	return m.store.List(ctx)
}

func (m *manager) create(ctx context.Context, p payload) error {
	model := &model.URL{URL: p.url}
	err := m.store.Create(ctx, model)
	if err != nil {
		return err
	}
	return m.broker.Send(ctx, model)
}

func (m *manager) youtubedl(ctx context.Context, model *model.URL) error {
	p, err := process.Start(ctx, "youtube-dl", model.URL)
	if err != nil {
		return err
	}
	for b := range p.Output() {
		if err := m.pubsub.Publish(string(b)); err != nil {
			log.Print(err)
			break
		}
	}
	return p.Wait()
}
