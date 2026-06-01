package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/aadityya4real/task-manager/internal/types"
)

// ValidateUsername checks if username meets requirements
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	if len(username) > 50 {
		return fmt.Errorf("username must be less than 50 characters")
	}
	return nil
}

// CreateUser inserts a new user with validation
func (s *Store) CreateUser(u types.User) (int64, error) {
	if err := ValidateUsername(u.Username); err != nil {
		return 0, err
	}

	// Check if username already exists
	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", u.Username).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("error checking username: %w", err)
	}
	if exists {
		return 0, fmt.Errorf("username already exists")
	}

	result, err := s.DB.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		u.Username, u.Password,
	)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}
	return result.LastInsertId()
}

// GetUser retrieves user by username for login
func (s *Store) GetUser(username string) (types.User, error) {
	var u types.User
	err := s.DB.QueryRow(
		"SELECT id, username, password FROM users WHERE username = ?",
		username,
	).Scan(&u.ID, &u.Username, &u.Password)

	if err == sql.ErrNoRows {
		return u, fmt.Errorf("user not found")
	}
	if err != nil {
		return u, fmt.Errorf("error fetching user: %w", err)
	}

	return u, nil
}
