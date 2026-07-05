package interfaces

import "social-network/backend/internal/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(user *models.User) error
	UpdateAvatar(userID, avatarPath string) error
	UpdatePrivacy(userID string, isPublic bool) error
	SearchUsers(query, currentUserID string) ([]*models.UserSearchResult, error)
	UpdatePassword(userID, hashedPassword string) error
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
	GetFollowers(userID string) ([]*models.FollowerProfile, error)
	GetFollowing(userID string) ([]*models.FollowerProfile, error)
	DeleteFollowRequest(requestID string) error
	GetFollowStatus(viewerID, targetID string) (string, error)
}

type NotificationRepository interface {
	CreateNotification(n *models.Notification) error
	GetNotificationsByUserID(userID string) ([]*models.NotificationDetail, error)
	MarkNotificationAsRead(notificationID, userID string) error
	GetUnreadCount(userID string) (int, error)
	DeleteNotificationByReferenceID(refID, notificationType string) error
}

type PostRepository interface {
	// posts
	CreatePost(post *models.Post) error
	CreateAllowedUser(postID, userID string) error
	GetPostsByUserID(userID, viewerID string) ([]*models.Post, error)
	GetPostByID(postID string) (*models.Post, error)
	GetFeed(userID string) ([]*models.FeedPost, error)
	UpdatePost(post *models.Post) error
	DeletePost(postID string) error
	IsAllowedToViewPost(postID, userID string) (bool, error)
	// comments
	CreateComment(comment *models.Comment) error
	GetCommentsByPostID(postID string) ([]*models.CommentDetail, error)
	// likes
	ToggleLike(postID, userID string) (liked bool, likeCount int, err error)
	HasLiked(postID, userID string) (bool, error)
	CountLikes(postID string) (int, error)
}

type ChatRepository interface {
	// private messages
	CreatePrivateMessage(msg *models.PrivateMessage) error
	GetPrivateMessages(userA, userB string) ([]*models.PrivateMessageDetail, error)
	// group messages
	CreateGroupMessage(msg *models.GroupMessage) error
	GetGroupMessages(groupID string) ([]*models.GroupMessageDetailFull, error)
	GetConversations(userID string) ([]*models.ConversationPreview, error)
}

type GroupRepository interface {
	// groups
	CreateGroup(group *models.Group) error
	GetGroupByID(groupID string) (*models.Group, error)
	GetAllGroups() ([]*models.Group, error)
	// members
	AddMember(groupID, userID string) error
	RemoveMember(groupID, userID string) error
	IsMember(groupID, userID string) (bool, error)
	GetMembers(groupID string) ([]*models.MemberProfile, error)
	GetMemberIDs(groupID string) ([]string, error)
	// invitations
	CreateInvitation(inv *models.GroupInvitation) error
	GetInvitationByID(invID string) (*models.GroupInvitation, error)
	GetPendingInvitation(groupID, inviteeID string) (*models.GroupInvitation, error)
	UpdateInvitationStatus(invID, status string) error
	// join requests
	CreateJoinRequest(req *models.GroupRequest) error
	GetJoinRequestByID(reqID string) (*models.GroupRequest, error)
	GetPendingJoinRequest(groupID, userID string) (*models.GroupRequest, error)
	UpdateJoinRequestStatus(reqID, status string) error
	// events
	CreateEvent(event *models.Event) error
	GetEventsByGroupID(groupID string) ([]*models.Event, error)
	GetEventByID(eventID string) (*models.Event, error)
	UpsertEventResponse(eventID, userID, response string) error
	GetEventResponses(eventID string) (going int, notGoing int, err error)
	GetMyEventResponse(eventID, userID string) (string, error)
	// group posts
	CreateGroupPost(post *models.GroupPost) error
	GetGroupPostByID(postID string) (*models.GroupPost, error)
	GetGroupPosts(groupID, userID string) ([]*models.GroupPostDetail, error)
	DeleteGroupPost(postID string) error
	CreateGroupComment(comment *models.GroupComment) error
	GetGroupComments(groupPostID string) ([]*models.GroupCommentDetail, error)
	ToggleGroupPostLike(postID, userID string) (liked bool, likeCount int, err error)
	HasLikedGroupPost(postID, userID string) (bool, error)
	CountGroupPostLikes(postID string) (int, error)
	// group chat
	CreateGroupChatMessage(msg *models.GroupMessage) error
	GetGroupChatMessages(groupID string) ([]*models.GroupMessageDetailFull, error)
	// invitation status helpers
	GetInvitationStatus(groupID, userID string) (string, error)
	GetJoinRequestStatus(groupID, userID string) (string, error)
	SearchGroups(query string) ([]*models.GroupSearchResult, error)
}