package models

import "time"

const (
	StatusPending  = "pending"
	StatusAccepted = "accepted"
	StatusDeclined = "declined"
)

type Follower struct {
	FollowerID  string
	FollowingID string
	CreatedAt   time.Time
}

type FollowRequest struct {
	ID         string
	SenderID   string
	ReceiverID string
	Status     string
	CreatedAt  time.Time
}

type Post struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	MediaPath string    `json:"media_path"`
	MediaType string    `json:"media_type"`
	Privacy   string    `json:"privacy"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FeedPost struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	AuthorName   string    `json:"author_name"`
	AuthorAvatar string    `json:"author_avatar"`
	Content      string    `json:"content"`
	MediaPath    string    `json:"media_path"`
	MediaType    string    `json:"media_type"`
	Privacy      string    `json:"privacy"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostAllowedUser struct {
	PostID string
	UserID string
}

type Comment struct {
	ID        string
	PostID    string
	UserID    string
	Content   string
	MediaPath string
	MediaType string
	CreatedAt time.Time
}

type Notification struct {
	ID          string
	UserID      string
	ActorID     string
	Type        string
	ReferenceID string
	IsRead      bool
	CreatedAt   time.Time
}