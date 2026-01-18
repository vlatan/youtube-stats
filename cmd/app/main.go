package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /static/", staticHandler)
	mux.HandleFunc("GET /{$}", getHandler)

	postHandler := applyMiddlewares(postHandler, limitMiddleware)
	mux.HandleFunc("POST /{$}", postHandler)

	// Simple health check
	mux.HandleFunc("GET /healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Robots-Tag", "noindex")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("Failed to write response on '%s'; %v", r.URL.Path, err)
		}
	})

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("no valid port; ", err)
	}

	fmt.Printf("Server running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
