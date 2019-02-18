package store

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/yansal/youtube-ar/model"
)

type Store interface {
	List(context.Context) ([]model.URL, error)
	Create(context.Context, *model.URL) error
	UpdateStatus(context.Context, *model.URL, string) error
}

func New(ctx context.Context, opts ...Option) (Store, error) {
	var o Options
	for i := range opts {
		opts[i](&o)
	}
	db, err := sqlx.ConnectContext(ctx, "postgres", o.url)
	return &store{db: db}, errors.WithStack(err)
}

type Option func(*Options)
type Options struct{ url string }

func URL(url string) Option { return func(o *Options) { o.url = url } }

type store struct{ db *sqlx.DB }

func (s *store) List(ctx context.Context) ([]model.URL, error) {
	// TODO: allow to filter, order and paginate
	var dest []model.URL
	q := `select id, url, status, created_at, updated_at from urls order by updated_at desc`
	return dest, s.db.SelectContext(ctx, &dest, q)
}

func (s *store) Create(ctx context.Context, m *model.URL) error {
	q := `insert into urls(url) values ($1) returning id, url, status, created_at, updated_at`
	return s.db.GetContext(ctx, m, q, m.URL)
}

func (s *store) UpdateStatus(ctx context.Context, m *model.URL, status string) error {
	q := `update urls set status = $1 where id = $2 returning id, url, status, created_at, updated_at`
	return s.db.GetContext(ctx, m, q, status, m.ID)
}
