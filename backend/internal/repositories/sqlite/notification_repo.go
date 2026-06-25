package sqlite

import (
	"database/sql"
	"encoding/json"
	"log"

	"social-network/backend/internal/models"
	"social-network/backend/internal/ws"

	"github.com/google/uuid"
)

type notificationRepository struct {
	db  *sql.DB
	hub *ws.Hub // Inject the real-time hub dependency
}

func NewNotificationRepository(db *sql.DB, hub *ws.Hub) *notificationRepository {
	return &notificationRepository{
		db:  db,
		hub: hub,
	}
}

func (r *notificationRepository) CreateNotification(n *models.Notification) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	_, err := r.db.Exec(`
                INSERT INTO notifications (id, user_id, actor_id, type, ref_id, is_read, created_at)
                VALUES (?, ?, ?, ?, ?, ?, ?)
        `, n.ID, n.UserID, n.ActorID, n.Type, n.ReferenceID, n.IsRead, n.CreatedAt)
	if err != nil {
		return err
	}

	// ── AUTOMATED REAL-TIME DISPATCH ──
	// Marshall the notification row into a clean JSON frame packet
	payload, err := json.Marshal(map[string]interface{}{
		"type":    "new_notification",
		"payload": n,
	})
	if err == nil {
		// Deliver instantly to the user's notification socket pool
		r.hub.SendNotification(n.UserID, payload)
		log.Printf("Automated real-time notification pushed for user %s", n.UserID)
	} else {
		log.Printf("Failed to marshal automatic notification: %v", err)
	}

	return nil
}

func (r *notificationRepository) GetNotificationsByUserID(userID string) ([]*models.Notification, error) {
	rows, err := r.db.Query(`
                SELECT id, user_id, actor_id, type, ref_id, is_read, created_at
                FROM notifications
                WHERE user_id = ?
                ORDER BY created_at DESC
        `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notifications []*models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(&n.ID, &n.UserID, &n.ActorID, &n.Type, &n.ReferenceID, &n.IsRead, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}
	return notifications, nil
}

func (r *notificationRepository) MarkNotificationAsRead(notificationID, userID string) error {
	_, err := r.db.Exec(`
                UPDATE notifications SET is_read = 1
                WHERE id = ? AND user_id = ?
        `, notificationID, userID)
	return err
}
