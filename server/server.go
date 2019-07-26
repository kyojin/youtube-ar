package server

import (
	"net/http"

	"github.com/yansal/youtube-ar/log"
	"github.com/yansal/youtube-ar/manager"
)

// Handler is the server handler.
func Handler(m manager.Manager, log log.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/", &apiHandler{manager: m})
	return logMiddleware(mux, log)
}

func logMiddleware(h http.Handler, l log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w}
		defer func() {
			fields := []log.Field{
				log.Int("code", rw.code),
			}
			msg := r.Method + " " + r.URL.Path
			l.Log(r.Context(), msg, fields...)
		}()

		h.ServeHTTP(rw, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	if rw.code == 0 {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(p)
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.ResponseWriter.WriteHeader(code)
	if rw.code != 0 {
		return
	}
	rw.code = code
}
