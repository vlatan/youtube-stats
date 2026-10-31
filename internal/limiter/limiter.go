package limiter

import (
	"net/http"
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
type Service struct {
	ips             map[string]*client
	mu              *sync.RWMutex
	rateLimit       rate.Limit
	burstSize       int
	cleanupInterval time.Duration
	maxIdleTime     time.Duration
}

// Produce new IP rate limiter with a given rate limit, burst size, cleanup interval and max idle time
func New(r rate.Limit, b int, cleanupInterval, maxIdleTime time.Duration) *Service {
	limiter := &Service{
		ips:             make(map[string]*client),
		mu:              new(sync.RWMutex),
		rateLimit:       r,
		burstSize:       b,
		cleanupInterval: cleanupInterval,
		maxIdleTime:     maxIdleTime,
	}

	// Start cleanup routine in the background
	go limiter.cleanupLoop()

	return limiter
}

// Middleware to limit the number of requests from a single IP address
func (s *Service) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		remoteAddr := remoteAddr(r)
		ipLimiter := s.get(remoteAddr)
		if !ipLimiter.Allow() {
			code := http.StatusTooManyRequests
			http.Error(w, http.StatusText(code), code)
			return
		}
		next(w, r)
	}
}

// Returns the rate limiter for the provided IP address if it exists.
// Otherwise ceates a new rate limiter and adds it to the ips map.
func (s *Service) get(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client, exists := s.ips[ip]; exists {
		client.lastSeen = time.Now()
		return client.limiter
	}

	limiter := rate.NewLimiter(s.rateLimit, s.burstSize)
	s.ips[ip] = &client{limiter: limiter, lastSeen: time.Now()}
	return limiter
}

// Cleanup old entries periodically
func (s *Service) cleanupLoop() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanup()
	}
}

// Remove clients that haven't been seen for a while
func (s *Service) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for ip, client := range s.ips {
		if time.Since(client.lastSeen) > s.maxIdleTime {
			delete(s.ips, ip)
		}
	}
}
