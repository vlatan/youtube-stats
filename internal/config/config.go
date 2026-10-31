package config

import (
	"errors"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Host   string
	Port   string
	ApiKey string `json:"-"`
}

const (
	hostKey = "HOST"
	portKey = "PORT"
	ytKey   = "YOUTUBE_API_KEY"
)

// New creates new config object
func New() (*Config, error) {

	apiKey := os.Getenv(ytKey)
	if apiKey == "" {
		return nil, errors.New("missing youtube api key")
	}

	host := os.Getenv(hostKey)
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv(portKey)
	if port == "" {
		port = "5000"
	}

	return &Config{host, port, apiKey}, nil
}
