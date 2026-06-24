package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"social-network/backend/internal/config"
	"social-network/backend/internal/models"
	"social-network/backend/internal/repositories/interfaces"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

type OAuthService struct {
	userRepo    interfaces.UserRepository
	oauthConfig *oauth2.Config
}

func NewOAuthService(userRepo interfaces.UserRepository, cfg *config.Config) *OAuthService {
	return &OAuthService{
		userRepo: userRepo,
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

// GetAuthURL returns the URL to redirect the user to Google's consent screen
func (s *OAuthService) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state)
}

// HandleCallback exchanges the code for a token, fetches the user's Google profile,
// and finds or creates a matching user in our database.
func (s *OAuthService) HandleCallback(code string) (*models.User, error) {
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("failed to exchange code")
	}

	client := s.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, errors.New("failed to fetch user info")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read user info")
	}

	var googleUser GoogleUserInfo
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, errors.New("failed to parse user info")
	}

	if !googleUser.VerifiedEmail {
		return nil, errors.New("google email not verified")
	}

	// check if a user with this email already exists
	existing, err := s.userRepo.GetUserByEmail(googleUser.Email)
	if err == nil && existing != nil {
		return existing, nil
	}

	// create a new user — no password since they log in via Google
	newUser := &models.User{
		ID:          uuid.New().String(),
		Email:       googleUser.Email,
		Password:    "", // no password for OAuth users
		FirstName:   googleUser.GivenName,
		LastName:    googleUser.FamilyName,
		DateOfBirth: "", // Google doesn't provide this — user can fill it in later via profile
		Avatar:      googleUser.Picture,
		IsPublic:    true,
	}

	if err := s.userRepo.CreateUser(newUser); err != nil {
		return nil, errors.New("failed to create user")
	}

	return newUser, nil
}
