package models

type Follow struct {
	FollowerID string
	FolloweeID string
	CreatedAt  int64
}

type Like struct {
	UserID    string
	PostID    string
	CreatedAt int64
}
