package models

import "time"

// MessageType tells the hub what to do with the payload
type MessageType string

const (
	TypePrivateMessage MessageType = "private_message"
	TypeGroupMessage   MessageType = "group_message"
	TypeNotification   MessageType = "notification"
	TypeError          MessageType = "error"
)

// PrivateMessage maps to the private_messages table
type PrivateMessage struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// GroupMessage maps to the group_messages table
type GroupMessage struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// WSMessage is the envelope sent over every WebSocket connection.
// The frontend encodes outgoing messages in this shape;
// the hub inspects Type and routes accordingly.
type WSMessage struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// IncomingChat is what the client sends when typing a message.
// For private chat: TargetID = receiverID.
// For group chat:   TargetID = groupID.
type IncomingChat struct {
	Content  string `json:"content"`
	TargetID string `json:"target_id"`
}
