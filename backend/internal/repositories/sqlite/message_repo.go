package sqlite

import (
	"database/sql"
	"social-network/backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *messageRepository {
	return &messageRepository{db: db}
}

// ── Private messages ─────────────────────────────────────────────────────────

func (r *messageRepository) SavePrivateMessage(msg *models.PrivateMessage) error {
	msg.ID = uuid.New().String()
	msg.CreatedAt = time.Now()

	_, err := r.db.Exec(
		`INSERT INTO private_messages (id, sender_id, receiver_id, content, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		msg.ID, msg.SenderID, msg.ReceiverID, msg.Content, msg.CreatedAt,
	)
	return err
}

func (r *messageRepository) GetPrivateMessages(userA, userB string, limit int) ([]*models.PrivateMessage, error) {
	rows, err := r.db.Query(
		`SELECT id, sender_id, receiver_id, content, created_at
		 FROM private_messages
		 WHERE (sender_id = ? AND receiver_id = ?)
		    OR (sender_id = ? AND receiver_id = ?)
		 ORDER BY created_at DESC
		 LIMIT ?`,
		userA, userB, userB, userA, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*models.PrivateMessage
	for rows.Next() {
		m := &models.PrivateMessage{}
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

// ── Group messages ───────────────────────────────────────────────────────────

func (r *messageRepository) SaveGroupMessage(msg *models.GroupMessage) error {
	msg.ID = uuid.New().String()
	msg.CreatedAt = time.Now()

	_, err := r.db.Exec(
		`INSERT INTO group_messages (id, group_id, sender_id, content, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		msg.ID, msg.GroupID, msg.SenderID, msg.Content, msg.CreatedAt,
	)
	return err
}

func (r *messageRepository) GetGroupMessages(groupID string, limit int) ([]*models.GroupMessage, error) {
	rows, err := r.db.Query(
		`SELECT id, group_id, sender_id, content, created_at
		 FROM group_messages
		 WHERE group_id = ?
		 ORDER BY created_at DESC
		 LIMIT ?`,
		groupID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*models.GroupMessage
	for rows.Next() {
		m := &models.GroupMessage{}
		if err := rows.Scan(&m.ID, &m.GroupID, &m.SenderID, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

// ── Group membership ─────────────────────────────────────────────────────────

func (r *messageRepository) IsGroupMember(userID, groupID string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM group_members WHERE user_id = ? AND group_id = ?`,
		userID, groupID,
	).Scan(&count)
	return count > 0, err
}
