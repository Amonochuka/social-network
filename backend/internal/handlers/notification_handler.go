package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/backend/internal/middleware"
	"social-network/backend/internal/models"
	"social-network/backend/internal/services"
	"social-network/backend/internal/ws"
)

// actionableTypes are notification types that stay active until the user
// explicitly accepts or declines — opening the panel alone won't mark them read.
var actionableTypes = map[string]bool{
	"follow_request":     true,
	"group_invitation":   true,
	"group_join_request": true,
}

type NotificationHandler struct {
	notificationService *services.NotificationService
	hub                 *ws.Hub
}

func NewNotificationHandler(notificationService *services.NotificationService, hub *ws.Hub) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		hub:                 hub,
	}
}

func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	notifications, err := h.notificationService.GetNotifications(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// auto-mark informational notifications as read — the user just opened the panel
	for _, n := range notifications {
		if !n.IsRead && !actionableTypes[n.Type] {
			h.notificationService.MarkAsRead(n.ID, userID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	notificationID := r.PathValue("notification_id")
	if notificationID == "" {
		http.Error(w, "notification_id is required", http.StatusBadRequest)
		return
	}

	if err := h.notificationService.MarkAsRead(notificationID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "notification marked as read"})
}

// PushNotification is called internally (by other handlers like FollowerHandler,
// GroupHandler) to push a real-time notification to a user via WebSocket after
// saving it to the DB. Not an HTTP endpoint — only called from Go code.
func (h *NotificationHandler) PushNotification(n *models.Notification) {
	h.hub.SendToUser(n.UserID, map[string]interface{}{
		"type":    "new_notification",
		"payload": n,
	})
}

// GetUnreadCount handles GET /api/notifications/unread — returns just the count
// of unread notifications for the badge in the navbar.
func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	notifications, err := h.notificationService.GetNotifications(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count := 0
	for _, n := range notifications {
		if !n.IsRead {
			count++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"unread_count": count})
}