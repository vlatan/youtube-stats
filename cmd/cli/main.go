package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	common "github.com/vlatan/youtube-stats/internal"
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

	video, err := common.GetVideo(apiKey, args[1])
	if err != nil {
		log.Fatal(err)
	}
	printVideoInfo(video)
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
		fmt.Sprint("Age Restricted", tab, common.AgeRestriction(video)),
		fmt.Sprint("Embeddable", tab, video.Status.Embeddable),
		fmt.Sprint("Region Restricted", tab, common.RegionRestriction(video)),
		fmt.Sprint("Default Labguage", tab, video.Snippet.DefaultLanguage),
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
