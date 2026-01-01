package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RazinShafayet2007/twitter-cli/internal/models"
	"github.com/oklog/ulid/v2"
)

type PostStore struct {
	db *sql.DB
}

func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{db: db}
}

// Create creates a new post
func (s *PostStore) Create(authorID, text string) (*models.Post, error) {
	id := ulid.Make().String()
	now := time.Now().Unix()

	query := `
		INSERT INTO posts (id, author_id, text, created_at, is_retweet)
		VALUES (?, ?, ?, ?, 0)
	`

	_, err := s.db.Exec(query, id, authorID, text, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return &models.Post{
		ID:        id,
		AuthorID:  authorID,
		Text:      text,
		CreatedAt: now,
		IsRetweet: false,
	}, nil
}

// GetByID retrieves a single post by ID
func (s *PostStore) GetByID(postID string) (*models.Post, error) {
	query := `
		SELECT id, author_id, text, created_at, is_retweet, original_post_id
		FROM posts
		WHERE id = ?
	`

	var post models.Post
	err := s.db.QueryRow(query, postID).Scan(
		&post.ID,
		&post.AuthorID,
		&post.Text,
		&post.CreatedAt,
		&post.IsRetweet,
		&post.OriginalPostID,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &post, nil
}

// PostWithAuthor represents a post with author information
type PostWithAuthor struct {
	Post     models.Post
	Username string
}

// GetByAuthorID retrieves all posts by a specific author
func (s *PostStore) GetByAuthorID(authorID string, limit int) ([]PostWithAuthor, error) {
	query := `
		SELECT 
			p.id, p.author_id, p.text, p.created_at, p.is_retweet, p.original_post_id,
			u.username
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.author_id = ?
		ORDER BY p.created_at DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, authorID, limit)
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

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, nil
}

// GetByUsername retrieves all posts by username
func (s *PostStore) GetByUsername(username string, limit int) ([]PostWithAuthor, error) {
	query := `
		SELECT 
			p.id, p.author_id, p.text, p.created_at, p.is_retweet, p.original_post_id,
			u.username
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE u.username = ?
		ORDER BY p.created_at DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, username, limit)
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

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, nil
}

// Delete deletes a post (only by the author)
func (s *PostStore) Delete(postID, authorID string) error {
	query := `
		DELETE FROM posts
		WHERE id = ? AND author_id = ?
	`

	result, err := s.db.Exec(query, postID, authorID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found or you don't own this post")
	}

	return nil
}

// GetFeed returns posts for a user's feed (posts from followed users + own posts)
func (s *PostStore) GetFeed(userID string, limit, offset int) ([]PostWithAuthor, error) {
	// This query gets:
	// 1. Posts from users that userID follows
	// 2. Posts from userID themselves
	// Ordered by creation time (newest first)

	query := `
		SELECT 
			p.id, 
			p.author_id, 
			p.text, 
			p.created_at, 
			p.is_retweet, 
			p.original_post_id,
			u.username
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.author_id IN (
			SELECT followee_id 
			FROM follows 
			WHERE follower_id = ?
		)
		OR p.author_id = ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, userID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query feed: %w", err)
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

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, nil
}
