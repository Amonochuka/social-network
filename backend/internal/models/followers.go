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

type FollowerProfile struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
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
	CommentCount int       `json:"comment_count"`
	LikeCount    int       `json:"like_count"`
	LikedByMe    bool      `json:"liked_by_me"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostAllowedUser struct {
	PostID string
	UserID string
}

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	MediaPath string    `json:"media_path"`
	MediaType string    `json:"media_type"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentDetail struct {
	ID          string    `json:"id"`
	PostID      string    `json:"post_id"`
	UserID      string    `json:"user_id"`
	AuthorName  string    `json:"author_name"`
	AuthorAvatar string   `json:"author_avatar"`
	Content     string    `json:"content"`
	MediaPath   string    `json:"media_path"`
	MediaType   string    `json:"media_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type Notification struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ActorID     string    `json:"actor_id"`
	Type        string    `json:"type"`
	ReferenceID string    `json:"reference_id"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

type NotificationDetail struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ActorID     string    `json:"actor_id"`
	ActorName   string    `json:"actor_name"`
	ActorAvatar string    `json:"actor_avatar"`
	Type        string    `json:"type"`
	ReferenceID string    `json:"reference_id"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

type FollowStatus struct {
    Status string `json:"status"`
}