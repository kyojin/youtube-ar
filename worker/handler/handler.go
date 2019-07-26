package handler

import (
	"context"
	"encoding/json"

	"github.com/yansal/youtube-ar/broker"
	"github.com/yansal/youtube-ar/event"
	"github.com/yansal/youtube-ar/manager"
)

// URLCreated is the url-created handler.
func URLCreated(m manager.Manager) broker.Handler {
	return func(ctx context.Context, payload string) error {
		var e event.URL
		if err := json.Unmarshal([]byte(payload), &e); err != nil {
			return err
		}

		return m.ProcessURL(ctx, e)
	}
}
