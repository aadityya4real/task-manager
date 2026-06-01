package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GetSecretKey retrieves JWT secret from environment or uses default for development
func GetSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default for development only - should be set in production
		secret = "mysecretkey"
	}
	return []byte(secret)
}

func GenerateToken(userID int, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	})

	return token.SignedString(GetSecretKey())
}

func ValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return GetSecretKey(), nil
	})
}
