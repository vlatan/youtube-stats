package youtube

import (
	"context"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Service struct {
	youtube *youtube.Service
}

// Create new YouTube service
func New(ctx context.Context, apiKey string) (*Service, error) {

	clientOption := option.WithAPIKey(apiKey)
	youtube, err := youtube.NewService(ctx, clientOption)

	if err != nil {
		return nil, err
	}

	return &Service{youtube}, nil
}
