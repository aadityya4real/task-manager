package middleware

import (
	"net/http"
	"strings"
)

// ValidateContentType ensures the request has the correct Content-Type
func ValidateContentType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
}

// SanitizeInput removes potentially dangerous characters from input
func SanitizeString(input string) string {
	// Remove null bytes and trim whitespace
	input = strings.ReplaceAll(input, "\x00", "")
	input = strings.TrimSpace(input)
	return input
}
