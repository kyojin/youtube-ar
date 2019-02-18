package main

import (
	"context"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/yansal/youtube-ar/broker"
	"github.com/yansal/youtube-ar/pubsub"
	"github.com/yansal/youtube-ar/redis"
	"github.com/yansal/youtube-ar/store"
	"golang.org/x/sync/errgroup"
)

const indexHTML = `<html>
<title>youtube-ar</title>
<form method="POST">
    <input name="url" placeholder="url">
    <button>Send</button>
</form>
{{if .}}
<table border="1">
<tr>
	<th align="center">id</th>
	<th align="center">url</th>
	<th align="center">status</th>
	<th align="center">updated at</th>
</tr>
{{range .}}
<tr valign="top">
	<td align="right">{{.ID}}</td>
	<td align="left">{{.URL}}</td>
	<td align="left">{{.Status}}</td>
	<td align="left">{{.UpdatedAt}}</td>
</tr>
{{end}}
{{end}}
`

func main() {
	var (
		manager  manager
		worker   worker
		handler  handler
		listener net.Listener
	)

	init, ctx := errgroup.WithContext(context.Background())
	init.Go(func() error {
		redisurl := os.Getenv("REDIS_URL")
		if redisurl == "" {
			redisurl = "redis://:6379"
		}

		poolsize, _ := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))

		b, err := broker.New(ctx, redis.URL(redisurl), redis.PoolSize(poolsize))
		manager.broker = b
		worker.broker = b
		return err
	})
	init.Go(func() error {
		redisurl := os.Getenv("REDIS_URL")
		if redisurl == "" {
			redisurl = "redis://:6379"
		}

		poolsize, _ := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))

		p, err := pubsub.New(ctx, redis.URL(redisurl), redis.PoolSize(poolsize))
		manager.pubsub = p
		return err
	})
	init.Go(func() error {
		pgurl := os.Getenv("DATABASE_URL")
		if pgurl == "" {
			pgurl = "dbname=youtube-ar sslmode=disable"
		}

		var err error
		manager.store, err = store.New(ctx, store.URL(pgurl))
		return err
	})
	init.Go(func() error {
		var err error
		handler.tmpl, err = template.New("").Parse(indexHTML)
		return errors.WithStack(err)
	})
	init.Go(func() error {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		var err error
		listener, err = net.Listen("tcp", ":"+port)
		return err
	})
	if err := init.Wait(); err != nil {
		log.Fatalf("%+v", err)
	}

	run, ctx := errgroup.WithContext(context.Background())
	run.Go(func() error {
		handler.manager = &manager
		http.Handle("/", &handler)
		http.Handle("/favicon.ico", http.NotFoundHandler())

		// TODO: stop server when ctx is canceled
		return http.Serve(listener, nil)
	})
	run.Go(func() error {
		worker.manager = &manager
		return worker.run(ctx)
	})
	if err := run.Wait(); err != nil {
		log.Fatalf("%+v", err)

	}
}
