package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

// Generate a secure random CSRF token
func generateCSRFToken() (string, error) {
	b := make([]byte, 32) // 32 bytes is a good size for a token
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Get CSRF token from the request's cookie.
// If cookie not valid generates a token and sets it in the response's HttpOnly cookie
func getCSRFToken(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(csrfCookieName)
	if err == nil {
		return cookie.Value, nil
	}

	token, err := generateCSRFToken()
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,                           // Javascript will not be able to access the cookie
		Secure:   true,                           // Set to true in production (HTTPS)
		SameSite: http.SameSiteLaxMode,           // Recommended: Lax or Strict
		Expires:  time.Now().Add(24 * time.Hour), // Token expiration (optional, but good practice)
	})

	return token, nil

}

func validateCSRFToken(r *http.Request) bool {
	cookie, err := r.Cookie(csrfCookieName)
	if err != nil || cookie.Value == "" {
		return false
	}

	return r.FormValue(csrfFieldName) == cookie.Value
}
