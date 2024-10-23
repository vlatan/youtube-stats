package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the .env file.", err)
	}

	var args []string = os.Args
	if len(args) <= 1 {
		log.Fatal("Please provide a YouTube video ID")
	}

	var apiKey string = os.Getenv("YOUTUBE_API_KEY")
	var videoID string = strings.TrimSpace(args[1])
	video := getVideo(apiKey, videoID)
	printVideoInfo(video)

}

func getVideo(apiKey string, videoID string) *youtube.Video {
	var ctx context.Context = context.Background()
	var co option.ClientOption = option.WithAPIKey(apiKey)
	youtubeService, err := youtube.NewService(ctx, co)

	if err != nil {
		log.Fatal("Unable to create a YouTube service.", err)
	}

	// https://pkg.go.dev/google.golang.org/api@v0.201.0/youtube/v3#VideoListResponse
	part := []string{"status", "snippet", "contentDetails"}
	response, err := youtubeService.Videos.List(part).Id(videoID).Do()
	if err != nil {
		log.Fatal("Unable to get a response from YouTube.", err)
	}

	var videoList []*youtube.Video = response.Items
	if len(videoList) == 0 {
		log.Fatal("Probably no such video:", videoID)
	}

	// https://pkg.go.dev/google.golang.org/api@v0.201.0/youtube/v3#Video
	return videoList[0]
}

func printVideoInfo(video *youtube.Video) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug)

	// fmt.Fprintln(w, "Title\t", video.Snippet.Title)
	fmt.Fprintln(w, "Privacy Status\t", video.Status.PrivacyStatus)
	rating := video.ContentDetails.ContentRating.YtRating
	fmt.Fprintln(w, "Age Restricted\t", rating == "ytAgeRestricted")
	fmt.Fprintln(w, "Embeddable\t", video.Status.Embeddable)

	restriction := video.ContentDetails.RegionRestriction
	switch restriction {
	case nil:
		fmt.Fprintln(w, "Region Restricted \t false")
	default:
		fmt.Fprintln(w, "Region Restricted\t", restriction.Blocked)
	}

	// fmt.Fprintln(w, "Default Labguage\t", video.Snippet.DefaultLanguage)
	fmt.Fprintln(w, "Live Broadcast\t", video.Snippet.LiveBroadcastContent)
	fmt.Fprintln(w, "Duration\t", video.ContentDetails.Duration)
	w.Flush()
}
