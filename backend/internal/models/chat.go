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

type GroupMessageDetailFull struct {
	ID           string    `json:"id"`
	GroupID      string    `json:"group_id"`
	SenderID     string    `json:"sender_id"`
	SenderName   string    `json:"sender_name"`
	SenderAvatar string    `json:"sender_avatar"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
}

type ConversationPreview struct {
	UserID       string    `json:"user_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Avatar       string    `json:"avatar"`
	LastMessage  string    `json:"last_message"`
	LastMessageAt time.Time `json:"last_message_at"`
}