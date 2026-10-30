package main

import (
	"log"

	"github.com/vlatan/youtube-stats/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	// Run tne app
	log.Fatal(app.Run())

}
