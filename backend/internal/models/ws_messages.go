package models

import "encoding/json"

// WSMessageFrame is the generic shell wrapper for any incoming socket packet.
type WSMessageFrame struct {
	Type    string          `json:"type"`    // e.g., "private_chat", "group_chat"
	Payload json.RawMessage `json:"payload"` // Delayed parsing until type is verified
}

// PrivateChatPayload maps out an explicit direct message structure.
type PrivateChatPayload struct {
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

// GroupChatPayload maps out an explicit community group message structure.
type GroupChatPayload struct {
	GroupID string `json:"group_id"`
	Content string `json:"content"`
}
