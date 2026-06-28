package services

import (
	"errors"
	"social-network/backend/internal/models"
	"social-network/backend/internal/repositories/interfaces"
)

type ChatService struct {
	chatRepo     interfaces.ChatRepository
	followerRepo interfaces.FollowerRepository
	userRepo     interfaces.UserRepository
}

func NewChatService(chatRepo interfaces.ChatRepository, followerRepo interfaces.FollowerRepository, userRepo interfaces.UserRepository) *ChatService {
	return &ChatService{
		chatRepo:     chatRepo,
		followerRepo: followerRepo,
		userRepo:     userRepo,
	}
}

// CanChat checks the spec's rule: at least one of the two users must follow the other.
func (s *ChatService) CanChat(userA, userB string) (bool, error) {
	aFollowsB, err := s.followerRepo.IsFollowing(userA, userB)
	if err != nil {
		return false, err
	}
	if aFollowsB {
		return true, nil
	}

	bFollowsA, err := s.followerRepo.IsFollowing(userB, userA)
	if err != nil {
		return false, err
	}
	return bFollowsA, nil
}

func (s *ChatService) SendPrivateMessage(senderID, receiverID, content string) (*models.PrivateMessage, error) {
	if senderID == receiverID {
		return nil, errors.New("cannot message yourself")
	}

	canChat, err := s.CanChat(senderID, receiverID)
	if err != nil {
		return nil, errors.New("could not verify chat permission")
	}
	if !canChat {
		return nil, errors.New("you must follow or be followed by this user to chat")
	}

	msg := &models.PrivateMessage{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}

	if err := s.chatRepo.CreatePrivateMessage(msg); err != nil {
		return nil, errors.New("could not send message")
	}

	return msg, nil
}

func (s *ChatService) GetPrivateMessages(userA, userB string) ([]*models.PrivateMessageDetail, error) {
	canChat, err := s.CanChat(userA, userB)
	if err != nil {
		return nil, errors.New("could not verify chat permission")
	}
	if !canChat {
		return nil, errors.New("you must follow or be followed by this user to view this chat")
	}

	messages, err := s.chatRepo.GetPrivateMessages(userA, userB)
	if err != nil {
		return nil, errors.New("could not get messages")
	}
	return messages, nil
}

func (s *ChatService) GetConversations(userID string) ([]*models.ConversationPreview, error) {
	conversations, err := s.chatRepo.GetConversations(userID)
	if err != nil {
		return nil, errors.New("could not get conversations")
	}
	return conversations, nil
}

type ChatPartnerInfo struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

// GetPartnerInfo returns basic profile info for the chat header.
// This bypasses profile privacy intentionally — visibility for chat purposes
// is governed by CanChat (the follow-based rule), not the profile privacy setting.
// Callers must call CanChat themselves before calling this.
// Never returns the full User struct — that would leak the password hash.
func (s *ChatService) GetPartnerInfo(userID string) (*ChatPartnerInfo, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &ChatPartnerInfo{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
	}, nil
}
