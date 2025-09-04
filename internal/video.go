package common

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Validate video ID
var validVideoID = regexp.MustCompile("^([-a-zA-Z0-9_]{11})$")

// Extract YouTube ID from URL
func ExtractYouTubeID(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	if parsedURL.Hostname() == "youtu.be" {
		return parsedURL.Path[1:], nil
	}

	if strings.HasSuffix(parsedURL.Hostname(), "youtube.com") {
		if parsedURL.Path == "/watch" {
			return parsedURL.Query().Get("v"), nil
		} else if parsedURL.Path[:7] == "/embed/" {
			return strings.Split(parsedURL.Path, "/")[2], nil
		}
	}

	return "", errors.New("could not extract the video ID")
}

// Get video object using the YouTube Golang package.
// https://pkg.go.dev/google.golang.org/api@v0.201.0/youtube/v3
func GetVideo(apiKey string, videoID string) (*youtube.Video, error) {

	if validVideoID.FindStringSubmatch(videoID) == nil {
		return nil, errors.New("could not validate the video ID")
	}

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
