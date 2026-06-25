package interfaces

import "social-network/backend/internal/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(user *models.User) error
	UpdateAvatar(userID, avatarPath string) error
	UpdatePrivacy(userID string, isPublic bool) error
}

type SessionRepository interface {
	CreateSession(session *models.Session) error
	GetSessionByID(id string) (*models.Session, error)
	DeleteSession(id string) error
}

type FollowerRepository interface {
	// follow requests
	CreateFollowRequest(req *models.FollowRequest) error
	GetFollowRequest(senderID, receiverID string) (*models.FollowRequest, error)
	GetFollowRequestByID(requestID string) (*models.FollowRequest, error)
	UpdateFollowRequest(requestID, status string) error
	// follow relationships
	CreateFollower(followerID, followingID string) error
	DeleteFollower(followerID, followingID string) error
	IsFollowing(followerID, followingID string) (bool, error)
	// lists
	GetFollowers(userID string) ([]*models.Follower, error)
	GetFollowing(userID string) ([]*models.Follower, error)
}

type NotificationRepository interface {
	CreateNotification(n *models.Notification) error
	GetNotificationsByUserID(userID string) ([]*models.Notification, error)
	MarkNotificationAsRead(notificationID, userID string) error
}

type PostRepository interface {
	// posts
	CreatePost(post *models.Post) error
	CreateAllowedUser(postID, userID string) error
	GetPostsByUserID(userID, viewerID string) ([]*models.Post, error)
	GetPostByID(postID string) (*models.Post, error)
	GetFeed(userID string) ([]*models.Post, error)
	UpdatePost(post *models.Post) error
	DeletePost(postID string) error
	IsAllowedToViewPost(postID, userID string) (bool, error)
	// comments
	CreateComment(comment *models.Comment) error
	GetCommentsByPostID(postID string) ([]*models.Comment, error)
}

// MessageRepository handles persistence for both private and group messages.
type MessageRepository interface {
	// Private messages
	SavePrivateMessage(msg *models.PrivateMessage) error
	GetPrivateMessages(userA, userB string, limit int) ([]*models.PrivateMessage, error)

	// Group messages
	SaveGroupMessage(msg *models.GroupMessage) error
	GetGroupMessages(groupID string, limit int) ([]*models.GroupMessage, error)

	// Group membership check (groups table owned by Dev 4)
	IsGroupMember(userID, groupID string) (bool, error)
}

// GroupRepository covers all group and group-content persistence.
type GroupRepository interface {
	// ── Core group CRUD ──────────────────────────────────────────────────────
	CreateGroup(g *models.Group) error
	GetGroupByID(groupID string) (*models.Group, error)
	GetAllGroups(viewerID string) ([]*models.Group, error)

	// ── Membership ───────────────────────────────────────────────────────────
	AddMember(groupID, userID string) error
	RemoveMember(groupID, userID string) error
	IsMember(groupID, userID string) (bool, error)
	GetMembers(groupID string) ([]*models.GroupMember, error)

	// ── Invitations ──────────────────────────────────────────────────────────
	CreateInvitation(inv *models.GroupInvitation) error
	GetInvitationByID(invID string) (*models.GroupInvitation, error)
	GetPendingInvitation(groupID, inviteeID string) (*models.GroupInvitation, error)
	UpdateInvitation(invID, status string) error

	// ── Join requests ────────────────────────────────────────────────────────
	CreateJoinRequest(req *models.GroupRequest) error
	GetJoinRequestByID(reqID string) (*models.GroupRequest, error)
	GetPendingJoinRequest(groupID, userID string) (*models.GroupRequest, error)
	UpdateJoinRequest(reqID, status string) error

	// ── Group posts ──────────────────────────────────────────────────────────
	CreateGroupPost(post *models.GroupPost) error
	GetGroupPosts(groupID string) ([]*models.GroupPost, error)
	GetGroupPostByID(postID string) (*models.GroupPost, error)

	// ── Group comments ───────────────────────────────────────────────────────
	CreateGroupComment(comment *models.GroupComment) error
	GetGroupComments(groupPostID string) ([]*models.GroupComment, error)
}

// EventRepository covers event persistence.
type EventRepository interface {
	CreateEvent(e *models.Event) error
	GetEventsByGroupID(groupID, viewerID string) ([]*models.Event, error)
	GetEventByID(eventID, viewerID string) (*models.Event, error)
	UpsertResponse(resp *models.EventResponse) error
}
