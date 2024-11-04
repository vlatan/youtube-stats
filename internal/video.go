package common

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const borderLen = 19

// Get video object using the youtybe Golang package.
// https://pkg.go.dev/google.golang.org/api@v0.201.0/youtube/v3
func GetVideo(apiKey string, videoID string) (*youtube.Video, error) {

	var ctx context.Context = context.Background()
	var co option.ClientOption = option.WithAPIKey(apiKey)

	youtubeService, err := youtube.NewService(ctx, co)
	if err != nil {
		msg := fmt.Sprint("Unable to create a YouTube service.", err)
		return nil, errors.New(msg)
	}

	part := []string{"status", "snippet", "contentDetails"}
	response, err := youtubeService.Videos.List(part).Id(videoID).Do()
	if err != nil {
		msg := fmt.Sprint("Unable to get a response from YouTube.", err)
		return nil, errors.New(msg)
	}

	var videoList []*youtube.Video = response.Items
	if len(videoList) == 0 {
		msg := fmt.Sprint("Probably no such video with this ID: ", videoID)
		return nil, errors.New(msg)
	}

	return videoList[0], nil
}

func PrintVideoInfo(video *youtube.Video) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug)
	yellow := color.New(color.FgYellow, color.Bold)
	border := strings.Repeat("-", borderLen)
	border = fmt.Sprint(border, "\t", border)
	tab := yellow.Sprint("\t ")

	stats := []string{
		// fmt.Sprint("Title", tab, video.Snippet.Title),
		fmt.Sprint("Privacy Status", tab, video.Status.PrivacyStatus),
		fmt.Sprint("Age Restricted", tab, AgeRestriction(video)),
		fmt.Sprint("Embeddable", tab, video.Status.Embeddable),
		fmt.Sprint("Region Restricted", tab, RegionRestriction(video)),
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

func RegionRestriction(video *youtube.Video) []string {
	restriction := video.ContentDetails.RegionRestriction
	if restriction == nil {
		return []string{}
	}
	return restriction.Blocked
}

func AgeRestriction(video *youtube.Video) bool {
	rating := video.ContentDetails.ContentRating.YtRating
	return rating == "ytAgeRestricted"
}
