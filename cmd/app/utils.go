package main

import (
	"log"
	"os"
	"strconv"
)

func getIntEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	port, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Invalid value for %s: %s. Using default port %d.", key, value, fallback)
		return fallback
	}

	return port
}

func getStringEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getBoolEnv(key string, fallback bool) bool {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Invalid value for %s: %s. Using default value %t.", key, value, fallback)
		return fallback
	}
	return b
}
