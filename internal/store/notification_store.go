package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RazinShafayet2007/twitter-cli/internal/models"
	"github.com/oklog/ulid/v2"
)

type NotificationStore struct {
	db *sql.DB
}

func NewNotificationStore(db *sql.DB) *NotificationStore {
	return &NotificationStore{db: db}
}

// Create creates a new notification
func (s *NotificationStore) Create(userID, actorID, notifType string, targetID *string) error {
	// Don't notify yourself
	if userID == actorID {
		return nil
	}

	id := ulid.Make().String()
	now := time.Now().Unix()

	query := `
		INSERT INTO notifications (id, user_id, actor_id, type, target_id, created_at, read)
		VALUES (?, ?, ?, ?, ?, ?, 0)
	`

	_, err := s.db.Exec(query, id, userID, actorID, notifType, targetID, now)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}

// GetNotifications retrieves notifications for a user
func (s *NotificationStore) GetNotifications(userID string, unreadOnly bool, limit int) ([]models.NotificationWithDetails, error) {
	query := `
		SELECT 
			n.id, n.user_id, n.actor_id, n.type, n.target_id, n.created_at, n.read,
			u.username as actor_name,
			CASE 
				WHEN n.type IN ('like', 'retweet') THEN p.text
				WHEN n.type = 'message' THEN m.text
				ELSE NULL
			END as target_text
		FROM notifications n
		JOIN users u ON n.actor_id = u.id
		LEFT JOIN posts p ON n.target_id = p.id AND n.type IN ('like', 'retweet')
		LEFT JOIN messages m ON n.target_id = m.id AND n.type = 'message'
		WHERE n.user_id = ?
	`

	if unreadOnly {
		query += " AND n.read = 0"
	}

	query += " ORDER BY n.created_at DESC LIMIT ?"

	rows, err := s.db.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.NotificationWithDetails
	for rows.Next() {
		var n models.NotificationWithDetails
		var readInt int
		err := rows.Scan(
			&n.Notification.ID,
			&n.Notification.UserID,
			&n.Notification.ActorID,
			&n.Notification.Type,
			&n.Notification.TargetID,
			&n.Notification.CreatedAt,
			&readInt,
			&n.ActorName,
			&n.TargetText,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		n.Notification.Read = readInt == 1
		notifications = append(notifications, n)
	}

	return notifications, nil
}

// MarkAsRead marks notifications as read
func (s *NotificationStore) MarkAsRead(userID string) error {
	query := `UPDATE notifications SET read = 1 WHERE user_id = ? AND read = 0`

	_, err := s.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to mark notifications as read: %w", err)
	}

	return nil
}

// GetUnreadCount returns count of unread notifications
func (s *NotificationStore) GetUnreadCount(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = ? AND read = 0`

	var count int
	err := s.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

// DeleteNotification deletes a notification
func (s *NotificationStore) DeleteNotification(notificationID, userID string) error {
	query := `DELETE FROM notifications WHERE id = ? AND user_id = ?`

	result, err := s.db.Exec(query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// DeleteAllRead deletes all read notifications for a user
func (s *NotificationStore) DeleteAllRead(userID string) error {
	query := `DELETE FROM notifications WHERE user_id = ? AND read = 1`

	_, err := s.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete notifications: %w", err)
	}

	return nil
}
