package youtube

import (
	"context"
	"errors"
	"net/url"
	"regexp"
	"strings"

	"google.golang.org/api/youtube/v3"
)

type Video struct {
	Id               string   `json:"id"`
	Title            string   `json:"title"`
	PrivacyStatus    string   `json:"privacyStatus"`
	AgeRestricted    bool     `json:"ageRestricted"`
	Embeddable       bool     `json:"embeddable"`
	RegionRestricted []string `json:"regionRestricted"`
	DefaultLanguage  string   `json:"defaultLanguage"`
	LiveBroadcast    string   `json:"liveBroadcast"`
	Duration         string   `json:"duration"`
}

// Validate video ID
var validVideoID = regexp.MustCompile("^([-a-zA-Z0-9_]{11})$")

// Video gets video metadata from YouTube
func (s *Service) Video(ctx context.Context, url string) (*Video, error) {

	videoID, err := extractYouTubeID(url)
	if err != nil {
		return nil, err
	}

	if validVideoID.FindStringSubmatch(videoID) == nil {
		return nil, errors.New("could not validate the video ID")
	}

	part := []string{"status", "snippet", "contentDetails"}
	resp, err := s.youtube.Videos.List(part).Id(videoID).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	var videoList []*youtube.Video = resp.Items
	if len(videoList) == 0 {
		return nil, errors.New("probably no such YouTube video")
	}

	video := Video{
		Id:               resp.Items[0].Id,
		Title:            resp.Items[0].Snippet.Title,
		PrivacyStatus:    resp.Items[0].Status.PrivacyStatus,
		AgeRestricted:    ageRestriction(resp.Items[0]),
		Embeddable:       resp.Items[0].Status.Embeddable,
		RegionRestricted: regionRestriction(resp.Items[0]),
		DefaultLanguage:  resp.Items[0].Snippet.DefaultLanguage,
		LiveBroadcast:    resp.Items[0].Snippet.LiveBroadcastContent,
		Duration:         resp.Items[0].ContentDetails.Duration,
	}

	return &video, nil
}

// Extract YouTube ID from URL
func extractYouTubeID(rawURL string) (string, error) {
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

func regionRestriction(video *youtube.Video) []string {
	restriction := video.ContentDetails.RegionRestriction
	if restriction == nil {
		return nil
	}
	return restriction.Blocked
}

func ageRestriction(video *youtube.Video) bool {
	rating := video.ContentDetails.ContentRating.YtRating
	return rating == "ytAgeRestricted"
}
