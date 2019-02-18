package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type handler struct {
	manager *manager
	tmpl    *template.Template
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.serveHTTP(w, r)
	if err == nil {
		return
	}

	herr, ok := err.(httpError)
	if !ok {
		herr = httpError{err: err, code: http.StatusInternalServerError}
	}
	http.Error(w, fmt.Sprintf("%+v", herr.err), herr.code)
}

func (h *handler) serveHTTP(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return h.serveGET(w, r)
	case http.MethodPost:
		return h.servePOST(w, r)
	}
	return httpError{code: http.StatusMethodNotAllowed}
}

func (h *handler) serveGET(w http.ResponseWriter, r *http.Request) error {
	data, err := h.manager.list(r.Context())
	if err != nil {
		return err
	}
	if err := h.tmpl.Execute(w, data); err != nil {
		log.Print(err)
	}
	return nil
}

func (h *handler) servePOST(w http.ResponseWriter, r *http.Request) error {
	payload, err := bind(r)
	if err != nil {
		return err
	}

	if err := h.manager.create(r.Context(), payload); err != nil {
		return err
	}

	http.Redirect(w, r, "", http.StatusFound)
	return nil
}

type httpError struct {
	err  error
	code int
}

func (e httpError) Error() string { return e.err.Error() }
