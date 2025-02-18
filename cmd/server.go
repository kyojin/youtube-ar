package cmd

import (
	"context"
	"net/http"
	"os"
	"regexp"

	"github.com/yansal/youtube-ar/broker"
	"github.com/yansal/youtube-ar/broker/redis"
	"github.com/yansal/youtube-ar/log"
	"github.com/yansal/youtube-ar/manager"
	"github.com/yansal/youtube-ar/server"
	"github.com/yansal/youtube-ar/server/handler"
	"github.com/yansal/youtube-ar/server/middleware"
	"github.com/yansal/youtube-ar/store"
	"github.com/yansal/youtube-ar/store/db"
)

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
	manager := manager.NewServer(broker, store)

	mux := server.NewMux()
	mux.HandleFunc(http.MethodGet, regexp.MustCompile(`^/api/urls$`), handler.ListURLs(manager))
	mux.HandleFunc(http.MethodGet, regexp.MustCompile(`^/api/urls/(\d+)$`), handler.DetailURL(manager))
	mux.HandleFunc(http.MethodPost, regexp.MustCompile(`^/api/urls$`), handler.CreateURL(manager))
	mux.HandleFunc(http.MethodGet, regexp.MustCompile(`^/api/urls/(\d+)/logs$`), handler.ListLogs(manager))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := middleware.Log(mux, log)
	handler = middleware.CORS(mux)
	return http.ListenAndServe(":"+port, handler)
}
