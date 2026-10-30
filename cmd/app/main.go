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

	if err := app.Run(); err != nil {
		log.Println(err)
	}
}
