package models

import "time"

// ── Groups ──────────────────────────────────────────────────────────────────

type Group struct {
	ID          string    `json:"id"`
	CreatorID   string    `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
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

// ── Events ──────────────────────────────────────────────────────────────────

type Event struct {
	ID          string    `json:"id"`
	GroupID     string    `json:"group_id"`
	CreatorID   string    `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	CreatedAt   time.Time `json:"created_at"`
}

type EventResponse struct {
	EventID   string    `json:"event_id"`
	UserID    string    `json:"user_id"`
	Response  string    `json:"response"` // "going" | "not_going"
	CreatedAt time.Time `json:"created_at"`
}

type EventDetail struct {
	ID           string    `json:"id"`
	GroupID      string    `json:"group_id"`
	CreatorID    string    `json:"creator_id"`
	CreatorName  string    `json:"creator_name"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	EventTime    time.Time `json:"event_time"`
	CreatedAt    time.Time `json:"created_at"`
	GoingCount   int       `json:"going_count"`
	NotGoingCount int      `json:"not_going_count"`
	MyResponse   string    `json:"my_response"` // empty if not yet responded
}

// ── Group Posts & Comments ───────────────────────────────────────────────────

type GroupPost struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	MediaPath string    `json:"media_path"`
	MediaType string    `json:"media_type"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupPostDetail struct {
	ID           string    `json:"id"`
	GroupID      string    `json:"group_id"`
	UserID       string    `json:"user_id"`
	AuthorName   string    `json:"author_name"`
	AuthorAvatar string    `json:"author_avatar"`
	Content      string    `json:"content"`
	MediaPath    string    `json:"media_path"`
	MediaType    string    `json:"media_type"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type GroupComment struct {
	ID          string    `json:"id"`
	GroupPostID string    `json:"group_post_id"`
	UserID      string    `json:"user_id"`
	Content     string    `json:"content"`
	MediaPath   string    `json:"media_path"`
	MediaType   string    `json:"media_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupCommentDetail struct {
	ID           string    `json:"id"`
	GroupPostID  string    `json:"group_post_id"`
	UserID       string    `json:"user_id"`
	AuthorName   string    `json:"author_name"`
	AuthorAvatar string    `json:"author_avatar"`
	Content      string    `json:"content"`
	MediaPath    string    `json:"media_path"`
	MediaType    string    `json:"media_type"`
	CreatedAt    time.Time `json:"created_at"`
}




type MemberProfile struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
}

type GroupDetail struct {
	Group
	Members       []*MemberProfile `json:"members"`
	IsMember      bool            `json:"is_member"`
	IsCreator     bool            `json:"is_creator"`
	InviteStatus  string          `json:"invite_status"`  // "pending" | "accepted" | ""
	RequestStatus string          `json:"request_status"` // "pending" | "accepted" | ""
}
