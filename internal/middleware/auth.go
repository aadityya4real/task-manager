package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/aadityya4real/task-manager/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Extract token
		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		// Validate token using centralized function
		token, err := utils.ValidateToken(tokenStr)

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims safely
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Safe extraction of user_id
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)

		// Extract username
		username, ok := claims["username"].(string)
		if !ok {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Store in context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "username", username)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
