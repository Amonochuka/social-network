package routes

import (
	"net/http"

	"social-network/backend/internal/handlers"
	"social-network/backend/internal/middleware"
	"social-network/backend/internal/services"
)

func Register(
	mux *http.ServeMux,
	authHandler *handlers.AuthHandler,
	followerHandler *handlers.FollowerHandler,
	postHandler *handlers.PostHandler,
	notificationHandler *handlers.NotificationHandler,
	oauthHandler *handlers.OAuthHandler,
	chatHandler *handlers.ChatHandler,
	groupHandler *handlers.GroupHandler,
	searchHandler *handlers.SearchHandler,
	sessionService *services.SessionService,
) {
	auth := middleware.NewAuthMiddleware(sessionService)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// ==========================
	// Public Auth Routes
	// ==========================
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/logout", authHandler.Logout)

	// ==========================
	// Google OAuth
	// ==========================
	mux.HandleFunc("/api/auth/google", oauthHandler.GoogleLogin)
	mux.HandleFunc("/api/auth/google/callback", oauthHandler.GoogleCallback)

	// ==========================
	// Profile Routes
	// ==========================
	mux.Handle("/api/auth/me", auth.Authenticate(http.HandlerFunc(authHandler.Me)))
	mux.Handle("/api/profile/{id}", auth.Authenticate(http.HandlerFunc(authHandler.GetProfile)))
	mux.Handle("/api/profile", auth.Authenticate(http.HandlerFunc(authHandler.UpdateProfile)))
	mux.Handle("/api/profile/privacy", auth.Authenticate(http.HandlerFunc(authHandler.UpdatePrivacy)))
	mux.Handle("/api/profile/avatar", auth.Authenticate(http.HandlerFunc(authHandler.UploadAvatar)))
	mux.Handle("/api/profile/password",
	auth.Authenticate(http.HandlerFunc(authHandler.ChangePassword)))

	// ==========================
	// Search
	// ==========================
	mux.Handle("/api/search", auth.Authenticate(http.HandlerFunc(searchHandler.Search)))

	// ==========================
	// Follow Routes
	// ==========================
	mux.Handle("/api/follow/requests", auth.Authenticate(http.HandlerFunc(followerHandler.SendFollowRequest)))
	mux.Handle("/api/follow/requests/{request_id}/accept", auth.Authenticate(http.HandlerFunc(followerHandler.AcceptFollowRequest)))
	mux.Handle("/api/follow/requests/{request_id}/decline", auth.Authenticate(http.HandlerFunc(followerHandler.DeclineFollowRequest)))
	mux.Handle("/api/follow/{user_id}", auth.Authenticate(http.HandlerFunc(followerHandler.Unfollow)))
	mux.Handle("/api/followers", auth.Authenticate(http.HandlerFunc(followerHandler.GetFollowers)))
	mux.Handle("/api/following", auth.Authenticate(http.HandlerFunc(followerHandler.GetFollowing)))
	mux.Handle("/api/follow/{user_id}/status",auth.Authenticate(http.HandlerFunc(followerHandler.GetFollowStatus)))

	// ==========================
	// Posts
	// ==========================
	mux.Handle("/api/posts", auth.Authenticate(http.HandlerFunc(postHandler.CreatePost)))
	mux.Handle("/api/posts/feed", auth.Authenticate(http.HandlerFunc(postHandler.GetFeed)))
	mux.Handle("/api/posts/{post_id}", auth.Authenticate(http.HandlerFunc(postHandler.GetPostByID)))
	mux.Handle("/api/posts/{post_id}/update", auth.Authenticate(http.HandlerFunc(postHandler.UpdatePost)))
	mux.Handle("/api/posts/{post_id}/delete", auth.Authenticate(http.HandlerFunc(postHandler.DeletePost)))
	mux.Handle("/api/users/{user_id}/posts", auth.Authenticate(http.HandlerFunc(postHandler.GetPostsByUserID)))
	mux.Handle("/api/posts/{post_id}/comments", auth.Authenticate(http.HandlerFunc(postHandler.CreateComment)))
	mux.Handle("/api/posts/{post_id}/comments/all", auth.Authenticate(http.HandlerFunc(postHandler.GetCommentsByPostID)))

	// ==========================
	// Notifications
	// ==========================
	mux.Handle("/api/notifications", auth.Authenticate(http.HandlerFunc(notificationHandler.GetNotifications)))
	mux.Handle("/api/notifications/{notification_id}/read", auth.Authenticate(http.HandlerFunc(notificationHandler.MarkAsRead)))
	mux.Handle("/api/notifications/unread", auth.Authenticate(http.HandlerFunc(notificationHandler.GetUnreadCount)))

	// ==========================
	// Private Chat
	// ==========================
	mux.Handle("/api/chat/private", auth.Authenticate(http.HandlerFunc(chatHandler.SendPrivateMessage)))
	mux.Handle("/api/chat/private/{user_id}", auth.Authenticate(http.HandlerFunc(chatHandler.GetPrivateMessages)))
	mux.Handle("/api/chat/ws", auth.Authenticate(http.HandlerFunc(chatHandler.ServeWebSocket)))
	mux.Handle("/api/chat/conversations", auth.Authenticate(http.HandlerFunc(chatHandler.GetConversations)))
	mux.Handle("/api/chat/partner/{user_id}", auth.Authenticate(http.HandlerFunc(chatHandler.GetChatPartnerInfo)))

	// ==========================
	// Groups
	// ==========================
	mux.Handle("/api/groups", auth.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			groupHandler.CreateGroup(w, r)
		} else {
			groupHandler.ListGroups(w, r)
		}
	})))

	mux.Handle("/api/groups/{group_id}", auth.Authenticate(http.HandlerFunc(groupHandler.GetGroup)))
	mux.Handle("/api/groups/{group_id}/invite", auth.Authenticate(http.HandlerFunc(groupHandler.InviteUser)))
	mux.Handle("/api/groups/{group_id}/join", auth.Authenticate(http.HandlerFunc(groupHandler.RequestToJoin)))

	mux.Handle("/api/groups/{group_id}/events", auth.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			groupHandler.CreateEvent(w, r)
		} else {
			groupHandler.GetEvents(w, r)
		}
	})))

	mux.Handle("/api/groups/events/{event_id}/rsvp", auth.Authenticate(http.HandlerFunc(groupHandler.RSVPEvent)))

	mux.Handle("/api/groups/{group_id}/posts", auth.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			groupHandler.CreateGroupPost(w, r)
		} else {
			groupHandler.GetGroupPosts(w, r)
		}
	})))

	mux.Handle("/api/groups/posts/{post_id}/comments", auth.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			groupHandler.CreateGroupComment(w, r)
		} else {
			groupHandler.GetGroupComments(w, r)
		}
	})))

	mux.Handle("/api/groups/{group_id}/chat", auth.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			groupHandler.SendGroupMessage(w, r)
		} else {
			groupHandler.GetGroupMessages(w, r)
		}
	})))

	mux.Handle("/api/groups/invitations/{inv_id}/accept", auth.Authenticate(http.HandlerFunc(groupHandler.AcceptInvitation)))
	mux.Handle("/api/groups/invitations/{inv_id}/decline", auth.Authenticate(http.HandlerFunc(groupHandler.DeclineInvitation)))
	mux.Handle("/api/groups/requests/{req_id}/accept", auth.Authenticate(http.HandlerFunc(groupHandler.AcceptJoinRequest)))
	mux.Handle("/api/groups/requests/{req_id}/decline", auth.Authenticate(http.HandlerFunc(groupHandler.DeclineJoinRequest)))
}