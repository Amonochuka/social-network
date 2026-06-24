package services

import (
	"errors"
	"social-network/backend/internal/models"
	"social-network/backend/internal/repositories/interfaces"
)

// MessageService enforces business rules before touching the database.
type MessageService struct {
	msgRepo      interfaces.MessageRepository
	followerRepo interfaces.FollowerRepository // reuse your existing interface
}

func NewMessageService(
	msgRepo interfaces.MessageRepository,
	followerRepo interfaces.FollowerRepository,
) *MessageService {
	return &MessageService{
		msgRepo:      msgRepo,
		followerRepo: followerRepo,
	}
}

// ── Private chat ─────────────────────────────────────────────────────────────

// CanChat returns true if at least one of the two users follows the other.
func (s *MessageService) CanChat(userA, userB string) (bool, error) {
	aFollowsB, err := s.followerRepo.IsFollowing(userA, userB)
	if err != nil {
		return false, err
	}
	if aFollowsB {
		return true, nil
	}
	bFollowsA, err := s.followerRepo.IsFollowing(userB, userA)
	return bFollowsA, err
}

func (s *MessageService) SendPrivateMessage(senderID, receiverID, content string) (*models.PrivateMessage, error) {
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	allowed, err := s.CanChat(senderID, receiverID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errors.New("chat not allowed: no follow relationship exists")
	}

	msg := &models.PrivateMessage{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}
	if err := s.msgRepo.SavePrivateMessage(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *MessageService) GetPrivateHistory(userA, userB string) ([]*models.PrivateMessage, error) {
	return s.msgRepo.GetPrivateMessages(userA, userB, 50)
}

// ── Group chat ───────────────────────────────────────────────────────────────

func (s *MessageService) SendGroupMessage(senderID, groupID, content string) (*models.GroupMessage, error) {
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	isMember, err := s.msgRepo.IsGroupMember(senderID, groupID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user is not a member of this group")
	}

	msg := &models.GroupMessage{
		GroupID:  groupID,
		SenderID: senderID,
		Content:  content,
	}
	if err := s.msgRepo.SaveGroupMessage(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *MessageService) GetGroupHistory(groupID string) ([]*models.GroupMessage, error) {
	return s.msgRepo.GetGroupMessages(groupID, 50)
}
