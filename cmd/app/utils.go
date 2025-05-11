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
