package store

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/RazinShafayet2007/twitter-cli/internal/models"
	"github.com/oklog/ulid/v2"
)

type MessageStore struct {
	db *sql.DB
}

func NewMessageStore(db *sql.DB) *MessageStore {
	return &MessageStore{db: db}
}

// GetDB returns the underlying *sql.DB instance
func (s *MessageStore) GetDB() *sql.DB {
	return s.db
}

// Send creates a new message
func (s *MessageStore) Send(senderID, receiverID, text string) (*models.Message, error) {
	id := ulid.Make().String()
	now := time.Now().Unix()

	query := `
		INSERT INTO messages (id, sender_id, receiver_id, text, created_at, read)
		VALUES (?, ?, ?, ?, ?, 0)
	`

	_, err := s.db.Exec(query, id, senderID, receiverID, text, now)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Check if sender is blocked by receiver
	blocked, err := s.IsBlocked(receiverID, senderID)
	if err != nil {
		return nil, err
	}
	if blocked {
		return nil, fmt.Errorf("you cannot send messages to this user")
	}

	return &models.Message{
		ID:         id,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Text:       text,
		CreatedAt:  now,
		Read:       false,
	}, nil
}

// GetInbox retrieves messages received by a user
func (s *MessageStore) GetInbox(userID string, limit int) ([]models.MessageWithUser, error) {
	query := `
		SELECT 
			m.id, m.sender_id, m.receiver_id, m.text, m.created_at, m.read,
			u.username as sender_name
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.receiver_id = ?
		ORDER BY m.created_at DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get inbox: %w", err)
	}
	defer rows.Close()

	var messages []models.MessageWithUser
	for rows.Next() {
		var m models.MessageWithUser
		var readInt int
		err := rows.Scan(
			&m.Message.ID,
			&m.Message.SenderID,
			&m.Message.ReceiverID,
			&m.Message.Text,
			&m.Message.CreatedAt,
			&readInt,
			&m.SenderName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		m.Message.Read = readInt == 1
		messages = append(messages, m)
	}

	return messages, nil
}

// GetConversation retrieves messages between two users
func (s *MessageStore) GetConversation(user1ID, user2ID string, limit int) ([]models.MessageWithUser, error) {
	query := `
		SELECT 
			m.id, m.sender_id, m.receiver_id, m.text, m.created_at, m.read,
			sender.username as sender_name,
			receiver.username as receiver_name
		FROM messages m
		JOIN users sender ON m.sender_id = sender.id
		JOIN users receiver ON m.receiver_id = receiver.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?)
		   OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at ASC
		LIMIT ?
	`

	rows, err := s.db.Query(query, user1ID, user2ID, user2ID, user1ID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	defer rows.Close()

	var messages []models.MessageWithUser
	for rows.Next() {
		var m models.MessageWithUser
		var readInt int
		err := rows.Scan(
			&m.Message.ID,
			&m.Message.SenderID,
			&m.Message.ReceiverID,
			&m.Message.Text,
			&m.Message.CreatedAt,
			&readInt,
			&m.SenderName,
			&m.ReceiverName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		m.Message.Read = readInt == 1
		messages = append(messages, m)
	}

	return messages, nil
}

// GetConversations retrieves list of conversations with unread counts
func (s *MessageStore) GetConversations(userID string) ([]models.Conversation, error) {
	// First, get all unique conversation partners
	query := `
		SELECT DISTINCT
			CASE 
				WHEN sender_id = ? THEN receiver_id 
				ELSE sender_id 
			END as other_user_id
		FROM messages
		WHERE sender_id = ? OR receiver_id = ?
	`

	rows, err := s.db.Query(query, userID, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation partners: %w", err)
	}
	defer rows.Close()

	var partnerIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan partner ID: %w", err)
		}
		partnerIDs = append(partnerIDs, id)
	}

	if len(partnerIDs) == 0 {
		return []models.Conversation{}, nil
	}

	// Now get details for each conversation
	var conversations []models.Conversation

	for _, partnerID := range partnerIDs {
		// Get partner username
		var username string
		userQuery := `SELECT username FROM users WHERE id = ?`
		if err := s.db.QueryRow(userQuery, partnerID).Scan(&username); err != nil {
			continue // Skip if user deleted
		}

		// Get last message
		var lastText string
		var lastTime int64
		lastMsgQuery := `
			SELECT text, created_at
			FROM messages
			WHERE (sender_id = ? AND receiver_id = ?)
			   OR (sender_id = ? AND receiver_id = ?)
			ORDER BY created_at DESC
			LIMIT 1
		`
		if err := s.db.QueryRow(lastMsgQuery, userID, partnerID, partnerID, userID).Scan(&lastText, &lastTime); err != nil {
			continue
		}

		// Get unread count
		var unreadCount int
		unreadQuery := `
			SELECT COUNT(*) FROM messages
			WHERE sender_id = ? AND receiver_id = ? AND read = 0
		`
		if err := s.db.QueryRow(unreadQuery, partnerID, userID).Scan(&unreadCount); err != nil {
			unreadCount = 0
		}

		conversations = append(conversations, models.Conversation{
			OtherUserID:   partnerID,
			OtherUsername: username,
			LastMessage:   lastText,
			LastMessageAt: lastTime,
			UnreadCount:   unreadCount,
		})
	}

	// Sort by last message time (newest first)
	sort.Slice(conversations, func(i, j int) bool {
		return conversations[i].LastMessageAt > conversations[j].LastMessageAt
	})

	return conversations, nil
}

// MarkAsRead marks all messages in a conversation as read
func (s *MessageStore) MarkAsRead(receiverID, senderID string) error {
	query := `
		UPDATE messages 
		SET read = 1
		WHERE receiver_id = ? AND sender_id = ? AND read = 0
	`

	_, err := s.db.Exec(query, receiverID, senderID)
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	return nil
}

// GetUnreadCount returns total unread message count for a user
func (s *MessageStore) GetUnreadCount(userID string) (int, error) {
	query := `
		SELECT COUNT(*) FROM messages
		WHERE receiver_id = ? AND read = 0
	`

	var count int
	err := s.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

// DeleteMessage deletes a message (only sender can delete)
func (s *MessageStore) DeleteMessage(messageID, senderID string) error {
	query := `
		DELETE FROM messages
		WHERE id = ? AND sender_id = ?
	`

	result, err := s.db.Exec(query, messageID, senderID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found or you don't own it")
	}

	return nil
}

func (s *MessageStore) SearchMessages(userID, query string) ([]models.MessageWithUser, error) {
	sqlQuery := `
		SELECT 
			m.id, m.sender_id, m.receiver_id, m.text, m.created_at, m.read,
			sender.username as sender_name,
			receiver.username as receiver_name
		FROM messages m
		JOIN users sender ON m.sender_id = sender.id
		JOIN users receiver ON m.receiver_id = receiver.id
		WHERE (m.sender_id = ? OR m.receiver_id = ?)
		  AND m.text LIKE ?
		ORDER BY m.created_at DESC
		LIMIT 50
	`

	searchTerm := "%" + query + "%"

	rows, err := s.db.Query(sqlQuery, userID, userID, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search messages: %w", err)
	}
	defer rows.Close()

	var messages []models.MessageWithUser
	for rows.Next() {
		var m models.MessageWithUser
		var readInt int
		err := rows.Scan(
			&m.Message.ID,
			&m.Message.SenderID,
			&m.Message.ReceiverID,
			&m.Message.Text,
			&m.Message.CreatedAt,
			&readInt,
			&m.SenderName,
			&m.ReceiverName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		m.Message.Read = readInt == 1
		messages = append(messages, m)
	}

	return messages, nil
}

func (s *MessageStore) IsBlocked(blockerID, blockedID string) (bool, error) {
	query := `SELECT COUNT(*) FROM blocks WHERE blocker_id = ? AND blocked_id = ?`

	var count int
	err := s.db.QueryRow(query, blockerID, blockedID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check block status: %w", err)
	}

	return count > 0, nil
}
