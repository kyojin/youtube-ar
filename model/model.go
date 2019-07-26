package model

import (
	"database/sql"
	"time"
)

// URL is the url model.
type URL struct {
	ID        int64
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    string
	Error     sql.NullString
	File      sql.NullString
}

// Log is the log model.
type Log struct {
	ID        int64
	URLID     int64
	Log       string
	CreatedAt time.Time
}

// YoutubeVideo is the youtube video model.
type YoutubeVideo struct {
	ID        int64
	YoutubeID string
	CreatedAt time.Time
}

// Page is the page model.
type Page struct {
	Limit  int64
	Cursor int64
}
