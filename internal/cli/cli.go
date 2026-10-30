package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	_ "github.com/joho/godotenv/autoload"
	"github.com/vlatan/youtube-stats/internal/youtube"
)

const borderLen = 19

type Service struct {
	youtube *youtube.Service
	url     string
}

func New(ctx context.Context) (*Service, error) {

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if len(apiKey) == 0 {
		return nil, errors.New("missing youtube api key")
	}

	ytsvc, err := youtube.New(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	var args []string = os.Args
	if len(args) <= 1 {
		return nil, errors.New("please provide a YouTube video URL")
	}

	return &Service{ytsvc, args[1]}, nil
}

func (s *Service) Run(ctx context.Context) error {

	video, err := s.youtube.Video(ctx, s.url)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug)
	yellow := color.New(color.FgYellow, color.Bold)
	border := strings.Repeat("-", borderLen)
	border = fmt.Sprint(border, "\t", border)
	tab := yellow.Sprint("\t ")

	stats := []string{
		// fmt.Sprint("Title", tab, video.Snippet.Title),
		fmt.Sprint("Privacy Status", tab, video.PrivacyStatus),
		fmt.Sprint("Age Restricted", tab, video.AgeRestricted),
		fmt.Sprint("Embeddable", tab, video.Embeddable),
		fmt.Sprint("Region Restricted", tab, video.RegionRestricted),
		fmt.Sprint("Default Labguage", tab, video.DefaultLanguage),
		fmt.Sprint("Live Broadcast", tab, video.LiveBroadcast),
		fmt.Sprint("Duration", tab, video.Duration),
	}

	yellow.Fprintln(w, border)
	for _, stat := range stats {
		fmt.Fprintln(w, stat)
		yellow.Fprintln(w, border)
	}

	w.Flush()
	return nil
}
