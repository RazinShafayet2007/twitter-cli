package models

type Message struct {
	ID         string
	SenderID   string
	ReceiverID string
	Text       string
	CreatedAt  int64
	Read       bool
}

// MessageWithUser represents a message with sender/receiver info
type MessageWithUser struct {
	Message      Message
	SenderName   string
	ReceiverName string
}

// Conversation represents a message thread between two users
type Conversation struct {
	OtherUserID   string
	OtherUsername string
	LastMessage   string
	LastMessageAt int64
	UnreadCount   int
}
