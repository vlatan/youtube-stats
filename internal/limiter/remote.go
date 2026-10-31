package limiter

import (
	"net/http"
	"strings"
)

// remoteAddr returns the client IP
func remoteAddr(r *http.Request) string {

	// Prioritize CF-Connecting-IP as recommended by Cloudflare
	addr := r.Header.Get("CF-Connecting-IP")

	// Fallback to True-Client-IP
	if addr == "" {
		addr = r.Header.Get("True-Client-IP")
	}

	// Fallback to X-Forwarded-For
	if addr == "" {
		xForwardedFor := r.Header.Get("X-Forwarded-For")
		parts := strings.Split(xForwardedFor, ",")
		addr = strings.TrimSpace(parts[0])
	}

	// Fallback to RemoteAddr
	if addr == "" {
		addr = r.RemoteAddr
	}

	return addr
}
