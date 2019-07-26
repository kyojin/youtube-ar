package manager

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/yansal/youtube-ar/broker"
	"github.com/yansal/youtube-ar/downloader"
	"github.com/yansal/youtube-ar/event"
	"github.com/yansal/youtube-ar/model"
	"github.com/yansal/youtube-ar/payload"
	"github.com/yansal/youtube-ar/storage"
	"github.com/yansal/youtube-ar/store"
	"github.com/yansal/youtube-ar/youtube"
)

// Manager is the manager interface.
type Manager interface {
	CreateURL(context.Context, payload.URL) (*model.URL, error)

	CreateURLsFromYoutube(context.Context, string) error

	ProcessURL(context.Context, event.URL) error

	ListURLs(context.Context, *model.Page) ([]model.URL, error)
	ListLogs(context.Context, int64, *model.Page) ([]model.Log, error)
}

// New returns a new manager.
func New(
	broker broker.Broker,
	downloader downloader.Downloader,
	storage storage.Storage,
	store store.Store,
	ytc youtube.Client,
) Manager {
	return &manager{broker: broker, downloader: downloader, storage: storage, store: store, ytc: ytc}
}

type manager struct {
	broker     broker.Broker
	downloader downloader.Downloader
	storage    storage.Storage
	store      store.Store
	ytc        youtube.Client
}

func (m *manager) CreateURL(ctx context.Context, p payload.URL) (*model.URL, error) {
	url := &model.URL{URL: p.URL}
	if err := m.store.CreateURL(ctx, url); err != nil {
		return nil, err
	}

	e := &event.URL{ID: url.ID}
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	if err := m.broker.Send(ctx, "url-created", string(b)); err != nil {
		return nil, err
	}
	return url, nil
}

func (m *manager) CreateURLsFromYoutube(ctx context.Context, playlistID string) error {
	videos, err := m.ytc.GetVideosFromPlaylist(ctx, playlistID)
	if err != nil {
		return err
	}

	for i := range videos {
		youtubeID := videos[i].ID
		v := &model.YoutubeVideo{YoutubeID: youtubeID}
		if err := m.store.CreateYoutubeVideo(ctx, v); err == sql.ErrNoRows {
			continue
		} else if err != nil {
			return err
		}

		p := payload.URL{URL: "https://www.youtube.com/watch?v=" + youtubeID}
		if _, err := m.CreateURL(ctx, p); err != nil {
			return err
		}
	}

	return nil
}

func (m *manager) ProcessURL(ctx context.Context, e event.URL) error {
	// TODO: process in db transaction?

	url := &model.URL{ID: e.ID, Status: "processing"}
	if err := m.store.LockURL(ctx, url); err != nil {
		return err
	}

	file, err := m.processURL(ctx, url)
	if err != nil {
		url.Error = sql.NullString{Valid: true, String: err.Error()}
		url.Status = "failure"
		if err := m.store.UnlockURL(ctx, url); err != nil {
			// TODO: log err
		}
		return err
	}

	url.File = sql.NullString{Valid: true, String: file}
	url.Status = "success"
	return m.store.UnlockURL(ctx, url)
}

func (m *manager) processURL(ctx context.Context, url *model.URL) (string, error) {
	var (
		path string
		err  error
	)
	stream := m.downloader.Download(ctx, url.URL)
	for event := range stream {
		switch event.Type {
		case downloader.Log:
			if err := m.store.AppendLog(ctx, &model.Log{
				URLID: url.ID,
				Log:   event.Log,
			}); err != nil {
				// TODO: log err
			}
		case downloader.Failure:
			err = event.Err
		case downloader.Success:
			path = event.Path
		}
	}
	defer os.Remove(path)
	if err != nil {
		return "", err
	}

	uploaded, err := m.storage.Upload(ctx, path)
	if err != nil {
		return "", err
	}
	return filepath.Base(uploaded), nil
}

func (m *manager) ListURLs(ctx context.Context, page *model.Page) ([]model.URL, error) {
	return m.store.ListURLs(ctx, page)
}

func (m *manager) ListLogs(ctx context.Context, urlID int64, page *model.Page) ([]model.Log, error) {
	return m.store.ListLogs(ctx, urlID, page)
}
