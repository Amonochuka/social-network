package middleware

import (
	"context"
	"net/http"
	"social-network/backend/internal/services"
)

type contextKey string

const UserIDKey contextKey = "user_id"


type AuthMiddleware struct {
	sessionService *services.SessionService
}

func NewAuthMiddleware(sessionService *services.SessionService) *AuthMiddleware {
	return &AuthMiddleware{
		sessionService: sessionService,
	}
}

// Authenticate is a middleware that validates the session cookie
func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		session, err := m.sessionService.GetSession(cookie.Value)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
		next(w, r.WithContext(ctx))
	}
}



// OptionalAuth is a middleware that optionally validates the session
func (m *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			session, err := m.sessionService.GetSession(cookie.Value)
			if err == nil {
				ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
				next(w, r.WithContext(ctx))
				return
			}
		}
		next(w, r)
	}
}



// LoggingMiddleware logs incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request
		println("Request:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func GetUserID(r *http.Request) string {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}