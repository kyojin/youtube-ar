package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/yansal/youtube-ar/broker"
	"github.com/yansal/youtube-ar/broker/redis"
	"github.com/yansal/youtube-ar/downloader"
	"github.com/yansal/youtube-ar/log"
	"github.com/yansal/youtube-ar/manager"
	"github.com/yansal/youtube-ar/model"
	"github.com/yansal/youtube-ar/payload"
	"github.com/yansal/youtube-ar/server"
	"github.com/yansal/youtube-ar/storage"
	"github.com/yansal/youtube-ar/store"
	"github.com/yansal/youtube-ar/store/db"
	"github.com/yansal/youtube-ar/worker"
	"github.com/yansal/youtube-ar/worker/handler"
	"github.com/yansal/youtube-ar/youtube"
)

// Cmd is the cmd functional type.
type Cmd func(ctx context.Context, args []string) error

// CreateURL is the create-url cmd.
func CreateURL(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("create-url", flag.ExitOnError)
	var url string
	fs.StringVar(&url, "url", "", "url to create")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if url == "" {
		return errors.New("url is required")
	}

	log := log.New()
	redis, err := redis.New(log)
	if err != nil {
		return err
	}
	broker := broker.New(redis, log)
	db, err := db.New(log)
	if err != nil {
		return err
	}
	store := store.New(db)
	m := manager.New(broker, nil, nil, store, nil)

	p := payload.URL{URL: url}
	if err := p.Validate(); err != nil {
		return err
	}

	if _, err := m.CreateURL(ctx, p); err != nil {
		return err
	}
	return nil
}

// CreateURLsFromYoutubePlaylist is the create-url-from-playlist cmd.
func CreateURLsFromYoutubePlaylist(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("create-urls-from-youtube-playlist", flag.ExitOnError)
	var playlist string
	fs.StringVar(&playlist, "playlist", "", "youtube playlist")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if playlist == "" {
		return errors.New("playlist is required")
	}

	log := log.New()
	redis, err := redis.New(log)
	if err != nil {
		return err
	}
	broker := broker.New(redis, log)
	db, err := db.New(log)
	if err != nil {
		return err
	}
	store := store.New(db)
	m := manager.New(broker, nil, nil, store, youtube.New(log))

	return m.CreateURLsFromYoutube(ctx, playlist)
}

// DownloadURL is the download-url cmd.
func DownloadURL(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("download-url", flag.ExitOnError)
	var url string
	fs.StringVar(&url, "url", "", "url to download")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if url == "" {
		return errors.New("url is required")
	}

	d := downloader.New()
	stream := d.Download(ctx, url)
	for event := range stream {
		switch event.Type {
		case downloader.Log:
			fmt.Println(event.Log)
		case downloader.Failure:
			return event.Err
		case downloader.Success:
			fmt.Printf("downloaded url to %s\n", event.Path)
		}
	}
	return nil
}

// ListLogs is the list-logs cmd.
func ListLogs(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("list-logs", flag.ExitOnError)
	var urlID, cursor, limit int64
	fs.Int64Var(&urlID, "url-id", 0, "url-id")
	fs.Int64Var(&cursor, "cursor", 0, "cursor")
	fs.Int64Var(&limit, "limit", 10, "limit")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if urlID == 0 {
		return errors.New("url-id is required")
	}

	log := log.New()
	db, err := db.New(log)
	if err != nil {
		return err
	}
	store := store.New(db)
	m := manager.New(nil, nil, nil, store, nil)

	logs, err := m.ListLogs(ctx, urlID, &model.Page{Cursor: cursor, Limit: limit})
	if err != nil {
		return err
	}
	for i := range logs {
		fmt.Printf("%+v\n", logs[i])
	}
	return nil
}

// ListURLs is the list-urls cmd.
func ListURLs(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("list-urls", flag.ExitOnError)
	var cursor, limit int64
	fs.Int64Var(&cursor, "cursor", 0, "cursor")
	fs.Int64Var(&limit, "limit", 10, "limit")
	if err := fs.Parse(args); err != nil {
		return err
	}

	log := log.New()
	db, err := db.New(log)
	if err != nil {
		return err
	}
	store := store.New(db)
	m := manager.New(nil, nil, nil, store, nil)

	urls, err := m.ListURLs(ctx, &model.Page{Cursor: cursor, Limit: limit})
	if err != nil {
		return err
	}
	for i := range urls {
		fmt.Printf("%+v\n", urls[i])
	}
	return nil
}

// Server is the server cmd.
func Server(ctx context.Context, args []string) error {
	log := log.New()
	redis, err := redis.New(log)
	if err != nil {
		return err
	}
	broker := broker.New(redis, log)
	db, err := db.New(log)
	if err != nil {
		return err
	}
	store := store.New(db)
	m := manager.New(broker, nil, nil, store, nil)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return http.ListenAndServe(":"+port, server.Handler(m, log))
}

// Worker is the worker cmd.
func Worker(ctx context.Context, args []string) error {
	log := log.New()
	redis, err := redis.New(log)
	if err != nil {
		return err
	}
	b := broker.New(redis, log)

	storage, err := storage.New()
	if err != nil {
		return err
	}
	db, err := db.New(log)
	if err != nil {
		return err
	}
	store := store.New(db)
	m := manager.New(nil, downloader.New(), storage, store, nil)

	w := worker.New(b, map[string]broker.Handler{
		"url-created": handler.URLCreated(m),
	})
	return w.Listen(ctx)
}
