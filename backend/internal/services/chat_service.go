package services

import (
	"errors"
	"social-network/backend/internal/models"
	"social-network/backend/internal/repositories/interfaces"
)

type ChatService struct {
	chatRepo     interfaces.ChatRepository
	followerRepo interfaces.FollowerRepository
}

func NewChatService(chatRepo interfaces.ChatRepository, followerRepo interfaces.FollowerRepository) *ChatService {
	return &ChatService{
		chatRepo:     chatRepo,
		followerRepo: followerRepo,
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