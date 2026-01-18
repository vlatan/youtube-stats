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

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("no valid port; ", err)
	}

	fmt.Printf("Server running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
