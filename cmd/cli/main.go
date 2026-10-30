package main

import (
	"context"
	"log"
	"time"

	"github.com/vlatan/youtube-stats/internal/cli"
)

func main() {

	ctx := context.Background()
	cli, err := cli.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := cli.Run(ctx); err != nil {
		log.Println(err)
	}
}
