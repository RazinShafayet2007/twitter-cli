package store

import (
	"database/sql"
	"errors"
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

// CreateReply creates a new reply to a post
func (s *PostStore) CreateReply(authorID, text, parentPostID string) (*models.Post, error) {
	// Verify parent exists
	_, err := s.GetByID(parentPostID)
	if err != nil {
		return nil, fmt.Errorf("parent post not found: %w", err)
	}

	id := ulid.Make().String()
	now := time.Now().Unix()

	query := `
		INSERT INTO posts (id, author_id, text, created_at, is_retweet, parent_post_id)
		VALUES (?, ?, ?, ?, 0, ?)
	`

	_, err = s.db.Exec(query, id, authorID, text, now, parentPostID)
	if err != nil {
		return nil, fmt.Errorf("failed to create reply: %w", err)
	}

	return &models.Post{
		ID:           id,
		AuthorID:     authorID,
		Text:         text,
		CreatedAt:    now,
		IsRetweet:    false,
		ParentPostID: &parentPostID,
	}, nil
}

// GetByID retrieves a single post by ID
func (s *PostStore) GetByID(postID string) (*models.Post, error) {
	query := `
		SELECT id, author_id, text, created_at, is_retweet, original_post_id, parent_post_id
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
		&post.ParentPostID,
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
			p.id, p.author_id, p.text, p.created_at, p.is_retweet, p.original_post_id, p.parent_post_id,
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
			&pwa.Post.ParentPostID,
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
			p.id, p.author_id, p.text, p.created_at, p.is_retweet, p.original_post_id, p.parent_post_id,
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
			&pwa.Post.ParentPostID,
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
			p.parent_post_id,
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
			&pwa.Post.ParentPostID,
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

// Retweet creates a retweet of an existing post
func (s *PostStore) Retweet(userID, originalPostID string) (*models.Post, error) {
	// Verify original post exists
	originalPost, err := s.GetByID(originalPostID)
	if err != nil {
		return nil, fmt.Errorf("original post not found")
	}

	// Check if user already retweeted this post
	hasRetweeted, err := s.HasRetweeted(userID, originalPostID)
	if err != nil {
		return nil, err
	}
	if hasRetweeted {
		return nil, errors.New("already retweeted this post")
	}

	// Can't retweet your own post
	if originalPost.AuthorID == userID {
		return nil, errors.New("cannot retweet your own post")
	}

	id := ulid.Make().String()
	now := time.Now().Unix()

	query := `
		INSERT INTO posts (id, author_id, text, created_at, is_retweet, original_post_id)
		VALUES (?, ?, ?, ?, 1, ?)
	`

	// Retweet uses the original post's text
	_, err = s.db.Exec(query, id, userID, originalPost.Text, now, originalPostID)
	if err != nil {
		return nil, fmt.Errorf("failed to create retweet: %w", err)
	}

	return &models.Post{
		ID:             id,
		AuthorID:       userID,
		Text:           originalPost.Text,
		CreatedAt:      now,
		IsRetweet:      true,
		OriginalPostID: &originalPostID,
	}, nil
}

// HasRetweeted checks if a user has retweeted a post
func (s *PostStore) HasRetweeted(userID, originalPostID string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM posts
		WHERE author_id = ? AND original_post_id = ? AND is_retweet = 1
	`

	var count int
	err := s.db.QueryRow(query, userID, originalPostID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check retweet status: %w", err)
	}

	return count > 0, nil
}

// GetRetweetCount returns the number of retweets for a post
func (s *PostStore) GetRetweetCount(postID string) (int, error) {
	query := `
		SELECT COUNT(*) FROM posts
		WHERE original_post_id = ? AND is_retweet = 1
	`

	var count int
	err := s.db.QueryRow(query, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get retweet count: %w", err)
	}

	return count, nil
}

// Search searches posts by text
func (s *PostStore) Search(query string, limit int) ([]PostWithAuthor, error) {
	sqlQuery := `
		SELECT 
			p.id, p.author_id, p.text, p.created_at, p.is_retweet, p.original_post_id, p.parent_post_id,
			u.username
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.text LIKE ?
		ORDER BY p.created_at DESC
		LIMIT ?
	`

	searchTerm := "%" + query + "%"

	rows, err := s.db.Query(sqlQuery, searchTerm, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
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
			&pwa.Post.ParentPostID,
			&pwa.Username,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, pwa)
	}

	return posts, nil
}

// GetThread retrieves the thread context for a post (ancestors + post + direct replies)
func (s *PostStore) GetThread(postID string) ([]PostWithAuthor, error) {
	// Reusable recursive CTE to get ancestors and children would be nice, but simple approach:
	// 1. Get the requested post
	// 2. Walk up to find ancestors (or use CTE)
	// 3. Get direct children

	// CTE for ancestors + self + children
	query := `
		WITH RECURSIVE ancestors AS (
			SELECT id, author_id, text, created_at, is_retweet, original_post_id, parent_post_id, 0 as level
			FROM posts
			WHERE id = ?
			
			UNION ALL
			
			SELECT p.id, p.author_id, p.text, p.created_at, p.is_retweet, p.original_post_id, p.parent_post_id, a.level - 1
			FROM posts p
			JOIN ancestors a ON p.id = a.parent_post_id
		),
		children AS (
			SELECT id, author_id, text, created_at, is_retweet, original_post_id, parent_post_id, 1 as level
			FROM posts
			WHERE parent_post_id = ?
		)
		SELECT 
			p.id, p.author_id, p.text, p.created_at, p.is_retweet, p.original_post_id, p.parent_post_id,
			u.username
		FROM (
			SELECT * FROM ancestors
			UNION ALL
			SELECT * FROM children
		) p
		JOIN users u ON p.author_id = u.id
		ORDER BY p.level ASC, p.created_at ASC
	`

	rows, err := s.db.Query(query, postID, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to query thread: %w", err)
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
			&pwa.Post.ParentPostID,
			&pwa.Username,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, pwa)
	}

	return posts, nil
}
