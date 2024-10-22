package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.", err)
	}

	var ctx context.Context = context.Background()
	var API_KEY string = os.Getenv("YOUTUBE_API_KEY")
	var co option.ClientOption = option.WithAPIKey(API_KEY)

	youtubeService, err := youtube.NewService(ctx, co)
	if err != nil {
		log.Fatal("Unable to create a YT service.", err)
	}

	var args []string = os.Args
	if len(args) <= 1 {
		log.Fatal("Please provide a YouTube video ID")
	}

	// https://pkg.go.dev/google.golang.org/api@v0.201.0/youtube/v3#VideoListResponse
	video_id := strings.TrimSpace(args[1])
	part := []string{"status", "snippet", "contentDetails"}
	response, err := youtubeService.Videos.List(part).Id(video_id).Do()
	if err != nil {
		log.Fatal("Unable to get a response from YouTube.", err)
	}

	var videoList []*youtube.Video = response.Items
	if len(videoList) == 0 {
		log.Fatal("Probably no such video:", video_id)
	}

	// https://pkg.go.dev/google.golang.org/api@v0.201.0/youtube/v3#Video
	var video *youtube.Video = videoList[0]
	fmt.Println(video.Snippet.Title)
	fmt.Println(video.Status.PrivacyStatus)
	fmt.Println(video.ContentDetails.ContentRating.YtRating)
	fmt.Println(video.Status.Embeddable)
	regionRestriction := video.ContentDetails.RegionRestriction
	if regionRestriction != nil {
		fmt.Println(video.ContentDetails.RegionRestriction.Blocked)
	}
	fmt.Println(video.Snippet.DefaultLanguage)
	fmt.Println(video.Snippet.LiveBroadcastContent)
	fmt.Println(video.ContentDetails.Duration)
}
