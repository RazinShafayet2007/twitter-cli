package store

import (
	"database/sql"
	"fmt"
	"time"
)

type HashtagStore struct {
	db *sql.DB
}

func NewHashtagStore(db *sql.DB) *HashtagStore {
	return &HashtagStore{db: db}
}

// GetOrCreateHashtag gets existing hashtag or creates new one
func (s *HashtagStore) GetOrCreateHashtag(tag string) (int64, error) {
	// Try to get existing
	var id int64
	query := `SELECT id FROM hashtags WHERE tag = ?`
	err := s.db.QueryRow(query, tag).Scan(&id)

	if err == nil {
		return id, nil
	}

	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to query hashtag: %w", err)
	}

	// Create new
	insertQuery := `INSERT INTO hashtags (tag, created_at) VALUES (?, ?)`
	result, err := s.db.Exec(insertQuery, tag, time.Now().Unix())
	if err != nil {
		return 0, fmt.Errorf("failed to create hashtag: %w", err)
	}

	return result.LastInsertId()
}

// LinkPostToHashtags links a post to multiple hashtags
func (s *HashtagStore) LinkPostToHashtags(postID string, hashtags []string) error {
	if len(hashtags) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, tag := range hashtags {
		// Get or create hashtag
		hashtagID, err := s.GetOrCreateHashtag(tag)
		if err != nil {
			return err
		}

		// Link to post
		query := `INSERT OR IGNORE INTO post_hashtags (post_id, hashtag_id) VALUES (?, ?)`
		_, err = tx.Exec(query, postID, hashtagID)
		if err != nil {
			return fmt.Errorf("failed to link hashtag: %w", err)
		}
	}

	return tx.Commit()
}

// GetPostsByHashtag retrieves posts with a specific hashtag
func (s *HashtagStore) GetPostsByHashtag(tag string, limit int) ([]PostWithAuthor, error) {
	query := `
		SELECT 
			p.id, p.author_id, p.text, p.created_at, p.is_retweet, p.original_post_id,
			u.username
		FROM posts p
		JOIN users u ON p.author_id = u.id
		JOIN post_hashtags ph ON p.id = ph.post_id
		JOIN hashtags h ON ph.hashtag_id = h.id
		WHERE h.tag = ?
		ORDER BY p.created_at DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, tag, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
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

// GetTrendingHashtags gets most used hashtags
func (s *HashtagStore) GetTrendingHashtags(limit int, since int64) ([]TrendingHashtag, error) {
	query := `
		SELECT 
			h.tag,
			COUNT(ph.post_id) as post_count
		FROM hashtags h
		JOIN post_hashtags ph ON h.id = ph.hashtag_id
		JOIN posts p ON ph.post_id = p.id
		WHERE p.created_at > ?
		GROUP BY h.id, h.tag
		ORDER BY post_count DESC, h.tag ASC
		LIMIT ?
	`

	rows, err := s.db.Query(query, since, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query trending hashtags: %w", err)
	}
	defer rows.Close()

	var trending []TrendingHashtag
	for rows.Next() {
		var t TrendingHashtag
		err := rows.Scan(&t.Tag, &t.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan hashtag: %w", err)
		}
		trending = append(trending, t)
	}

	return trending, nil
}

type TrendingHashtag struct {
	Tag   string
	Count int
}
