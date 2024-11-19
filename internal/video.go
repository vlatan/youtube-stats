package common

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Get video object using the YouTube Golang package.
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
