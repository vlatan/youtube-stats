package app

import (
	"log"
	"net/http"
)

func (a *App) Handler() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /static/", a.Static)
	mux.HandleFunc("GET /{$}", a.Get)
	mux.HandleFunc("POST /{$}", a.limiter.Limit(a.Post))

	// Simple health check
	mux.HandleFunc("GET /healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Robots-Tag", "noindex")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("Failed to write response on '%s'; %v", r.URL.Path, err)
		}
	})

	cop := http.NewCrossOriginProtection()
	return cop.Handler(mux)
}
