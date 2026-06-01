package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORS middleware with configurable allowed origins
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Get allowed origins from environment or use default for development
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			// Development default - restrict in production!
			allowedOrigins = "http://localhost:3000,http://localhost:8080"
		}

		// Check if origin is allowed
		origins := strings.Split(allowedOrigins, ",")
		allowed := false
		for _, o := range origins {
			if strings.TrimSpace(o) == origin || strings.TrimSpace(o) == "*" {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
