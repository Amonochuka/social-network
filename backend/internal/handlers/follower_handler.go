package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/backend/internal/middleware"
	"social-network/backend/internal/services"
	"social-network/backend/internal/ws"
)

type FollowerHandler struct {
	followerService *services.FollowerService
	hub             *ws.Hub
}

func NewFollowerHandler(
	followerService *services.FollowerService,
	hub *ws.Hub,
) *FollowerHandler {
	return &FollowerHandler{
		followerService: followerService,
		hub:             hub,
	}
}

func (h *FollowerHandler) SendFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	senderID := middleware.GetUserID(r)
	if senderID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ReceiverID string `json:"receiver_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ReceiverID == "" {
		http.Error(w, "receiver_id is required", http.StatusBadRequest)
		return
	}

	notification, err := h.followerService.SendFollowRequest(senderID, req.ReceiverID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if notification != nil {
		h.hub.SendToUser(notification.UserID, notification)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "follow request sent",
	})
}

func (h *FollowerHandler) AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	receiverID := middleware.GetUserID(r)
	if receiverID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	requestID := r.PathValue("request_id")
	if requestID == "" {
		http.Error(w, "request_id is required", http.StatusBadRequest)
		return
	}

	notification, err := h.followerService.AcceptFollowRequest(requestID, receiverID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if notification != nil {
		h.hub.SendToUser(notification.UserID, notification)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "follow request accepted",
	})
}

func (h *FollowerHandler) DeclineFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	receiverID := middleware.GetUserID(r)
	if receiverID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	requestID := r.PathValue("request_id")
	if requestID == "" {
		http.Error(w, "request_id is required", http.StatusBadRequest)
		return
	}

	notification, err := h.followerService.DeclineFollowRequest(requestID, receiverID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if notification != nil {
		h.hub.SendToUser(notification.UserID, notification)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "follow request declined",
	})
}

func (h *FollowerHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	followerID := middleware.GetUserID(r)
	if followerID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	followingID := r.PathValue("user_id")
	if followingID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	if err := h.followerService.Unfollow(followerID, followingID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "unfollowed successfully",
	})
}

func (h *FollowerHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	followers, err := h.followerService.GetFollowers(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func (h *FollowerHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	following, err := h.followerService.GetFollowing(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}

func (h *FollowerHandler) GetFollowStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	viewerID := middleware.GetUserID(r)
	if viewerID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	targetID := r.PathValue("user_id")
	if targetID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	status, err := h.followerService.GetFollowStatus(viewerID, targetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status": status,
	})
}