package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/internal/handlers"
	"social-network/backend/internal/middleware"
	"social-network/backend/internal/models"
	"social-network/backend/internal/services"
	"social-network/backend/internal/ws"
)

func Register(
	mux *http.ServeMux,
	authHandler *handlers.AuthHandler,
	followerHandler *handlers.FollowerHandler,
	postHandler *handlers.PostHandler,
	notificationHandler *handlers.NotificationHandler,
	sessionService *services.SessionService,
) {
	auth := middleware.NewAuthMiddleware(sessionService)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Public Auth Routes
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/logout", authHandler.Logout)

	// Private Profile Routes
	mux.Handle("/api/auth/me", auth.Authenticate(http.HandlerFunc(authHandler.Me)))
	mux.Handle("/api/profile/{id}", auth.Authenticate(http.HandlerFunc(authHandler.GetProfile)))
	mux.Handle("/api/profile", auth.Authenticate(http.HandlerFunc(authHandler.UpdateProfile)))
	mux.Handle("/api/profile/privacy", auth.Authenticate(http.HandlerFunc(authHandler.UpdatePrivacy)))
	mux.Handle("/api/profile/avatar", auth.Authenticate(http.HandlerFunc(authHandler.UploadAvatar)))

	// Private Follower Routes
	mux.Handle("/api/follow/requests", auth.Authenticate(http.HandlerFunc(followerHandler.SendFollowRequest)))
	mux.Handle("/api/follow/requests/{request_id}/accept", auth.Authenticate(http.HandlerFunc(followerHandler.AcceptFollowRequest)))
	mux.Handle("/api/follow/requests/{request_id}/decline", auth.Authenticate(http.HandlerFunc(followerHandler.DeclineFollowRequest)))
	mux.Handle("/api/follow/{user_id}", auth.Authenticate(http.HandlerFunc(followerHandler.Unfollow)))
	mux.Handle("/api/followers", auth.Authenticate(http.HandlerFunc(followerHandler.GetFollowers)))
	mux.Handle("/api/following", auth.Authenticate(http.HandlerFunc(followerHandler.GetFollowing)))

	// Private Post Routes
	mux.Handle("/api/posts", auth.Authenticate(http.HandlerFunc(postHandler.CreatePost)))
	mux.Handle("/api/posts/feed", auth.Authenticate(http.HandlerFunc(postHandler.GetFeed)))
	mux.Handle("/api/posts/{post_id}", auth.Authenticate(http.HandlerFunc(postHandler.GetPostByID)))
	mux.Handle("/api/posts/{post_id}/update", auth.Authenticate(http.HandlerFunc(postHandler.UpdatePost)))
	mux.Handle("/api/posts/{post_id}/delete", auth.Authenticate(http.HandlerFunc(postHandler.DeletePost)))
	mux.Handle("/api/users/{user_id}/posts", auth.Authenticate(http.HandlerFunc(postHandler.GetPostsByUserID)))
	mux.Handle("/api/posts/{post_id}/comments", auth.Authenticate(http.HandlerFunc(postHandler.CreateComment)))
	mux.Handle("/api/posts/{post_id}/comments/all", auth.Authenticate(http.HandlerFunc(postHandler.GetCommentsByPostID)))

	// Private Notification Routes
	mux.Handle("/api/notifications", auth.Authenticate(http.HandlerFunc(notificationHandler.GetNotifications)))
	mux.Handle("/api/notifications/{notification_id}/read", auth.Authenticate(http.HandlerFunc(notificationHandler.MarkAsRead)))
}

func RegisterWSRoutes(
	mux *http.ServeMux,
	hub *ws.Hub,
	msgSvc *services.MessageService,
	sessionService *services.SessionService,
) {
	auth := middleware.NewAuthMiddleware(sessionService)

	mux.Handle("/ws", auth.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		if userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ws.ServeWS(w, r,
			func(client *ws.Client) {
				hub.RegisterPrivate(userID, client)
			},
			func() {
				hub.UnregisterPrivate(userID)
			},
			func(msg []byte) {
				// 1. Parse the outer framework shell
				var frame models.WSMessageFrame
				if err := json.Unmarshal(msg, &frame); err != nil {
					log.Printf("Malformed WS envelope error: %v", err)
					return
				}

				// 2. Route payload according to type specification
				switch frame.Type {
				case "private_chat":
					var chat models.PrivateChatPayload
					if err := json.Unmarshal(frame.Payload, &chat); err != nil {
						log.Printf("Invalid private payload structure: %v", err)
						return
					}

					// Persist to database utilizing your business rule validation
					savedMsg, err := msgSvc.SendPrivateMessage(userID, chat.ReceiverID, chat.Content)
					if err != nil {
						log.Printf("Business logic chat rejection: %v", err)
						return
					}

					// Serialize saved record back to JSON for transmission
					outboundPayload, _ := json.Marshal(map[string]interface{}{
						"type":    "private_chat",
						"payload": savedMsg,
					})

					// Deliver to recipient real-time via Hub registry
					hub.SendToUser(chat.ReceiverID, outboundPayload)
					// Echo back to sender to confirm successful delivery status
					hub.SendToUser(userID, outboundPayload)

				case "group_chat":
					var group models.GroupChatPayload
					if err := json.Unmarshal(frame.Payload, &group); err != nil {
						log.Printf("Invalid group payload structure: %v", err)
						return
					}

					// Persist to database utilizing group validation criteria
					savedMsg, err := msgSvc.SendGroupMessage(userID, group.GroupID, group.Content)
					if err != nil {
						log.Printf("Business logic group rejection: %v", err)
						return
					}

					outboundPayload, _ := json.Marshal(map[string]interface{}{
						"type":    "group_chat",
						"payload": savedMsg,
					})

					// Dispatch transmission dynamically to all active members in the room
					hub.BroadcastToGroup(group.GroupID, outboundPayload)

				default:
					log.Printf("Unhandled WebSocket transmission packet frame type: %s", frame.Type)
				}
			},
		)
	})))
}
