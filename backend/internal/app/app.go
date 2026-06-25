package app

import (
	"log"
	"net/http"
	"os"

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

	// Ensure upload directories exist
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
	messageRepo := repoSqlite.NewMessageRepository(db)
	groupRepo := repoSqlite.NewGroupRepository(db)
	eventRepo := repoSqlite.NewEventRepository(db)

	// 4. services
	userService := services.NewUserService(userRepo, followerRepo)
	sessionService := services.NewSessionService(sessionRepo)
	followerService := services.NewFollowerService(followerRepo, userRepo, notificationRepo)
	followerService.SetHub(hub)
	postService := services.NewPostService(postRepo, followerRepo, notificationRepo, userRepo)
	notificationService := services.NewNotificationService(notificationRepo)
	messageService := services.NewMessageService(messageRepo, followerRepo)
	groupService := services.NewGroupService(groupRepo)
	eventService := services.NewEventService(eventRepo, groupRepo)

	// 5. handlers
	authHandler := handlers.NewAuthHandler(userService, sessionService)
	followerHandler := handlers.NewFollowerHandler(followerService)
	postHandler := handlers.NewPostHandler(postService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	groupHandler := handlers.NewGroupHandler(groupService)
	eventHandler := handlers.NewEventHandler(eventService)

	// 6. routes
	mux := http.NewServeMux()
	routes.Register(mux, authHandler, followerHandler, postHandler, notificationHandler, sessionService)
	routes.RegisterGroupRoutes(mux, groupHandler, eventHandler, sessionService)
	routes.RegisterWSRoutes(mux, hub, messageService, sessionService)

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
