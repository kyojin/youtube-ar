package handler

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/yansal/youtube-ar/model"
	"github.com/yansal/youtube-ar/query"
)

func assertf(t *testing.T, ok bool, msg string, args ...interface{}) {
	t.Helper()
	if !ok {
		t.Errorf(msg, args...)
	}
}

type mockManager struct {
	listURLsFunc func(context.Context, *query.URLs) ([]model.URL, error)
}

func (m mockManager) ListURLs(ctx context.Context, q *query.URLs) ([]model.URL, error) {
	return m.listURLsFunc(ctx, q)
}

func TestListURLs(t *testing.T) {
	h := listURLs(mockManager{
		listURLsFunc: func(ctx context.Context, q *query.URLs) ([]model.URL, error) {
			assertf(t, q.Cursor == 0, "expected cursor to be 0, got %d", q.Cursor)
			assertf(t, q.Limit == query.DefaultLimit, "expected limit to be %d, got %d", query.DefaultLimit, q.Limit)
			assertf(t, q.Status == nil, "expected status to be nil, got %v", q.Status)
			return nil, nil
		},
	})

	var u url.URL
	req, err := http.NewRequest("", u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := h(req); err != nil {
		t.Fatal(err)
	}
}

func TestListURLsQuery(t *testing.T) {
	var (
		cursor int64 = 10
		limit  int64 = 20
		status       = []string{"failure", "success"}
	)
	h := listURLs(mockManager{
		listURLsFunc: func(ctx context.Context, q *query.URLs) ([]model.URL, error) {
			assertf(t, q.Cursor == cursor, "expected cursor to be %d, got %d", cursor, q.Cursor)
			assertf(t, q.Limit == limit, "expected limit to be %d, got %d", limit, q.Limit)
			assertf(t, len(q.Status) == len(status), "expected %d status, got %v", len(status), len(q.Status))
			for i := range status {
				assertf(t, q.Status[i] == status[i], "expected status %v, got %v", status[i], q.Status[i])

			}
			return nil, nil
		},
	})

	v := url.Values{
		"status": status,
		"cursor": []string{strconv.FormatInt(cursor, 10)},
		"limit":  []string{strconv.FormatInt(limit, 10)},
	}
	u := url.URL{RawQuery: v.Encode()}

	req, err := http.NewRequest("", u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := h(req); err != nil {
		t.Fatal(err)
	}
}
