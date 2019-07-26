package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/yansal/youtube-ar/manager"
	"github.com/yansal/youtube-ar/model"
	"github.com/yansal/youtube-ar/payload"
	"github.com/yansal/youtube-ar/resource"
)

type apiHandler struct {
	manager manager.Manager
}

func (h *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := h.serveHTTP(r)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.code)
		w.Write(resp.body)
		return
	}
	herr, ok := err.(httpError)
	if !ok {
		herr = httpError{
			code: http.StatusInternalServerError,
			err:  err,
		}
	}
	w.WriteHeader(herr.code)
	w.Write([]byte(herr.Error()))
}

type response struct {
	body []byte
	code int
}

type httpError struct {
	code int
	err  error
}

func (e httpError) Error() string {
	if e.err == nil {
		return http.StatusText(e.code)
	}
	return e.err.Error()
}

func (h *apiHandler) serveHTTP(r *http.Request) (*response, error) {
	if !strings.HasPrefix(r.URL.Path, "/api/urls") {
		return nil, httpError{code: http.StatusNotFound}
	}

	if r.URL.Path == "/api/urls" {
		switch r.Method {
		case http.MethodGet:
			return h.listURLs(r)
		case http.MethodPost:
			return h.createURL(r)
		default:
			return nil, httpError{code: http.StatusMethodNotAllowed}
		}
	}

	split := strings.Split(r.URL.Path, "/")
	if len(split) != 5 || split[4] != "logs" {
		return nil, httpError{code: http.StatusNotFound}
	}
	id, err := strconv.ParseInt(split[3], 0, 0)
	if err != nil {
		return nil, httpError{code: http.StatusNotFound}
	}
	if r.Method != http.MethodGet {
		return nil, httpError{code: http.StatusMethodNotAllowed}
	}
	return h.listLogs(r, id)
}

func (h *apiHandler) listURLs(r *http.Request) (*response, error) {
	ctx := r.Context()
	page, err := parseQuery(r.URL.Query())
	if err != nil {
		return nil, httpError{
			err:  err,
			code: http.StatusBadRequest,
		}
	}

	urls, err := h.manager.ListURLs(ctx, page)
	if err != nil {
		return nil, err
	}
	resource := resource.NewURLs(urls)
	b, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}
	return &response{body: b, code: http.StatusOK}, nil
}

func (h *apiHandler) createURL(r *http.Request) (*response, error) {
	ctx := r.Context()
	var payload payload.URL
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, httpError{
			err:  err,
			code: http.StatusBadRequest,
		}
	}
	if err := payload.Validate(); err != nil {
		return nil, httpError{
			err:  err,
			code: http.StatusBadRequest,
		}
	}

	url, err := h.manager.CreateURL(ctx, payload)
	if err != nil {
		return nil, err
	}
	resource := resource.NewURL(url)
	b, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}
	return &response{body: b, code: http.StatusCreated}, nil
}

func (h *apiHandler) listLogs(r *http.Request, id int64) (*response, error) {
	ctx := r.Context()
	page, err := parseQuery(r.URL.Query())
	if err != nil {
		return nil, httpError{
			err:  err,
			code: http.StatusBadRequest,
		}
	}

	logs, err := h.manager.ListLogs(ctx, id, page)
	if err != nil {
		return nil, err
	}
	resource := resource.NewLogs(logs)
	b, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}
	return &response{body: b, code: http.StatusOK}, nil
}

func parseQuery(v url.Values) (*model.Page, error) {
	page := model.Page{Limit: 10}
	limit := v.Get("limit")
	if limit != "" {
		var err error
		page.Limit, err = strconv.ParseInt(limit, 0, 0)
		if err != nil {
			return nil, httpError{
				err:  err,
				code: http.StatusBadRequest,
			}
		}
	}
	cursor := v.Get("cursor")
	if cursor != "" {
		var err error
		page.Cursor, err = strconv.ParseInt(cursor, 0, 0)
		if err != nil {
			return nil, httpError{
				err:  err,
				code: http.StatusBadRequest,
			}
		}
	}
	return &page, nil
}
