package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/backend/internal/middleware"
	"social-network/backend/internal/services"
)

type GroupHandler struct {
	groupService *services.GroupService
}

func NewGroupHandler(groupService *services.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

func (h *GroupHandler) jsonOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// ── Core ─────────────────────────────────────────────────────────────────────

// POST /api/groups
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := middleware.GetUserID(r)
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	group, err := h.groupService.CreateGroup(userID, body.Title, body.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	h.jsonOK(w, group)
}

// GET /api/groups
func (h *GroupHandler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := middleware.GetUserID(r)
	groups, err := h.groupService.GetAllGroups(userID)
	if err != nil {
		http.Error(w, "could not fetch groups", http.StatusInternalServerError)
		return
	}
	h.jsonOK(w, groups)
}

// GET /api/groups/{group_id}
func (h *GroupHandler) GetGroupByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	groupID := r.PathValue("group_id")
	group, err := h.groupService.GetGroupByID(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.jsonOK(w, group)
}

// GET /api/groups/{group_id}/members
func (h *GroupHandler) GetMembers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	groupID := r.PathValue("group_id")
	members, err := h.groupService.GetMembers(groupID)
	if err != nil {
		http.Error(w, "could not fetch members", http.StatusInternalServerError)
		return
	}
	h.jsonOK(w, members)
}

// ── Invitations ──────────────────────────────────────────────────────────────

// POST /api/groups/{group_id}/invite
func (h *GroupHandler) InviteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	inviterID := middleware.GetUserID(r)
	groupID := r.PathValue("group_id")
	var body struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.groupService.InviteUser(groupID, inviterID, body.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.jsonOK(w, map[string]string{"message": "invitation sent"})
}

// POST /api/groups/invitations/{inv_id}/accept
func (h *GroupHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := middleware.GetUserID(r)
	invID := r.PathValue("inv_id")
	if err := h.groupService.AcceptInvitation(invID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.jsonOK(w, map[string]string{"message": "invitation accepted"})
}

// POST /api/groups/invitations/{inv_id}/reject
func (h *GroupHandler) RejectInvitation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := middleware.GetUserID(r)
	invID := r.PathValue("inv_id")
	if err := h.groupService.RejectInvitation(invID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.jsonOK(w, map[string]string{"message": "invitation rejected"})
}

// ── Join requests ────────────────────────────────────────────────────────────

// POST /api/groups/{group_id}/join
func (h *GroupHandler) RequestToJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := middleware.GetUserID(r)
	groupID := r.PathValue("group_id")
	if err := h.groupService.RequestToJoin(groupID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.jsonOK(w, map[string]string{"message": "join request sent"})
}

// POST /api/groups/requests/{req_id}/accept
func (h *GroupHandler) AcceptJoinRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	adminID := middleware.GetUserID(r)
	reqID := r.PathValue("req_id")
	if err := h.groupService.AcceptJoinRequest(reqID, adminID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.jsonOK(w, map[string]string{"message": "join request accepted"})
}

// POST /api/groups/requests/{req_id}/reject
func (h *GroupHandler) RejectJoinRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	adminID := middleware.GetUserID(r)
	reqID := r.PathValue("req_id")
	if err := h.groupService.RejectJoinRequest(reqID, adminID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.jsonOK(w, map[string]string{"message": "join request rejected"})
}

// ── Group posts & comments ─────────────────────────────────────────────────

// POST /api/groups/{group_id}/posts
func (h *GroupHandler) CreateGroupPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := middleware.GetUserID(r)
	groupID := r.PathValue("group_id")
	var body struct {
		Content   string `json:"content"`
		MediaPath string `json:"media_path"`
		MediaType string `json:"media_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	post, err := h.groupService.CreateGroupPost(groupID, userID, body.Content, body.MediaPath, body.MediaType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	h.jsonOK(w, post)
}

// GET /api/groups/{group_id}/posts
func (h *GroupHandler) GetGroupPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := middleware.GetUserID(r)
	groupID := r.PathValue("group_id")
	posts, err := h.groupService.GetGroupPosts(groupID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	h.jsonOK(w, posts)
}

// POST /api/groups/posts/{post_id}/comments
func (h *GroupHandler) CreateGroupComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := middleware.GetUserID(r)
	postID := r.PathValue("post_id")
	var body struct {
		Content   string `json:"content"`
		MediaPath string `json:"media_path"`
		MediaType string `json:"media_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	comment, err := h.groupService.CreateGroupComment(postID, userID, body.Content, body.MediaPath, body.MediaType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	h.jsonOK(w, comment)
}

// GET /api/groups/posts/{post_id}/comments
func (h *GroupHandler) GetGroupComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := middleware.GetUserID(r)
	postID := r.PathValue("post_id")
	comments, err := h.groupService.GetGroupComments(postID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	h.jsonOK(w, comments)
}
