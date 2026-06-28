package models

import "time"

type PrivateMessage struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type PrivateMessageDetail struct {
	ID           string    `json:"id"`
	SenderID     string    `json:"sender_id"`
	SenderName   string    `json:"sender_name"`
	SenderAvatar string    `json:"sender_avatar"`
	ReceiverID   string    `json:"receiver_id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
}

type GroupMessage struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupMessageDetail struct {
	ID           string    `json:"id"`
	GroupID      string    `json:"group_id"`
	SenderID     string    `json:"sender_id"`
	SenderName   string    `json:"sender_name"`
	SenderAvatar string    `json:"sender_avatar"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
}