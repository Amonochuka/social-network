package services

import (
	"errors"
	"social-network/backend/internal/models"
	"social-network/backend/internal/repositories/interfaces"
	"time"

	"github.com/google/uuid"
)

type GroupService struct {
	groupRepo        interfaces.GroupRepository
	userRepo         interfaces.UserRepository
	followerRepo     interfaces.FollowerRepository
	notificationRepo interfaces.NotificationRepository
}

func NewGroupService(
	groupRepo interfaces.GroupRepository,
	userRepo interfaces.UserRepository,
	followerRepo interfaces.FollowerRepository,
	notificationRepo interfaces.NotificationRepository,
) *GroupService {
	return &GroupService{
		groupRepo:        groupRepo,
		userRepo:         userRepo,
		followerRepo:     followerRepo,
		notificationRepo: notificationRepo,
	}
}

// ── Groups ───────────────────────────────────────────────────────────────────

func (s *GroupService) CreateGroup(creatorID, title, description string) (*models.Group, error) {
	if title == "" || description == "" {
		return nil, errors.New("title and description are required")
	}
	group := &models.Group{
		ID:          uuid.New().String(),
		CreatorID:   creatorID,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
	}
	if err := s.groupRepo.CreateGroup(group); err != nil {
		return nil, errors.New("could not create group")
	}
	// creator is automatically a member
	if err := s.groupRepo.AddMember(group.ID, creatorID); err != nil {
		return nil, errors.New("could not add creator as member")
	}
	return group, nil
}

func (s *GroupService) GetGroup(groupID, viewerID string) (*models.GroupDetail, error) {
	group, err := s.groupRepo.GetGroupByID(groupID)
	if err != nil {
		return nil, errors.New("group not found")
	}
	members, _ := s.groupRepo.GetMembers(groupID)
	isMember, _ := s.groupRepo.IsMember(groupID, viewerID)
	invStatus, _ := s.groupRepo.GetInvitationStatus(groupID, viewerID)
	reqStatus, _ := s.groupRepo.GetJoinRequestStatus(groupID, viewerID)

	return &models.GroupDetail{
		Group:         *group,
		Members:       members,
		IsMember:      isMember,
		IsCreator:     group.CreatorID == viewerID,
		InviteStatus:  invStatus,
		RequestStatus: reqStatus,
	}, nil
}

func (s *GroupService) GetAllGroups() ([]*models.Group, error) {
	groups, err := s.groupRepo.GetAllGroups()
	if err != nil {
		return nil, errors.New("could not list groups")
	}
	return groups, nil
}

// ── Invitations ──────────────────────────────────────────────────────────────

func (s *GroupService) InviteUser(groupID, inviterID, inviteeID string) error {
	// inviter must be a member
	isMember, err := s.groupRepo.IsMember(groupID, inviterID)
	if err != nil || !isMember {
		return errors.New("you must be a group member to invite")
	}
	// invitee must not already be a member
	alreadyMember, _ := s.groupRepo.IsMember(groupID, inviteeID)
	if alreadyMember {
		return errors.New("user is already a member")
	}
	// no pending invite already
	if existing, _ := s.groupRepo.GetPendingInvitation(groupID, inviteeID); existing != nil {
		return errors.New("invitation already sent")
	}

	inv := &models.GroupInvitation{
		ID:        uuid.New().String(),
		GroupID:   groupID,
		InviterID: inviterID,
		InviteeID: inviteeID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	if err := s.groupRepo.CreateInvitation(inv); err != nil {
		return errors.New("could not create invitation")
	}

	// notify invitee
	group, _ := s.groupRepo.GetGroupByID(groupID)
	groupTitle := ""
	if group != nil {
		groupTitle = group.ID
	}
	_ = groupTitle
	s.notificationRepo.CreateNotification(&models.Notification{
		ID:          uuid.New().String(),
		UserID:      inviteeID,
		ActorID:     inviterID,
		Type:        "group_invitation",
		ReferenceID: inv.ID,
		IsRead:      false,
		CreatedAt:   time.Now(),
	})
	return nil
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
	if err := s.groupRepo.UpdateInvitationStatus(invID, "accepted"); err != nil {
		return errors.New("could not update invitation")
	}
	return s.groupRepo.AddMember(inv.GroupID, userID)
}

func (s *GroupService) DeclineInvitation(invID, userID string) error {
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
	return s.groupRepo.UpdateInvitationStatus(invID, "declined")
}

// ── Join Requests ─────────────────────────────────────────────────────────────

func (s *GroupService) RequestToJoin(groupID, userID string) error {
	isMember, _ := s.groupRepo.IsMember(groupID, userID)
	if isMember {
		return errors.New("you are already a member")
	}
	if existing, _ := s.groupRepo.GetPendingJoinRequest(groupID, userID); existing != nil {
		return errors.New("join request already sent")
	}
	group, err := s.groupRepo.GetGroupByID(groupID)
	if err != nil {
		return errors.New("group not found")
	}

	req := &models.GroupRequest{
		ID:        uuid.New().String(),
		GroupID:   groupID,
		UserID:    userID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	if err := s.groupRepo.CreateJoinRequest(req); err != nil {
		return errors.New("could not send join request")
	}

	// notify group creator
	s.notificationRepo.CreateNotification(&models.Notification{
		ID:          uuid.New().String(),
		UserID:      group.CreatorID,
		ActorID:     userID,
		Type:        "group_join_request",
		ReferenceID: req.ID,
		IsRead:      false,
		CreatedAt:   time.Now(),
	})
	return nil
}

func (s *GroupService) AcceptJoinRequest(reqID, creatorID string) error {
	req, err := s.groupRepo.GetJoinRequestByID(reqID)
	if err != nil {
		return errors.New("request not found")
	}
	group, err := s.groupRepo.GetGroupByID(req.GroupID)
	if err != nil || group.CreatorID != creatorID {
		return errors.New("only the group creator can accept requests")
	}
	if req.Status != "pending" {
		return errors.New("request is no longer pending")
	}
	if err := s.groupRepo.UpdateJoinRequestStatus(reqID, "accepted"); err != nil {
		return errors.New("could not update request")
	}
	return s.groupRepo.AddMember(req.GroupID, req.UserID)
}

func (s *GroupService) DeclineJoinRequest(reqID, creatorID string) error {
	req, err := s.groupRepo.GetJoinRequestByID(reqID)
	if err != nil {
		return errors.New("request not found")
	}
	group, err := s.groupRepo.GetGroupByID(req.GroupID)
	if err != nil || group.CreatorID != creatorID {
		return errors.New("only the group creator can decline requests")
	}
	if req.Status != "pending" {
		return errors.New("request is no longer pending")
	}
	return s.groupRepo.UpdateJoinRequestStatus(reqID, "declined")
}

// ── Events ────────────────────────────────────────────────────────────────────

func (s *GroupService) CreateEvent(groupID, creatorID, title, description, eventTimeStr string) (*models.Event, error) {
	isMember, _ := s.groupRepo.IsMember(groupID, creatorID)
	if !isMember {
		return nil, errors.New("you must be a group member to create an event")
	}
	eventTime, err := time.Parse(time.RFC3339, eventTimeStr)
	if err != nil {
		// try date-only format
		eventTime, err = time.Parse("2006-01-02T15:04", eventTimeStr)
		if err != nil {
			return nil, errors.New("invalid event_time format, use RFC3339 or YYYY-MM-DDTHH:MM")
		}
	}
	event := &models.Event{
		ID:          uuid.New().String(),
		GroupID:     groupID,
		CreatorID:   creatorID,
		Title:       title,
		Description: description,
		EventTime:   eventTime,
		CreatedAt:   time.Now(),
	}
	if err := s.groupRepo.CreateEvent(event); err != nil {
		return nil, errors.New("could not create event")
	}

	// notify all group members (except creator)
	memberIDs, _ := s.groupRepo.GetMemberIDs(groupID)
	for _, memberID := range memberIDs {
		if memberID == creatorID {
			continue
		}
		s.notificationRepo.CreateNotification(&models.Notification{
			ID:          uuid.New().String(),
			UserID:      memberID,
			ActorID:     creatorID,
			Type:        "event_created",
			ReferenceID: event.ID,
			IsRead:      false,
			CreatedAt:   time.Now(),
		})
	}
	return event, nil
}

func (s *GroupService) GetEvents(groupID, userID string) ([]*models.EventDetail, error) {
	isMember, _ := s.groupRepo.IsMember(groupID, userID)
	if !isMember {
		return nil, errors.New("you must be a group member to view events")
	}
	events, err := s.groupRepo.GetEventsByGroupID(groupID)
	if err != nil {
		return nil, errors.New("could not load events")
	}
	var details []*models.EventDetail
	for _, e := range events {
		going, notGoing, _ := s.groupRepo.GetEventResponses(e.ID)
		myResp, _ := s.groupRepo.GetMyEventResponse(e.ID, userID)
		creator, _ := s.userRepo.GetUserByID(e.CreatorID)
		creatorName := ""
		if creator != nil {
			creatorName = creator.FirstName + " " + creator.LastName
		}
		details = append(details, &models.EventDetail{
			ID:            e.ID,
			GroupID:       e.GroupID,
			CreatorID:     e.CreatorID,
			CreatorName:   creatorName,
			Title:         e.Title,
			Description:   e.Description,
			EventTime:     e.EventTime,
			CreatedAt:     e.CreatedAt,
			GoingCount:    going,
			NotGoingCount: notGoing,
			MyResponse:    myResp,
		})
	}
	return details, nil
}

func (s *GroupService) RSVPEvent(eventID, userID, response string) error {
	if response != "going" && response != "not_going" {
		return errors.New("response must be 'going' or 'not_going'")
	}
	event, err := s.groupRepo.GetEventByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}
	isMember, _ := s.groupRepo.IsMember(event.GroupID, userID)
	if !isMember {
		return errors.New("you must be a group member to respond to events")
	}
	return s.groupRepo.UpsertEventResponse(eventID, userID, response)
}

// ── Group Posts ───────────────────────────────────────────────────────────────

func (s *GroupService) CreateGroupPost(groupID, userID, content, mediaPath, mediaType string) (*models.GroupPost, error) {
	isMember, _ := s.groupRepo.IsMember(groupID, userID)
	if !isMember {
		return nil, errors.New("you must be a group member to post")
	}
	post := &models.GroupPost{
		ID:        uuid.New().String(),
		GroupID:   groupID,
		UserID:    userID,
		Content:   content,
		MediaPath: mediaPath,
		MediaType: mediaType,
		CreatedAt: time.Now(),
	}
	if err := s.groupRepo.CreateGroupPost(post); err != nil {
		return nil, errors.New("could not create group post")
	}
	return post, nil
}

func (s *GroupService) GetGroupPosts(groupID, userID string) ([]*models.GroupPostDetail, error) {
	isMember, _ := s.groupRepo.IsMember(groupID, userID)
	if !isMember {
		return nil, errors.New("you must be a group member to view posts")
	}
	return s.groupRepo.GetGroupPosts(groupID)
}

func (s *GroupService) CreateGroupComment(groupPostID, userID, content, mediaPath, mediaType string) (*models.GroupComment, error) {
	comment := &models.GroupComment{
		ID:          uuid.New().String(),
		GroupPostID: groupPostID,
		UserID:      userID,
		Content:     content,
		MediaPath:   mediaPath,
		MediaType:   mediaType,
		CreatedAt:   time.Now(),
	}
	if err := s.groupRepo.CreateGroupComment(comment); err != nil {
		return nil, errors.New("could not create comment")
	}
	return comment, nil
}

func (s *GroupService) GetGroupComments(groupPostID string) ([]*models.GroupCommentDetail, error) {
	return s.groupRepo.GetGroupComments(groupPostID)
}

// ── Group Chat ────────────────────────────────────────────────────────────────

func (s *GroupService) SendGroupMessage(groupID, senderID, content string) (*models.GroupMessage, error) {
	isMember, _ := s.groupRepo.IsMember(groupID, senderID)
	if !isMember {
		return nil, errors.New("you must be a group member to send messages")
	}
	msg := &models.GroupMessage{
		ID:        uuid.New().String(),
		GroupID:   groupID,
		SenderID:  senderID,
		Content:   content,
		CreatedAt: time.Now(),
	}
	if err := s.groupRepo.CreateGroupChatMessage(msg); err != nil {
		return nil, errors.New("could not send message")
	}
	return msg, nil
}

func (s *GroupService) GetGroupMessages(groupID, userID string) ([]*models.GroupMessageDetailFull, error) {
	isMember, _ := s.groupRepo.IsMember(groupID, userID)
	if !isMember {
		return nil, errors.New("you must be a group member to view messages")
	}
	return s.groupRepo.GetGroupChatMessages(groupID)
}

func (s *GroupService) GetMemberIDs(groupID string) ([]string, error) {
	return s.groupRepo.GetMemberIDs(groupID)
}
