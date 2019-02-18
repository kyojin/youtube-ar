package main

import (
	"errors"
	"net/http"
)

type payload struct{ url string }

func bind(r *http.Request) (payload, error) {
	p := payload{url: r.FormValue("url")}
	if p.url == "" {
		return p, httpError{err: errors.New("url is required"), code: http.StatusBadRequest}
	}
	// TODO: validate url?
	return p, nil
}
