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
	return err
}

func (r *notificationRepository) GetNotificationsByUserID(userID string) ([]*models.NotificationDetail, error) {
	rows, err := r.db.Query(`
		SELECT n.id, n.user_id, n.actor_id, u.first_name || ' ' || u.last_name AS actor_name, u.avatar,
		       n.type, n.ref_id, n.is_read, n.created_at
		FROM notifications n
		JOIN users u ON u.id = n.actor_id
		WHERE n.user_id = ?
		ORDER BY n.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notifications []*models.NotificationDetail
	for rows.Next() {
		var n models.NotificationDetail
		err := rows.Scan(&n.ID, &n.UserID, &n.ActorID, &n.ActorName, &n.ActorAvatar, &n.Type, &n.ReferenceID, &n.IsRead, &n.CreatedAt)
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
