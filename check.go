package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const borderLen = 19

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

	var video *youtube.Video = getVideo(apiKey, args[1])
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
		log.Fatal("Probably no such video with this ID: ", videoID)
	}

	// https://pkg.go.dev/google.golang.org/api@v0.201.0/youtube/v3#Video
	return videoList[0]
}

func printVideoInfo(video *youtube.Video) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug)
	yellow := color.New(color.FgYellow, color.Bold)
	border := strings.Repeat("-", borderLen)
	border = fmt.Sprint(border, "\t", border)
	tab := yellow.Sprint("\t ")

	stats := []string{
		// fmt.Sprint("Title", tab, video.Snippet.Title),
		fmt.Sprint("Privacy Status", tab, video.Status.PrivacyStatus),
		fmt.Sprint("Age Restricted", tab, ageRestriction(video)),
		fmt.Sprint("Embeddable", tab, video.Status.Embeddable),
		fmt.Sprint("Region Restricted", tab, regionRestriction(video)),
		// fmt.Sprint("Default Labguage", tab, video.Snippet.DefaultLanguage),
		fmt.Sprint("Live Broadcast", tab, video.Snippet.LiveBroadcastContent),
		fmt.Sprint("Duration", tab, video.ContentDetails.Duration),
	}

	yellow.Fprintln(w, border)
	for _, stat := range stats {
		fmt.Fprintln(w, stat)
		yellow.Fprintln(w, border)
	}
	w.Flush()
}

func regionRestriction(video *youtube.Video) []string {
	restriction := video.ContentDetails.RegionRestriction
	if restriction == nil {
		return []string{}
	}
	return restriction.Blocked
}

func ageRestriction(video *youtube.Video) bool {
	rating := video.ContentDetails.ContentRating.YtRating
	return rating == "ytAgeRestricted"
}
