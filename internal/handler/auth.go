package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aadityya4real/task-manager/internal/middleware"
	"github.com/aadityya4real/task-manager/internal/storage"
	"github.com/aadityya4real/task-manager/internal/types"
	"github.com/aadityya4real/task-manager/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

// SignupHandler handles user registration
func SignupHandler(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var u types.User

		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Sanitize input
		u.Username = middleware.SanitizeString(u.Username)
		u.Password = middleware.SanitizeString(u.Password)

		// Validate input
		if u.Username == "" || u.Password == "" {
			http.Error(w, "Username and password required", http.StatusBadRequest)
			return
		}

		if len(u.Password) < 8 {
			http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			http.Error(w, "Error processing request", http.StatusInternalServerError)
			return
		}

		u.Password = string(hashedPassword)

		// Save user
		id, err := store.CreateUser(u)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			log.Printf("Error creating user: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		u.ID = int(id)
		u.Password = "" // Don't send password back

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User created successfully",
			"user":    u,
		})
	}
}

// LoginHandler handles user authentication
func LoginHandler(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var u types.User

		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Sanitize input
		u.Username = middleware.SanitizeString(u.Username)
		u.Password = middleware.SanitizeString(u.Password)

		if u.Username == "" || u.Password == "" {
			http.Error(w, "Username and password required", http.StatusBadRequest)
			return
		}

		dbUser, err := store.GetUser(u.Username)
		if err != nil {
			// Generic error message to prevent username enumeration
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(u.Password))
		if err != nil {
			// Generic error message to prevent username enumeration
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateToken(dbUser.ID, dbUser.Username)
		if err != nil {
			log.Printf("Error generating token: %v", err)
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}
