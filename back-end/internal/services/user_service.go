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
	if user.Type == models.UserTypeExternal && user.AgeGroup == models.UserAgeGroupChild {
		return errors.New("external users cannot be child") // example
	}

	if user.AgeGroup == models.UserAgeGroupChild && ((user.Role == models.UserRoleAdvisor) || (user.Role == models.UserRoleManager)) {
		return errors.New("children cannot be managers or advisors") // example
	}

	if !isValidUserAgeGroup(user.AgeGroup) {
		return errors.New("ageGroup is invalid")
	}

	if !isValidUserType(user.Type) {
		return errors.New("type is invalid")
	}

	return nil
}

func isValidUserType(userType models.UserType) bool {
	return userType == models.UserTypeResident || userType == models.UserTypeExternal
}

func isValidUserAgeGroup(userAgeGroup models.UserAgeGroup) bool {
	return userAgeGroup == models.UserAgeGroupAdult || userAgeGroup == models.UserAgeGroupChild
}
