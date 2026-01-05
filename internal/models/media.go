package models

type Media struct {
	ID        string
	PostID    string
	FilePath  string
	FileName  string
	FileType  string // "image/jpeg", "image/png", "image/gif"
	FileSize  int64
	Width     *int
	Height    *int
	Position  int // 0, 1, 2, 3 for multiple images
	CreatedAt int64
}