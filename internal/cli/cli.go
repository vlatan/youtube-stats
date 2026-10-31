package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/vlatan/youtube-stats/internal/config"
	"github.com/vlatan/youtube-stats/internal/youtube"
)

const borderLen = 19

type Service struct {
	youtube *youtube.Service
	url     string
}

func New(ctx context.Context) (*Service, error) {

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	ytsvc, err := youtube.New(ctx, cfg.ApiKey)
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
