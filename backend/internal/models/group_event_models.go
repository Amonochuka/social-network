package models

import "time"

// ── Groups ───────────────────────────────────────────────────────────────────

type Group struct {
	ID          string    `json:"id"`
	CreatorID   string    `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	MemberCount int       `json:"member_count,omitempty"`
	IsMember    bool      `json:"is_member,omitempty"`
}

type GroupMember struct {
	GroupID  string    `json:"group_id"`
	UserID   string    `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
}

type GroupInvitation struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	InviterID string    `json:"inviter_id"`
	InviteeID string    `json:"invitee_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupRequest struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ── Group Posts & Comments ───────────────────────────────────────────────────

type GroupPost struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	MediaPath string    `json:"media_path,omitempty"`
	MediaType string    `json:"media_type,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupComment struct {
	ID          string    `json:"id"`
	GroupPostID string    `json:"group_post_id"`
	UserID      string    `json:"user_id"`
	Content     string    `json:"content"`
	MediaPath   string    `json:"media_path,omitempty"`
	MediaType   string    `json:"media_type,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// ── Events ───────────────────────────────────────────────────────────────────

type Event struct {
	ID          string    `json:"id"`
	GroupID     string    `json:"group_id"`
	CreatorID   string    `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	CreatedAt   time.Time `json:"created_at"`
	GoingCount  int       `json:"going_count,omitempty"`
	NotGoing    int       `json:"not_going_count,omitempty"`
	MyResponse  string    `json:"my_response,omitempty"` // "going", "not_going", or ""
}

type EventResponse struct {
	EventID   string    `json:"event_id"`
	UserID    string    `json:"user_id"`
	Response  string    `json:"response"` // "going" | "not_going"
	CreatedAt time.Time `json:"created_at"`
}
