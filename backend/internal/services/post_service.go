package services

import (
	"errors"
	"time"

	"social-network/backend/internal/models"
	"social-network/backend/internal/repositories/interfaces"

	"github.com/google/uuid"
)

type PostService struct {
	postRepo         interfaces.PostRepository
	followerRepo     interfaces.FollowerRepository
	notificationRepo interfaces.NotificationRepository
	userRepo         interfaces.UserRepository
}

func NewPostService(
	postRepo interfaces.PostRepository,
	followerRepo interfaces.FollowerRepository,
	notificationRepo interfaces.NotificationRepository,
	userRepo interfaces.UserRepository,
) *PostService {
	return &PostService{
		postRepo:         postRepo,
		followerRepo:     followerRepo,
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

func (s *PostService) CreatePost(userID, content, mediaPath, mediaType, privacy string, allowedUsers []string) (*models.Post, error) {
	if content == "" && mediaPath == "" {
		return nil, errors.New("post must have content or media")
	}
	if privacy != "public" && privacy != "followers" && privacy != "selected" {
		return nil, errors.New("invalid privacy setting")
	}
	if privacy == "selected" && len(allowedUsers) == 0 {
		return nil, errors.New("selected privacy requires at least one allowed user")
	}
	post := &models.Post{
		ID:        uuid.New().String(),
		UserID:    userID,
		Content:   content,
		MediaPath: mediaPath,
		MediaType: mediaType,
		Privacy:   privacy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.postRepo.CreatePost(post); err != nil {
		return nil, errors.New("could not create post")
	}
	if privacy == "selected" {
		for _, allowedUserID := range allowedUsers {
			if err := s.postRepo.CreateAllowedUser(post.ID, allowedUserID); err != nil {
				return nil, errors.New("could not add allowed user")
			}
		}
	}
	return post, nil
}

func (s *PostService) GetPostsByUserID(userID, viewerID string) ([]*models.Post, error) {
	owner, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Own profile
	if userID == viewerID {
		return s.postRepo.GetPostsByUserID(userID, viewerID)
	}

	// Public profile
	if owner.IsPublic {
		return s.postRepo.GetPostsByUserID(userID, viewerID)
	}

	// Private profile
	following, err := s.followerRepo.IsFollowing(viewerID, userID)
	if err != nil {
		return nil, errors.New("could not check follow status")
	}

	// Not following -> profile exists but no posts
	if !following {
		return []*models.Post{}, nil
	}

	return s.postRepo.GetPostsByUserID(userID, viewerID)
}

func (s *PostService) GetFeed(userID string) ([]*models.FeedPost, error) {
	posts, err := s.postRepo.GetFeed(userID)
	if err != nil {
		return nil, errors.New("could not get feed")
	}
	return posts, nil
}

func (s *PostService) GetPostByID(postID, viewerID string) (*models.Post, error) {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}
	if post.Privacy == "public" {
		return post, nil
	}
	if post.UserID == viewerID {
		return post, nil
	}
	if post.Privacy == "followers" {
		following, err := s.followerRepo.IsFollowing(viewerID, post.UserID)
		if err != nil {
			return nil, errors.New("could not check follow status")
		}
		if !following {
			return nil, errors.New("unauthorized to view this post")
		}
		return post, nil
	}
	if post.Privacy == "selected" {
		allowed, err := s.postRepo.IsAllowedToViewPost(postID, viewerID)
		if err != nil {
			return nil, errors.New("could not check permissions")
		}
		if !allowed {
			return nil, errors.New("unauthorized to view this post")
		}
		return post, nil
	}
	return nil, errors.New("unauthorized to view this post")
}

func (s *PostService) UpdatePost(postID, userID, content, privacy string) (*models.Post, error) {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}
	if post.UserID != userID {
		return nil, errors.New("unauthorized to update this post")
	}
	if content != "" {
		post.Content = content
	}
	if privacy != "" {
		if privacy != "public" && privacy != "followers" && privacy != "selected" {
			return nil, errors.New("invalid privacy setting")
		}
		post.Privacy = privacy
	}
	post.UpdatedAt = time.Now()
	if err := s.postRepo.UpdatePost(post); err != nil {
		return nil, errors.New("could not update post")
	}
	return post, nil
}

func (s *PostService) DeletePost(postID, userID string) error {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return errors.New("post not found")
	}
	if post.UserID != userID {
		return errors.New("unauthorized to delete this post")
	}
	if err := s.postRepo.DeletePost(postID); err != nil {
		return errors.New("could not delete post")
	}
	return nil
}

func (s *PostService) CreateComment(postID, userID, content, mediaPath, mediaType string) (*models.Comment, error) {
	if content == "" && mediaPath == "" {
		return nil, errors.New("comment must have content or media")
	}
	_, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}
	comment := &models.Comment{
		ID:        uuid.New().String(),
		PostID:    postID,
		UserID:    userID,
		Content:   content,
		MediaPath: mediaPath,
		MediaType: mediaType,
		CreatedAt: time.Now(),
	}
	if err := s.postRepo.CreateComment(comment); err != nil {
		return nil, errors.New("could not create comment")
	}
	return comment, nil
}

func (s *PostService) GetCommentsByPostID(postID string) ([]*models.CommentDetail, error) {
	comments, err := s.postRepo.GetCommentsByPostID(postID)
	if err != nil {
		return nil, errors.New("could not get comments")
	}
	return comments, nil
}

func (s *PostService) ToggleLike(postID, userID string) (liked bool, likeCount int, err error) {
	// Ensure the post exists before liking
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return false, 0, errors.New("post not found")
	}

	liked, likeCount, err = s.postRepo.ToggleLike(postID, userID)
	if err != nil {
		return false, 0, err
	}

	// Don't notify yourself
	if post.UserID == userID {
		return liked, likeCount, nil
	}

	if liked {
		// Send "post_liked" notification to the post owner
		notification := &models.Notification{
			ID:          uuid.New().String(),
			UserID:      post.UserID,
			ActorID:     userID,
			Type:        "post_liked",
			ReferenceID: postID,
			IsRead:      false,
			CreatedAt:   time.Now(),
		}
		// Non-fatal — like still succeeds even if notification fails
		_ = s.notificationRepo.CreateNotification(notification)
	} else {
		// Remove the notification when unliking
		_ = s.notificationRepo.DeleteNotificationByReferenceID(postID, "post_liked")
	}

	return liked, likeCount, nil
}
