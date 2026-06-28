package sqlite

import (
	"database/sql"
	"time"

	"social-network/backend/internal/models"

	"github.com/google/uuid"
)

type chatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *chatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) CreatePrivateMessage(msg *models.PrivateMessage) error {
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}
	_, err := r.db.Exec(`
		INSERT INTO private_messages (id, sender_id, receiver_id, content, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, msg.ID, msg.SenderID, msg.ReceiverID, msg.Content, msg.CreatedAt)
	return err
}

func (r *chatRepository) GetPrivateMessages(userA, userB string) ([]*models.PrivateMessageDetail, error) {
	rows, err := r.db.Query(`
		SELECT pm.id, pm.sender_id, u.first_name || ' ' || u.last_name AS sender_name, u.avatar,
		       pm.receiver_id, pm.content, pm.created_at
		FROM private_messages pm
		JOIN users u ON u.id = pm.sender_id
		WHERE (pm.sender_id = ? AND pm.receiver_id = ?)
		   OR (pm.sender_id = ? AND pm.receiver_id = ?)
		ORDER BY pm.created_at ASC
	`, userA, userB, userB, userA)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.PrivateMessageDetail
	for rows.Next() {
		var m models.PrivateMessageDetail
		err := rows.Scan(&m.ID, &m.SenderID, &m.SenderName, &m.SenderAvatar, &m.ReceiverID, &m.Content, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}
	return messages, nil
}

func (r *chatRepository) CreateGroupMessage(msg *models.GroupMessage) error {
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}
	_, err := r.db.Exec(`
		INSERT INTO group_messages (id, group_id, sender_id, content, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, msg.ID, msg.GroupID, msg.SenderID, msg.Content, msg.CreatedAt)
	return err
}

func (r *chatRepository) GetGroupMessages(groupID string) ([]*models.GroupMessageDetail, error) {
	rows, err := r.db.Query(`
		SELECT gm.id, gm.group_id, gm.sender_id, u.first_name || ' ' || u.last_name AS sender_name, u.avatar,
		       gm.content, gm.created_at
		FROM group_messages gm
		JOIN users u ON u.id = gm.sender_id
		WHERE gm.group_id = ?
		ORDER BY gm.created_at ASC
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.GroupMessageDetail
	for rows.Next() {
		var m models.GroupMessageDetail
		err := rows.Scan(&m.ID, &m.GroupID, &m.SenderID, &m.SenderName, &m.SenderAvatar, &m.Content, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}
	return messages, nil
}
