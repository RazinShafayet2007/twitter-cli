package models

type Post struct {
	ID             string
	AuthorID       string
	Text           string
	CreatedAt      int64
	IsRetweet      bool
	OriginalPostID *string // pointer because it can be NULL
	ParentPostID   *string // pointer because it can be NULL
}
