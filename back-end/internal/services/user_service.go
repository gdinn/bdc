package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"bdc/internal/models"
	"bdc/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
	// Validate business rules
	if err := s.validateUserCreation(user); err != nil {
		return nil, err
	}

	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking email existence: %w", err)
	}

	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Create user
	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return createdUser, nil
}

func (s *UserService) validateUserCreation(user *models.User) error {
	// Business rules validation
	if user.IsManager && user.IsAdvisor {
		return errors.New("user cannot be both manager and advisor")
	}

	// Children cannot be managers or advisors
	if user.AgeGroup == models.UserAgeGroupChild && (user.IsManager || user.IsAdvisor) {
		return errors.New("children cannot be managers or advisors")
	}

	// External users typically shouldn't be managers (business rule)
	if user.Type == models.UserTypeExternal && user.IsManager {
		return errors.New("external users cannot be managers")
	}

	return nil
}
