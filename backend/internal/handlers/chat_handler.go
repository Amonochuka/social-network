package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/backend/internal/middleware"
	"social-network/backend/internal/services"
	"social-network/backend/internal/ws"
)

type ChatHandler struct {
	chatService *services.ChatService
	hub         *ws.Hub
}

func NewChatHandler(chatService *services.ChatService, hub *ws.Hub) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		hub:         hub,
	}
}

type sendMessageRequest struct {
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

// SendPrivateMessage handles POST /api/chat/private
func (h *ChatHandler) SendPrivateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	senderID := r.Context().Value(middleware.UserIDKey).(string)

	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ReceiverID == "" || req.Content == "" {
		http.Error(w, "receiver_id and content are required", http.StatusBadRequest)
		return
	}

	msg, err := h.chatService.SendPrivateMessage(senderID, req.ReceiverID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// push the new message out live to the receiver if they're connected
	h.hub.SendToUser(req.ReceiverID, map[string]interface{}{
		"type":    "private_message",
		"message": msg,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

// GetPrivateMessages handles GET /api/chat/private/{user_id}
func (h *ChatHandler) GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	otherUserID := r.PathValue("user_id")

	if otherUserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	messages, err := h.chatService.GetPrivateMessages(userID, otherUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// ServeWebSocket handles GET /api/chat/ws — upgrades to a websocket connection.
func (h *ChatHandler) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	ws.ServeWS(h.hub, w, r, userID)
}

// GetConversations handles GET /api/chat/conversations
func (h *ChatHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)

	conversations, err := h.chatService.GetConversations(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}