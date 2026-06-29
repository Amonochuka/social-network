package services

import (
	"errors"
	"time"

	"social-network/backend/internal/models"
	"social-network/backend/internal/repositories/interfaces"
)

type EventService struct {
	eventRepo interfaces.EventRepository
	groupRepo interfaces.GroupRepository
}

func NewEventService(eventRepo interfaces.EventRepository, groupRepo interfaces.GroupRepository) *EventService {
	return &EventService{eventRepo: eventRepo, groupRepo: groupRepo}
}

func (s *EventService) CreateEvent(groupID, creatorID, title, description string, eventTime time.Time) (*models.Event, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	isMember, _ := s.groupRepo.IsMember(groupID, creatorID)
	if !isMember {
		return nil, errors.New("only group members can create events")
	}
	e := &models.Event{
		GroupID:     groupID,
		CreatorID:   creatorID,
		Title:       title,
		Description: description,
		EventTime:   eventTime,
	}
	if err := s.eventRepo.CreateEvent(e); err != nil {
		return nil, errors.New("could not create event")
	}
	return e, nil
}

func (s *EventService) GetEventsByGroup(groupID, viewerID string) ([]*models.Event, error) {
	isMember, _ := s.groupRepo.IsMember(groupID, viewerID)
	if !isMember {
		return nil, errors.New("only group members can view events")
	}
	return s.eventRepo.GetEventsByGroupID(groupID, viewerID)
}

func (s *EventService) RespondToEvent(eventID, userID, response string) error {
	if response != "going" && response != "not_going" {
		return errors.New("response must be 'going' or 'not_going'")
	}
	e, err := s.eventRepo.GetEventByID(eventID, userID)
	if err != nil {
		return errors.New("event not found")
	}
	isMember, _ := s.groupRepo.IsMember(e.GroupID, userID)
	if !isMember {
		return errors.New("only group members can respond to events")
	}
	return s.eventRepo.UpsertResponse(&models.EventResponse{
		EventID:  eventID,
		UserID:   userID,
		Response: response,
	})
}
