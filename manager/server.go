package manager

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/yansal/youtube-ar/event"
	"github.com/yansal/youtube-ar/model"
	"github.com/yansal/youtube-ar/payload"
)

// Server is the manager used for server features.
type Server struct {
	broker BrokerServer
	store  StoreServer
}

// BrokerServer is the broker interface required by Server.
type BrokerServer interface {
	Send(context.Context, string, string) error
}

// StoreServer is the store interface required by Server.
type StoreServer interface {
	CreateURL(context.Context, *model.URL) error
	ListURLs(context.Context, *model.Page) ([]model.URL, error)
	ListLogs(context.Context, int64, *model.Page) ([]model.Log, error)
}

// NewServer returns a new Server.
func NewServer(broker BrokerServer, store StoreServer) *Server {
	return &Server{broker: broker, store: store}
}

// CreateURL creates an URL.
func (m *Server) CreateURL(ctx context.Context, p payload.URL) (*model.URL, error) {
	url := &model.URL{URL: p.URL}
	if p.Retries != 0 {
		url.Retries = sql.NullInt64{Valid: true, Int64: p.Retries}
	}
	if err := m.store.CreateURL(ctx, url); err != nil {
		return nil, err
	}

	e := &event.URL{ID: url.ID}
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return url, m.broker.Send(ctx, "url-created", string(b))
}

// ListURLs lists urls.
func (m *Server) ListURLs(ctx context.Context, page *model.Page) ([]model.URL, error) {
	return m.store.ListURLs(ctx, page)
}

// ListLogs lists logs.
func (m *Server) ListLogs(ctx context.Context, urlID int64, page *model.Page) ([]model.Log, error) {
	return m.store.ListLogs(ctx, urlID, page)
}
