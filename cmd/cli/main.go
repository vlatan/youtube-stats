package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vlatan/youtube-stats/internal/cli"
)

func main() {

	ctx := context.Background()
	cli, err := cli.New(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := cli.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
