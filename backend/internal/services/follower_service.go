package services

import (
	"errors"
	"time"

	"social-network/backend/internal/models"
	"social-network/backend/internal/repositories/interfaces"

	"github.com/google/uuid"
)

type FollowerService struct {
	followerRepo     interfaces.FollowerRepository
	userRepo         interfaces.UserRepository
	notificationRepo interfaces.NotificationRepository
}

func NewFollowerService(
	followerRepo interfaces.FollowerRepository,
	userRepo interfaces.UserRepository,
	notificationRepo interfaces.NotificationRepository,
) *FollowerService {
	return &FollowerService{
		followerRepo:     followerRepo,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *FollowerService) SendFollowRequest(senderID, receiverID string) error {
	// Rule 1 — does the receiver exist?
	_, err := s.userRepo.GetUserByID(receiverID)
	if err != nil {
		return errors.New("receiver not found")
	}

	// Rule 2 — cannot follow yourself
	if senderID == receiverID {
		return errors.New("cannot follow yourself")
	}

	// Rule 3 — already pending?
	existing, _ := s.followerRepo.GetFollowRequest(senderID, receiverID)
	if existing != nil && existing.Status == models.StatusPending {
		return errors.New("follow request already sent")
	}

	// Rule 4 — already following?
	already, err := s.followerRepo.IsFollowing(senderID, receiverID)
	if err != nil {
		return errors.New("could not check follow status")
	}
	if already {
		return errors.New("already following this user")
	}

	// Fetch receiver (privacy check)
	receiver, err := s.userRepo.GetUserByID(receiverID)
	if err != nil {
		return errors.New("receiver not found")
	}

	// Public profile → direct follow
	if receiver.IsPublic {
		if err := s.followerRepo.CreateFollower(senderID, receiverID); err != nil {
			return errors.New("could not follow user")
		}
		return nil
	}

	// Private profile → follow request
	req := &models.FollowRequest{
		ID:         uuid.New().String(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     models.StatusPending,
		CreatedAt:  time.Now(),
	}

	if err := s.followerRepo.CreateFollowRequest(req); err != nil {
		return errors.New("could not send follow request")
	}

	// Notification (DB only, NO websocket here)
	notification := &models.Notification{
		ID:          uuid.New().String(),
		UserID:      receiverID,
		ActorID:     senderID,
		Type:        "follow_request",
		ReferenceID: req.ID,
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	if err := s.notificationRepo.CreateNotification(notification); err != nil {
		return errors.New("could not create notification")
	}

	return nil
}

func (s *FollowerService) AcceptFollowRequest(requestID, receiverID string) error {
	req, err := s.followerRepo.GetFollowRequestByID(requestID)
	if err != nil {
		return errors.New("follow request not found")
	}

	if req.Status != models.StatusPending {
		return errors.New("follow request is no longer pending")
	}

	if req.ReceiverID != receiverID {
		return errors.New("unauthorized")
	}

	if err := s.followerRepo.UpdateFollowRequest(requestID, models.StatusAccepted); err != nil {
		return errors.New("could not update follow request status")
	}

	if err := s.followerRepo.CreateFollower(req.SenderID, req.ReceiverID); err != nil {
		return errors.New("could not create follower relationship")
	}

	notification := &models.Notification{
		ID:          uuid.New().String(),
		UserID:      req.SenderID,
		ActorID:     receiverID,
		Type:        "follow_accepted",
		ReferenceID: req.ID,
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	if err := s.notificationRepo.CreateNotification(notification); err != nil {
		return errors.New("could not create notification")
	}

	return nil
}

func (s *FollowerService) DeclineFollowRequest(requestID, receiverID string) error {
	req, err := s.followerRepo.GetFollowRequestByID(requestID)
	if err != nil {
		return errors.New("follow request not found")
	}

	if req.Status != models.StatusPending {
		return errors.New("follow request is no longer pending")
	}

	if req.ReceiverID != receiverID {
		return errors.New("unauthorized")
	}

	if err := s.followerRepo.UpdateFollowRequest(requestID, models.StatusDeclined); err != nil {
		return errors.New("could not update follow request status")
	}

	notification := &models.Notification{
		ID:          uuid.New().String(),
		UserID:      req.SenderID,
		ActorID:     receiverID,
		Type:        "follow_declined",
		ReferenceID: req.ID,
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	if err := s.notificationRepo.CreateNotification(notification); err != nil {
		return errors.New("could not create notification")
	}

	return nil
}

func (s *FollowerService) Unfollow(followerID, followingID string) error {
	already, err := s.followerRepo.IsFollowing(followerID, followingID)
	if err != nil {
		return errors.New("could not check follow status")
	}
	if !already {
		return errors.New("you are not following this user")
	}
	return s.followerRepo.DeleteFollower(followerID, followingID)
}

func (s *FollowerService) GetFollowers(userID string) ([]*models.FollowerProfile, error) {
	return s.followerRepo.GetFollowers(userID)
}

func (s *FollowerService) GetFollowing(userID string) ([]*models.FollowerProfile, error) {
	return s.followerRepo.GetFollowing(userID)
}