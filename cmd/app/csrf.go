package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

// generateCSRFToken generates a secure random CSRF token
func generateCSRFToken() (string, error) {
	b := make([]byte, 32) // 32 bytes is a good size for a token
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Set the CSRF token in an HttpOnly cookie
func setCSRFCookie(w http.ResponseWriter, cookieName, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,                           // Javascript will not be able to access the cookie
		Secure:   true,                           // Set to true in production (HTTPS)
		SameSite: http.SameSiteLaxMode,           // Recommended: Lax or Strict
		Expires:  time.Now().Add(24 * time.Hour), // Token expiration (optional, but good practice)
	})
}

// Retrieve the CSRF token from the request's cookie
func getCSRFCookie(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", nil // No cookie found
		}
		return "", err // Other error (e.g., malformed cookie)
	}
	return cookie.Value, nil
}
