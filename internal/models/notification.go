package models

type Notification struct {
	ID        string
	UserID    string  // Who receives the notification
	ActorID   string  // Who performed the action
	Type      string  // "like", "retweet", "follow", "message"
	TargetID  *string // Post ID, message ID, etc. (can be NULL)
	CreatedAt int64
	Read      bool
}

// NotificationWithDetails includes actor username and target details
type NotificationWithDetails struct {
	Notification Notification
	ActorName    string
	TargetText   *string // Text of the post/message if applicable
}
