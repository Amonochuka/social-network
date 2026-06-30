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

func (s *FollowerService) SendFollowRequest(senderID, receiverID string) (*models.Notification, error) {
	// Receiver exists?
	_, err := s.userRepo.GetUserByID(receiverID)
	if err != nil {
		return nil, errors.New("receiver not found")
	}

	// Can't follow yourself
	if senderID == receiverID {
		return nil, errors.New("cannot follow yourself")
	}

	// Already pending?
	existing, _ := s.followerRepo.GetFollowRequest(senderID, receiverID)
	if existing != nil && existing.Status == models.StatusPending {
		return nil, errors.New("follow request already sent")
	}

	// Already following?
	already, err := s.followerRepo.IsFollowing(senderID, receiverID)
	if err != nil {
		return nil, errors.New("could not check follow status")
	}
	if already {
		return nil, errors.New("already following this user")
	}

	receiver, err := s.userRepo.GetUserByID(receiverID)
	if err != nil {
		return nil, errors.New("receiver not found")
	}

	// PUBLIC ACCOUNT
	if receiver.IsPublic {
		if err := s.followerRepo.CreateFollower(senderID, receiverID); err != nil {
			return nil, errors.New("could not follow user")
		}

		notification := &models.Notification{
			ID:          uuid.New().String(),
			UserID:      receiverID,
			ActorID:     senderID,
			Type:        "new_follower",
			ReferenceID: senderID,
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		if err := s.notificationRepo.CreateNotification(notification); err != nil {
			return nil, errors.New("could not create notification")
		}

		return notification, nil
	}

	// PRIVATE ACCOUNT
	req := &models.FollowRequest{
		ID:         uuid.New().String(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     models.StatusPending,
		CreatedAt:  time.Now(),
	}

	if err := s.followerRepo.CreateFollowRequest(req); err != nil {
		return nil, errors.New("could not send follow request")
	}

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
		return nil, errors.New("could not create notification")
	}

	return notification, nil
}

func (s *FollowerService) AcceptFollowRequest(requestID, receiverID string) (*models.Notification, error) {
	req, err := s.followerRepo.GetFollowRequestByID(requestID)
	if err != nil {
		return nil, errors.New("follow request not found")
	}

	if req.Status != models.StatusPending {
		return nil, errors.New("follow request is no longer pending")
	}

	if req.ReceiverID != receiverID {
		return nil, errors.New("unauthorized")
	}

	if err := s.followerRepo.DeleteFollowRequest(requestID); err != nil {
		return nil, errors.New("could not remove follow request")
	}

	if err := s.followerRepo.CreateFollower(req.SenderID, req.ReceiverID); err != nil {
		return nil, errors.New("could not create follower relationship")
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
		return nil, errors.New("could not create notification")
	}

	return notification, nil
}

func (s *FollowerService) DeclineFollowRequest(requestID, receiverID string) (*models.Notification, error) {
	req, err := s.followerRepo.GetFollowRequestByID(requestID)
	if err != nil {
		return nil, errors.New("follow request not found")
	}

	if req.Status != models.StatusPending {
		return nil, errors.New("follow request is no longer pending")
	}

	if req.ReceiverID != receiverID {
		return nil, errors.New("unauthorized")
	}

	if err := s.followerRepo.DeleteFollowRequest(requestID); err != nil {
		return nil, errors.New("could not remove follow request")
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
		return nil, errors.New("could not create notification")
	}

	return notification, nil
}

func (s *FollowerService) Unfollow(followerID, followingID string) error {
	already, err := s.followerRepo.IsFollowing(followerID, followingID)
	if err != nil {
		return errors.New("could not check follow status")
	}

	if !already {
		return errors.New("you are not following this user")
	}

	if err := s.followerRepo.DeleteFollower(followerID, followingID); err != nil {
		return errors.New("could not unfollow user")
	}

	return nil
}

func (s *FollowerService) GetFollowers(userID string) ([]*models.FollowerProfile, error) {
	return s.followerRepo.GetFollowers(userID)
}

func (s *FollowerService) GetFollowing(userID string) ([]*models.FollowerProfile, error) {
	return s.followerRepo.GetFollowing(userID)
}
