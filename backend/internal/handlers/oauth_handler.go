package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/backend/internal/models"
	"social-network/backend/internal/services"
)

type OAuthHandler struct {
	oauthService   *services.OAuthService
	sessionService *services.SessionService
}

func NewOAuthHandler(oauthService *services.OAuthService, sessionService *services.SessionService) *OAuthHandler {
	return &OAuthHandler{
		oauthService:   oauthService,
		sessionService: sessionService,
	}
}

// GoogleLogin redirects the user to Google's consent screen
func (h *OAuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.oauthService.GetAuthURL("state-token") // in production, generate and verify a random state per request
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handles Google's redirect back with the auth code
func (h *OAuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	user, err := h.oauthService.HandleCallback(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := h.sessionService.CreateSession(user.ID)
	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})

	// redirect back to the frontend after successful login
	http.Redirect(w, r, "http://localhost:3000", http.StatusTemporaryRedirect)
}

// Helper used only if you want a JSON response instead of redirect (not used by default flow above)
func writeUserJSON(w http.ResponseWriter, user *models.User) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DateOfBirth: user.DateOfBirth,
		Avatar:      user.Avatar,
		NickName:    user.NickName,
		AboutMe:     user.AboutMe,
		IsPublic:    user.IsPublic,
	})
}