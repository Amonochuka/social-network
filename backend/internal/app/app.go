package app

import (
	"log"
	"net/http"
	"os"
	"social-network/backend/internal/config"
	"social-network/backend/internal/handlers"
	"social-network/backend/internal/middleware"
	repoSqlite "social-network/backend/internal/repositories/sqlite"
	"social-network/backend/internal/routes"
	"social-network/backend/internal/services"
	"social-network/backend/internal/ws"
	dbSqlite "social-network/backend/pkg/db/sqlite"
)

type App struct {
	Router *http.ServeMux
	Hub    *ws.Hub
}

func New() (*App, error) {
	// 1. connect to DB + run migrations
	db := dbSqlite.NewDB()

	// load env vars (.env file) — needed for Google OAuth credentials
	cfg := config.Load()

	// Ensure upload directory exists for avatars before server starts
	if err := os.MkdirAll("uploads/avatars", 0755); err != nil {
		return nil, err
	}

	// 2. Real-time hub (instantiated early — injected into notification repo)
	hub := ws.NewHub()

	// 3. repos
	userRepo := repoSqlite.NewUserRepository(db)
	sessionRepo := repoSqlite.NewSessionRepository(db)
	followerRepo := repoSqlite.NewFollowerRepository(db)
	notificationRepo := repoSqlite.NewNotificationRepository(db, hub)
	postRepo := repoSqlite.NewPostRepository(db)
	chatRepo := repoSqlite.NewChatRepository(db)
	groupRepo := repoSqlite.NewGroupRepository(db)

	// 4. services
	userService := services.NewUserService(userRepo, followerRepo)
	sessionService := services.NewSessionService(sessionRepo)
	followerService := services.NewFollowerService(followerRepo, userRepo, notificationRepo)
	followerService.SetHub(hub)
	postService := services.NewPostService(postRepo, followerRepo, notificationRepo, userRepo)
	notificationService := services.NewNotificationService(notificationRepo)
	oauthService := services.NewOAuthService(userRepo, cfg)
	chatService := services.NewChatService(chatRepo, followerRepo, userRepo)
	groupService := services.NewGroupService(groupRepo, userRepo, followerRepo, notificationRepo)

	
	// 5. handlers
	authHandler := handlers.NewAuthHandler(userService, sessionService)
	followerHandler := handlers.NewFollowerHandler(followerService)
	postHandler := handlers.NewPostHandler(postService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	oauthHandler := handlers.NewOAuthHandler(oauthService, sessionService)
	chatHandler := handlers.NewChatHandler(chatService, hub)
	groupHandler := handlers.NewGroupHandler(groupService, hub)

	
	// 6. routes
	mux := http.NewServeMux()

	routes.Register(mux, authHandler, followerHandler, postHandler, notificationHandler, oauthHandler, chatHandler, groupHandler, sessionService)
	log.Println("app initialised")
	return &App{Router: mux, Hub: hub}, nil
}

func (a *App) ChainMiddlewares() http.Handler {
	nextMiddleware := middleware.ChainMiddlewares(
		a.Router,
		middleware.CORSMiddleware,
	)
	return nextMiddleware
}

