package service

import (
	"time"

	"hybrid-app/backend/internal/domain"
	"hybrid-app/backend/internal/repository"
)

type ProfileService struct {
	repo repository.Repository
}

func NewProfileService(repo repository.Repository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) Me(userID string) (*domain.User, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *ProfileService) UpdateMe(userID string, input domain.ProfileInput, hideFromCities []string) (*domain.User, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Bio != "" {
		user.Bio = input.Bio
	}
	if input.Gender != "" {
		user.Gender = input.Gender
	}
	if input.InterestedIn != nil {
		user.InterestedIn = input.InterestedIn
	}
	if input.City != "" {
		user.City = input.City
	}
	if input.DateOfBirth != "" {
		user.DateOfBirth = input.DateOfBirth
	}
	if input.Photos != nil {
		user.Photos = input.Photos
	}
	if input.Interests != nil {
		user.Interests = input.Interests
	}
	user.InstagramHandle = input.InstagramHandle
	user.HideFromCities = hideFromCities
	user.UpdatedAt = time.Now().UTC()
	return user, s.repo.SaveUser(user)
}
