package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aadityya4real/Task-manager/internal/storage"
	"github.com/aadityya4real/Task-manager/internal/types"
	"github.com/aadityya4real/Task-manager/internal/utils"
)

func AuthHandler(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 🔹 SIGNUP
		if r.Method == "POST" && r.URL.Path == "/signup" {
			var u types.User

			err := json.NewDecoder(r.Body).Decode(&u)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			if u.Username == "" || u.Password == "" {
				http.Error(w, "Username and password required", http.StatusBadRequest)
				return
			}

			id, err := store.CreateUser(u)
			if err != nil {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}

			u.ID = int(id)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(u)
			return
		}

		// 🔹 LOGIN
		if r.Method == "POST" && r.URL.Path == "/login" {
			var u types.User

			err := json.NewDecoder(r.Body).Decode(&u)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			if u.Username == "" || u.Password == "" {
				http.Error(w, "Username and password required", http.StatusBadRequest)
				return
			}

			// 🔹 Fetch user from DB
			dbUser, err := store.GetUser(u.Username)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			// 🔹 Check password
			if dbUser.Password != u.Password {
				http.Error(w, "Invalid password", http.StatusUnauthorized)
				return
			}

			// 🔥 Generate JWT Token
			token, err := utils.GenerateToken(dbUser.Username)
			if err != nil {
				http.Error(w, "Failed to generate token", http.StatusInternalServerError)
				return
			}

			// 🔹 Send token
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"token": token,
			})

			return
		}

		// ❌ Method not allowed
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
