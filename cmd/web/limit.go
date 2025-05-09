package main

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Client that holds the rate limiter and the last seen time
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter
type ipRateLimiter struct {
	ips             map[string]*client
	mu              *sync.RWMutex
	r               rate.Limit
	b               int
	cleanupInterval time.Duration
	maxIdleTime     time.Duration
}

// Returns the rate limiter for the provided IP address if it exists.
// Otherwise ceates a new rate limiter and adds it to the ips map.
func (i *ipRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	if client, exists := i.ips[ip]; exists {
		client.lastSeen = time.Now()
		return client.limiter
	}

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = &client{limiter: limiter, lastSeen: time.Now()}
	return limiter
}

// Cleanup old entries periodically
func (i *ipRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(i.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		i.cleanup()
	}
}

// Remove clients that haven't been seen for a while
func (i *ipRateLimiter) cleanup() {
	i.mu.Lock()
	defer i.mu.Unlock()

	for ip, client := range i.ips {
		if time.Since(client.lastSeen) > i.maxIdleTime {
			delete(i.ips, ip)
		}
	}
}

// Produce new IP rate limiter with a given rate, burst size, cleanup interval and max idle time
func newIPRateLimiter(r rate.Limit, b int, cleanupInterval, maxIdleTime time.Duration) *ipRateLimiter {
	limiter := &ipRateLimiter{
		ips:             make(map[string]*client),
		mu:              &sync.RWMutex{},
		r:               r,
		b:               b,
		cleanupInterval: cleanupInterval,
		maxIdleTime:     maxIdleTime,
	}

	// Start cleanup routine
	go limiter.cleanupLoop()

	return limiter
}
