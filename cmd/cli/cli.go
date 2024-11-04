package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	common "github.com/vlatan/youtube-stats/internal"
)

func main() {
	godotenv.Load()
	var apiKey string = os.Getenv("YOUTUBE_API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("Please set YOUTUBE_API_KEY environment variable.")
	}

	var args []string = os.Args
	if len(args) <= 1 {
		log.Fatal("Please provide a YouTube video ID")
	}

	video, err := common.GetVideo(apiKey, args[1])
	if err != nil {
		log.Fatal(err)
	}
	common.PrintVideoInfo(video)
}
