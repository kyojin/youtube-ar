package resource

import (
	"time"

	"github.com/yansal/youtube-ar/model"
)

// URL is the url resource.
type URL struct {
	ID        int64     `json:"id,omitempty"`
	URL       string    `json:"url,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Status    string    `json:"status,omitempty"`
	Error     string    `json:"error,omitempty"`
	File      string    `json:"file,omitempty"`
}

// NewURL returns a new URL.
func NewURL(url *model.URL) *URL {
	resource := URL{
		ID:        url.ID,
		URL:       url.URL,
		CreatedAt: url.CreatedAt,
		UpdatedAt: url.UpdatedAt,
		Status:    url.Status,
	}
	if url.Error.Valid {
		resource.Error = url.Error.String
	}
	if url.File.Valid {
		resource.File = url.File.String
	}
	return &resource
}

// NewURLs returns a new URL list.
func NewURLs(urls []model.URL) []URL {
	var resources []URL
	for i := range urls {
		resource := NewURL(&urls[i])
		resources = append(resources, *resource)
	}
	return resources
}

// Log is the log resource.
type Log struct {
	ID        int64     `json:"id,omitempty"`
	Log       string    `json:"log,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// NewLog returns a new Log.
func NewLog(log *model.Log) *Log {
	resource := Log{
		ID:        log.ID,
		Log:       log.Log,
		CreatedAt: log.CreatedAt,
	}
	return &resource
}

// NewLogs returns a new Log list.
func NewLogs(logs []model.Log) []Log {
	var resources []Log
	for i := range logs {
		resource := NewLog(&logs[i])
		resources = append(resources, *resource)
	}
	return resources
}
