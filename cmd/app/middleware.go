package main

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// Middleware to limit the number of requests from a single IP address
func limitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limiter := limiter.getLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

// Helper function to apply multiple middlewares to a handler function
func applyMiddlewares(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}
