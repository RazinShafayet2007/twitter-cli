package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RazinShafayet2007/twitter-cli/internal/models"
	"github.com/oklog/ulid/v2"
)

type MediaStore struct {
	db *sql.DB
}

func NewMediaStore(db *sql.DB) *MediaStore {
	return &MediaStore{db: db}
}

// Create creates a media record
func (s *MediaStore) Create(media *models.Media) error {
	if media.ID == "" {
		media.ID = ulid.Make().String()
	}
	if media.CreatedAt == 0 {
		media.CreatedAt = time.Now().Unix()
	}

	query := `
		INSERT INTO media (
			id, post_id, file_path, file_name, file_type, file_size,
			width, height, position, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		media.ID,
		media.PostID,
		media.FilePath,
		media.FileName,
		media.FileType,
		media.FileSize,
		media.Width,
		media.Height,
		media.Position,
		media.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create media: %w", err)
	}

	return nil
}

// GetByPostID retrieves all media for a post
func (s *MediaStore) GetByPostID(postID string) ([]models.Media, error) {
	query := `
		SELECT 
			id, post_id, file_path, file_name, file_type, file_size,
			width, height, position, created_at
		FROM media
		WHERE post_id = ?
		ORDER BY position
	`

	rows, err := s.db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to query media: %w", err)
	}
	defer rows.Close()

	var mediaList []models.Media
	for rows.Next() {
		var m models.Media
		err := rows.Scan(
			&m.ID,
			&m.PostID,
			&m.FilePath,
			&m.FileName,
			&m.FileType,
			&m.FileSize,
			&m.Width,
			&m.Height,
			&m.Position,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan media: %w", err)
		}
		mediaList = append(mediaList, m)
	}

	return mediaList, nil
}

// Delete deletes a media record
func (s *MediaStore) Delete(mediaID string) error {
	query := `DELETE FROM media WHERE id = ?`
	_, err := s.db.Exec(query, mediaID)
	if err != nil {
		return fmt.Errorf("failed to delete media: %w", err)
	}
	return nil
}

// DeleteByPostID deletes all media for a post
func (s *MediaStore) DeleteByPostID(postID string) error {
	query := `DELETE FROM media WHERE post_id = ?`
	_, err := s.db.Exec(query, postID)
	if err != nil {
		return fmt.Errorf("failed to delete media: %w", err)
	}
	return nil
}

// GetMediaCount returns number of media items for a post
func (s *MediaStore) GetMediaCount(postID string) (int, error) {
	query := `SELECT COUNT(*) FROM media WHERE post_id = ?`

	var count int
	err := s.db.QueryRow(query, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get media count: %w", err)
	}

	return count, nil
}
