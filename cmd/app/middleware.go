package main

import (
	"net/http"

	"github.com/rs/cors"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// CORS options for the CORS
func newCorsOptions(allowedOrigins []string, debug bool) cors.Options {
	return cors.Options{
		AllowedOrigins:   allowedOrigins,                            // What origins are allowed to access the API
		AllowedMethods:   []string{"GET", "POST"},                   // All methods the API uses
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Custom headers the frontend sends
		AllowCredentials: true,                                      // Important if the frontend sends cookies or Authorization headers (e.g., JWTs in a secure cookie)
		MaxAge:           86400,                                     // Cache preflight requests for 24 hours
		Debug:            debug,                                     // Set to false in production
	}
}

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
