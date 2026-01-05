package media

import (
	"crypto/sha256"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	MaxImageSize     = 5 * 1024 * 1024 // 5MB
	MaxImagesPerPost = 4
)

var AllowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

// GetMediaDir returns the media storage directory
func GetMediaDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./media"
	}
	return filepath.Join(home, ".twitter-cli", "media")
}

// EnsureMediaDir creates media directory if it doesn't exist
func EnsureMediaDir() error {
	dir := GetMediaDir()
	return os.MkdirAll(dir, 0755)
}

// ValidateImage checks if file is a valid image
func ValidateImage(filePath string) error {
	// Check file exists
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Check file size
	if info.Size() > MaxImageSize {
		return fmt.Errorf("image too large (max 5MB)")
	}

	// Check file type by opening and decoding
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Try to decode as image
	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return fmt.Errorf("not a valid image file")
	}

	// Check format
	mimeType := "image/" + format
	if !AllowedTypes[mimeType] {
		return fmt.Errorf("unsupported image format: %s (only JPEG, PNG, GIF allowed)", format)
	}

	return nil
}

// GetImageDimensions returns width and height of an image
func GetImageDimensions(filePath string) (int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return config.Width, config.Height, nil
}

// GetFileType returns MIME type of an image
func GetFileType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return "", err
	}

	return "image/" + format, nil
}

// CopyImageToMedia copies an image to media directory with unique name
func CopyImageToMedia(sourcePath, postID string, position int) (string, string, error) {
	if err := EnsureMediaDir(); err != nil {
		return "", "", fmt.Errorf("failed to create media directory: %w", err)
	}

	// Open source file
	source, err := os.Open(sourcePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	// Get file extension
	ext := strings.ToLower(filepath.Ext(sourcePath))
	if ext == "" {
		ext = ".jpg" // default
	}

	// Generate unique filename: postID_position_hash.ext
	hash := sha256.New()
	io.Copy(hash, source)
	hashStr := fmt.Sprintf("%x", hash.Sum(nil))[:8]

	fileName := fmt.Sprintf("%s_%d_%s%s", postID, position, hashStr, ext)
	destPath := filepath.Join(GetMediaDir(), fileName)

	// Reset file pointer
	source.Seek(0, 0)

	// Create destination file
	dest, err := os.Create(destPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dest.Close()

	// Copy file
	_, err = io.Copy(dest, source)
	if err != nil {
		return "", "", fmt.Errorf("failed to copy file: %w", err)
	}

	return destPath, fileName, nil
}

// DeleteMediaFile deletes a media file from disk
func DeleteMediaFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	// Only delete if it's in our media directory
	mediaDir := GetMediaDir()
	if !strings.HasPrefix(filePath, mediaDir) {
		return fmt.Errorf("file not in media directory")
	}

	return os.Remove(filePath)
}
