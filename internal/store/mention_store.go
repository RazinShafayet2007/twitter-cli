package store

import (
	"database/sql"
	"fmt"
	"time"
)

type MentionStore struct {
	db *sql.DB
}

func NewMentionStore(db *sql.DB) *MentionStore {
	return &MentionStore{db: db}
}

// CreateMentions creates mention records for a post
func (s *MentionStore) CreateMentions(postID string, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	now := time.Now().Unix()
	query := `INSERT OR IGNORE INTO mentions (post_id, mentioned_user_id, created_at) VALUES (?, ?, ?)`

	for _, userID := range userIDs {
		_, err := tx.Exec(query, postID, userID, now)
		if err != nil {
			return fmt.Errorf("failed to create mention: %w", err)
		}
	}

	return tx.Commit()
}

// GetMentions retrieves posts that mention a user
func (s *MentionStore) GetMentions(userID string, limit int) ([]PostWithAuthor, error) {
	query := `
		SELECT 
			p.id, p.author_id, p.text, p.created_at, p.is_retweet, p.original_post_id,
			u.username
		FROM posts p
		JOIN users u ON p.author_id = u.id
		JOIN mentions m ON p.id = m.post_id
		WHERE m.mentioned_user_id = ?
		ORDER BY p.created_at DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query mentions: %w", err)
	}
	defer rows.Close()

	var posts []PostWithAuthor
	for rows.Next() {
		var pwa PostWithAuthor
		err := rows.Scan(
			&pwa.Post.ID,
			&pwa.Post.AuthorID,
			&pwa.Post.Text,
			&pwa.Post.CreatedAt,
			&pwa.Post.IsRetweet,
			&pwa.Post.OriginalPostID,
			&pwa.Username,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, pwa)
	}

	return posts, nil
}

// GetMentionedUsers gets user IDs from usernames
func (s *MentionStore) GetMentionedUsers(usernames []string) ([]string, error) {
	if len(usernames) == 0 {
		return nil, nil
	}

	// Build placeholders for IN clause
	placeholders := ""
	for i := range usernames {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
	}

	query := fmt.Sprintf(`SELECT id FROM users WHERE username IN (%s)`, placeholders)

	// Convert usernames to interface{} for variadic args
	args := make([]interface{}, len(usernames))
	for i, username := range usernames {
		args[i] = username
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan user ID: %w", err)
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}
