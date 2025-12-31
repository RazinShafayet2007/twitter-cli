package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/RazinShafayet2007/twitter-cli/internal/models"
	"github.com/oklog/ulid/v2"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

// Create creates a new user
func (s *UserStore) Create(username string) (*models.User, error) {
	// Generate ULID for ID
	id := ulid.Make().String()
	now := time.Now().Unix()

	query := `
		INSERT INTO users (id, username, created_at)
		VALUES (?, ?, ?)
	`

	_, err := s.db.Exec(query, id, username, now)
	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "UNIQUE constraint failed: users.username" {
			return nil, errors.New("username already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &models.User{
		ID:        id,
		Username:  username,
		CreatedAt: now,
	}, nil
}

// GetByUsername retrieves a user by username
func (s *UserStore) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, created_at
		FROM users
		WHERE username = ?
	`

	var user models.User
	err := s.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByID retrieves a user by ID
func (s *UserStore) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, username, created_at
		FROM users
		WHERE id = ?
	`

	var user models.User
	err := s.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
