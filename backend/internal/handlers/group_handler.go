package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"social-network/backend/internal/middleware"
	"social-network/backend/internal/services"
	"social-network/backend/internal/ws"
)

type GroupHandler struct {
	groupService *services.GroupService
	hub          *ws.Hub
}

func NewGroupHandler(groupService *services.GroupService, hub *ws.Hub) *GroupHandler {
	return &GroupHandler{groupService: groupService, hub: hub}
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func (h *GroupHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *GroupHandler) getUserID(r *http.Request) string {
	return middleware.GetUserID(r)
}

// ── POST /api/groups ──────────────────────────────────────────────────────────

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	group, err := h.groupService.CreateGroup(userID, req.Title, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.respondJSON(w, http.StatusCreated, group)
}

// ── GET /api/groups ───────────────────────────────────────────────────────────

func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	groups, err := h.groupService.GetAllGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.respondJSON(w, http.StatusOK, groups)
}

// ── GET /api/groups/{group_id} ────────────────────────────────────────────────

func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	groupID := r.PathValue("group_id")
	detail, err := h.groupService.GetGroup(groupID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.respondJSON(w, http.StatusOK, detail)
}

// ── POST /api/groups/{group_id}/invite ────────────────────────────────────────

func (h *GroupHandler) InviteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserID(r)
	groupID := r.PathValue("group_id")

	var req struct {
		InviteeID string `json:"invitee_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	notification, err := h.groupService.InviteUser(groupID, userID, req.InviteeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if notification != nil {
		h.hub.SendToUser(notification.UserID, notification)
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "invitation sent",
	})
}

// ── POST /api/groups/{group_id}/join ─────────────────────────────────────────

func (h *GroupHandler) RequestToJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	groupID := r.PathValue("group_id")
	if err := h.groupService.RequestToJoin(groupID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]string{"message": "join request sent"})
}

// ── POST /api/groups/invitations/{inv_id}/accept ─────────────────────────────

func (h *GroupHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserID(r)
	invID := r.PathValue("inv_id")

	if err := h.groupService.AcceptInvitation(invID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "invitation accepted",
	})
}

// ── POST /api/groups/invitations/{inv_id}/decline ─────────────────────────────

func (h *GroupHandler) DeclineInvitation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserID(r)
	invID := r.PathValue("inv_id")

	if err := h.groupService.DeclineInvitation(invID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "invitation declined",
	})
}

// ── POST /api/groups/requests/{req_id}/accept ─────────────────────────────────

func (h *GroupHandler) AcceptJoinRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	reqID := r.PathValue("req_id")
	notification, err := h.groupService.AcceptJoinRequest(reqID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if notification != nil {
		h.hub.SendToUser(notification.UserID, notification)
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "join request accepted",
	})
}

// ── POST /api/groups/requests/{req_id}/decline ────────────────────────────────

func (h *GroupHandler) DeclineJoinRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	reqID := r.PathValue("req_id")
	notification, err := h.groupService.DeclineJoinRequest(reqID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if notification != nil {
		h.hub.SendToUser(notification.UserID, notification)
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "join request declined",
	})
}

// ── POST /api/groups/{group_id}/events ───────────────────────────────────────

func (h *GroupHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	groupID := r.PathValue("group_id")
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		EventTime   string `json:"event_time"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	event, err := h.groupService.CreateEvent(groupID, userID, req.Title, req.Description, req.EventTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.respondJSON(w, http.StatusCreated, event)
}

// ── GET /api/groups/{group_id}/events ────────────────────────────────────────

func (h *GroupHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	groupID := r.PathValue("group_id")
	events, err := h.groupService.GetEvents(groupID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	h.respondJSON(w, http.StatusOK, events)
}

// ── POST /api/groups/events/{event_id}/rsvp ──────────────────────────────────

func (h *GroupHandler) RSVPEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	eventID := r.PathValue("event_id")
	var req struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.groupService.RSVPEvent(eventID, userID, req.Response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]string{"message": "response saved"})
}

// ── POST /api/groups/{group_id}/posts ────────────────────────────────────────

func (h *GroupHandler) CreateGroupPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	groupID := r.PathValue("group_id")

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}
	content := r.FormValue("content")
	mediaPath, mediaType := "", ""
	file, header, err := r.FormFile("media")
	if err == nil {
		defer file.Close()
		ext := strings.ToLower(filepath.Ext(header.Filename))
		allowedExts := map[string]string{".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png", ".gif": "image/gif"}
		mime, ok := allowedExts[ext]
		if !ok {
			http.Error(w, "unsupported file type", http.StatusBadRequest)
			return
		}
		mediaType = mime
		filename := userID + "_" + header.Filename
		mediaPath = "uploads/group_posts/" + filename
		if err := os.MkdirAll("uploads/group_posts", 0755); err != nil {
			http.Error(w, "could not create upload directory", http.StatusInternalServerError)
			return
		}
		dst, err := os.Create(mediaPath)
		if err != nil {
			http.Error(w, "could not save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		buf := make([]byte, 32*1024)
		for {
			n, err := file.Read(buf)
			if n > 0 {
				dst.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
	}

	post, err := h.groupService.CreateGroupPost(groupID, userID, content, mediaPath, mediaType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.respondJSON(w, http.StatusCreated, post)
}

// ── GET /api/groups/{group_id}/posts ─────────────────────────────────────────

func (h *GroupHandler) GetGroupPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	groupID := r.PathValue("group_id")
	posts, err := h.groupService.GetGroupPosts(groupID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	h.respondJSON(w, http.StatusOK, posts)
}

// ── POST /api/groups/posts/{post_id}/like ─────────────────────────────────────

func (h *GroupHandler) ToggleGroupPostLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	postID := r.PathValue("post_id")
	if postID == "" {
		http.Error(w, "post_id is required", http.StatusBadRequest)
		return
	}
	liked, likeCount, err := h.groupService.ToggleGroupPostLike(postID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]any{
		"liked":      liked,
		"like_count": likeCount,
	})
}

// ── DELETE /api/groups/posts/{post_id} ───────────────────────────────────────

func (h *GroupHandler) DeleteGroupPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	postID := r.PathValue("post_id")
	if err := h.groupService.DeleteGroupPost(postID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]string{"message": "post deleted"})
}

// ── POST /api/groups/posts/{post_id}/comments ─────────────────────────────────

func (h *GroupHandler) CreateGroupComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	postID := r.PathValue("post_id")
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}
	content := r.FormValue("content")
	comment, err := h.groupService.CreateGroupComment(postID, userID, content, "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.respondJSON(w, http.StatusCreated, comment)
}

// ── GET /api/groups/posts/{post_id}/comments ──────────────────────────────────

func (h *GroupHandler) GetGroupComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	postID := r.PathValue("post_id")
	comments, err := h.groupService.GetGroupComments(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.respondJSON(w, http.StatusOK, comments)
}

// ── POST /api/groups/{group_id}/chat ─────────────────────────────────────────

func (h *GroupHandler) SendGroupMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	groupID := r.PathValue("group_id")
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	msg, err := h.groupService.SendGroupMessage(groupID, userID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	// broadcast to all online group members
	memberIDs, _ := h.groupService.GetMemberIDs(groupID)
	h.hub.SendToGroup(memberIDs, map[string]interface{}{
		"type":    "group_message",
		"message": msg,
	})
	h.respondJSON(w, http.StatusCreated, msg)
}

// ── GET /api/groups/{group_id}/chat ──────────────────────────────────────────

func (h *GroupHandler) GetGroupMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := h.getUserID(r)
	groupID := r.PathValue("group_id")
	messages, err := h.groupService.GetGroupMessages(groupID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	h.respondJSON(w, http.StatusOK, messages)
}
