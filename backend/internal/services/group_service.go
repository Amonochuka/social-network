package services

import (
	"errors"

	"social-network/backend/internal/models"
	"social-network/backend/internal/repositories/interfaces"
)

type GroupService struct {
	groupRepo interfaces.GroupRepository
}

func NewGroupService(groupRepo interfaces.GroupRepository) *GroupService {
	return &GroupService{groupRepo: groupRepo}
}

// ── Core ─────────────────────────────────────────────────────────────────────

func (s *GroupService) CreateGroup(creatorID, title, description string) (*models.Group, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	g := &models.Group{
		CreatorID:   creatorID,
		Title:       title,
		Description: description,
	}
	if err := s.groupRepo.CreateGroup(g); err != nil {
		return nil, errors.New("could not create group")
	}
	// Creator automatically becomes a member
	if err := s.groupRepo.AddMember(g.ID, creatorID); err != nil {
		return nil, errors.New("could not add creator as member")
	}
	return g, nil
}

func (s *GroupService) GetAllGroups(viewerID string) ([]*models.Group, error) {
	return s.groupRepo.GetAllGroups(viewerID)
}

func (s *GroupService) GetGroupByID(groupID string) (*models.Group, error) {
	g, err := s.groupRepo.GetGroupByID(groupID)
	if err != nil {
		return nil, errors.New("group not found")
	}
	return g, nil
}

func (s *GroupService) GetMembers(groupID string) ([]*models.GroupMember, error) {
	return s.groupRepo.GetMembers(groupID)
}

// ── Invitations ──────────────────────────────────────────────────────────────

func (s *GroupService) InviteUser(groupID, inviterID, inviteeID string) error {
	// Inviter must be a member
	isMember, err := s.groupRepo.IsMember(groupID, inviterID)
	if err != nil || !isMember {
		return errors.New("only group members can invite users")
	}
	// Invitee must not already be a member
	alreadyMember, _ := s.groupRepo.IsMember(groupID, inviteeID)
	if alreadyMember {
		return errors.New("user is already a member")
	}
	// No duplicate pending invitation
	existing, _ := s.groupRepo.GetPendingInvitation(groupID, inviteeID)
	if existing != nil {
		return errors.New("invitation already sent")
	}
	inv := &models.GroupInvitation{
		GroupID:   groupID,
		InviterID: inviterID,
		InviteeID: inviteeID,
	}
	return s.groupRepo.CreateInvitation(inv)
}

func (s *GroupService) AcceptInvitation(invID, userID string) error {
	inv, err := s.groupRepo.GetInvitationByID(invID)
	if err != nil {
		return errors.New("invitation not found")
	}
	if inv.InviteeID != userID {
		return errors.New("unauthorized")
	}
	if inv.Status != "pending" {
		return errors.New("invitation is no longer pending")
	}
	if err := s.groupRepo.UpdateInvitation(invID, "accepted"); err != nil {
		return errors.New("could not update invitation")
	}
	return s.groupRepo.AddMember(inv.GroupID, userID)
}

func (s *GroupService) RejectInvitation(invID, userID string) error {
	inv, err := s.groupRepo.GetInvitationByID(invID)
	if err != nil {
		return errors.New("invitation not found")
	}
	if inv.InviteeID != userID {
		return errors.New("unauthorized")
	}
	if inv.Status != "pending" {
		return errors.New("invitation is no longer pending")
	}
	return s.groupRepo.UpdateInvitation(invID, "rejected")
}

// ── Join requests ────────────────────────────────────────────────────────────

func (s *GroupService) RequestToJoin(groupID, userID string) error {
	isMember, _ := s.groupRepo.IsMember(groupID, userID)
	if isMember {
		return errors.New("already a member")
	}
	existing, _ := s.groupRepo.GetPendingJoinRequest(groupID, userID)
	if existing != nil {
		return errors.New("join request already sent")
	}
	req := &models.GroupRequest{GroupID: groupID, UserID: userID}
	return s.groupRepo.CreateJoinRequest(req)
}

func (s *GroupService) AcceptJoinRequest(reqID, adminID string) error {
	req, err := s.groupRepo.GetJoinRequestByID(reqID)
	if err != nil {
		return errors.New("join request not found")
	}
	// Only group creator / existing member with admin rights can accept
	// For now: any existing member can accept (group creator is always a member)
	isMember, _ := s.groupRepo.IsMember(req.GroupID, adminID)
	if !isMember {
		return errors.New("unauthorized")
	}
	if req.Status != "pending" {
		return errors.New("request is no longer pending")
	}
	if err := s.groupRepo.UpdateJoinRequest(reqID, "accepted"); err != nil {
		return err
	}
	return s.groupRepo.AddMember(req.GroupID, req.UserID)
}

func (s *GroupService) RejectJoinRequest(reqID, adminID string) error {
	req, err := s.groupRepo.GetJoinRequestByID(reqID)
	if err != nil {
		return errors.New("join request not found")
	}
	isMember, _ := s.groupRepo.IsMember(req.GroupID, adminID)
	if !isMember {
		return errors.New("unauthorized")
	}
	if req.Status != "pending" {
		return errors.New("request is no longer pending")
	}
	return s.groupRepo.UpdateJoinRequest(reqID, "rejected")
}

// ── Group posts & comments ────────────────────────────────────────────────────

func (s *GroupService) CreateGroupPost(groupID, userID, content, mediaPath, mediaType string) (*models.GroupPost, error) {
	isMember, _ := s.groupRepo.IsMember(groupID, userID)
	if !isMember {
		return nil, errors.New("only group members can post")
	}
	post := &models.GroupPost{
		GroupID:   groupID,
		UserID:    userID,
		Content:   content,
		MediaPath: mediaPath,
		MediaType: mediaType,
	}
	if err := s.groupRepo.CreateGroupPost(post); err != nil {
		return nil, errors.New("could not create post")
	}
	return post, nil
}

func (s *GroupService) GetGroupPosts(groupID, viewerID string) ([]*models.GroupPost, error) {
	isMember, _ := s.groupRepo.IsMember(groupID, viewerID)
	if !isMember {
		return nil, errors.New("only group members can view posts")
	}
	return s.groupRepo.GetGroupPosts(groupID)
}

func (s *GroupService) CreateGroupComment(groupPostID, userID, content, mediaPath, mediaType string) (*models.GroupComment, error) {
	post, err := s.groupRepo.GetGroupPostByID(groupPostID)
	if err != nil {
		return nil, errors.New("post not found")
	}
	isMember, _ := s.groupRepo.IsMember(post.GroupID, userID)
	if !isMember {
		return nil, errors.New("only group members can comment")
	}
	comment := &models.GroupComment{
		GroupPostID: groupPostID,
		UserID:      userID,
		Content:     content,
		MediaPath:   mediaPath,
		MediaType:   mediaType,
	}
	if err := s.groupRepo.CreateGroupComment(comment); err != nil {
		return nil, errors.New("could not create comment")
	}
	return comment, nil
}

func (s *GroupService) GetGroupComments(groupPostID, viewerID string) ([]*models.GroupComment, error) {
	post, err := s.groupRepo.GetGroupPostByID(groupPostID)
	if err != nil {
		return nil, errors.New("post not found")
	}
	isMember, _ := s.groupRepo.IsMember(post.GroupID, viewerID)
	if !isMember {
		return nil, errors.New("only group members can view comments")
	}
	return s.groupRepo.GetGroupComments(groupPostID)
}
