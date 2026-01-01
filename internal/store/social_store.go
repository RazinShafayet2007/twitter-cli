package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/RazinShafayet2007/twitter-cli/internal/models"
)

type SocialStore struct {
	db *sql.DB
}

func NewSocialStore(db *sql.DB) *SocialStore {
	return &SocialStore{db: db}
}

// Follow creates a follow relationship
func (s *SocialStore) Follow(followerID, followeeID string) error {
	// Check if trying to follow yourself
	if followerID == followeeID {
		return errors.New("cannot follow yourself")
	}

	// Check if already following
	exists, err := s.IsFollowing(followerID, followeeID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already following this user")
	}

	now := time.Now().Unix()
	query := `
		INSERT INTO follows (follower_id, followee_id, created_at)
		VALUES (?, ?, ?)
	`

	_, err = s.db.Exec(query, followerID, followeeID, now)
	if err != nil {
		return fmt.Errorf("failed to follow: %w", err)
	}

	return nil
}

// Unfollow removes a follow relationship
func (s *SocialStore) Unfollow(followerID, followeeID string) error {
	query := `
		DELETE FROM follows
		WHERE follower_id = ? AND followee_id = ?
	`

	result, err := s.db.Exec(query, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("failed to unfollow: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("not following this user")
	}

	return nil
}

// IsFollowing checks if follower follows followee
func (s *SocialStore) IsFollowing(followerID, followeeID string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM follows
		WHERE follower_id = ? AND followee_id = ?
	`

	var count int
	err := s.db.QueryRow(query, followerID, followeeID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check follow status: %w", err)
	}

	return count > 0, nil
}

// GetFollowing returns list of users that the given user follows
func (s *SocialStore) GetFollowing(userID string) ([]models.User, error) {
	query := `
		SELECT u.id, u.username, u.created_at
		FROM users u
		JOIN follows f ON u.id = f.followee_id
		WHERE f.follower_id = ?
		ORDER BY u.username
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// GetFollowers returns list of users that follow the given user
func (s *SocialStore) GetFollowers(userID string) ([]models.User, error) {
	query := `
		SELECT u.id, u.username, u.created_at
		FROM users u
		JOIN follows f ON u.id = f.follower_id
		WHERE f.followee_id = ?
		ORDER BY u.username
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// GetFollowingByUsername returns list of users that the given username follows
func (s *SocialStore) GetFollowingByUsername(username string) ([]models.User, error) {
	// First get the user ID
	userStore := NewUserStore(s.db)
	user, err := userStore.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	return s.GetFollowing(user.ID)
}

// GetFollowersByUsername returns list of users that follow the given username
func (s *SocialStore) GetFollowersByUsername(username string) ([]models.User, error) {
	// First get the user ID
	userStore := NewUserStore(s.db)
	user, err := userStore.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	return s.GetFollowers(user.ID)
}

// GetFollowCounts returns the number of following and followers for a user
func (s *SocialStore) GetFollowCounts(userID string) (following int, followers int, err error) {
	// Get following count
	query := `SELECT COUNT(*) FROM follows WHERE follower_id = ?`
	err = s.db.QueryRow(query, userID).Scan(&following)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get following count: %w", err)
	}

	// Get followers count
	query = `SELECT COUNT(*) FROM follows WHERE followee_id = ?`
	err = s.db.QueryRow(query, userID).Scan(&followers)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get followers count: %w", err)
	}

	return following, followers, nil
}
