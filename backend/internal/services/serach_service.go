package services

import (
	"strings"

	"social-network/backend/internal/models"
	"social-network/backend/internal/repositories/interfaces"
)

type SearchService struct {
	userRepo  interfaces.UserRepository
	groupRepo interfaces.GroupRepository
}

func NewSearchService(
	userRepo interfaces.UserRepository,
	groupRepo interfaces.GroupRepository,
) *SearchService {
	return &SearchService{
		userRepo:  userRepo,
		groupRepo: groupRepo,
	}
}

func (s *SearchService) Search(query, currentUserID string) (*models.SearchResponse, error) {
	query = strings.TrimSpace(query)

	if query == "" {
		return &models.SearchResponse{
			Users:  []*models.UserSearchResult{},
			Groups: []*models.GroupSearchResult{},
		}, nil
	}

	users, err := s.userRepo.SearchUsers(query, currentUserID)
	if err != nil {
		return nil, err
	}

	groups, err := s.groupRepo.SearchGroups(query)
	if err != nil {
		return nil, err
	}

	return &models.SearchResponse{
		Users:  users,
		Groups: groups,
	}, nil
}
