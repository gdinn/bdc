package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"bdc/internal/domain"
	"bdc/internal/models"
	"bdc/internal/repositories"
	"bdc/internal/utils"
)

type UserService struct {
	userRepo       *repositories.UserRepository
	cognitoService *CognitoService
}

func NewUserService(userRepo *repositories.UserRepository, cognitoService *CognitoService) *UserService {
	return &UserService{
		userRepo:       userRepo,
		cognitoService: cognitoService,
	}
}

func (s *UserService) CreateUserWithContext(ctx *domain.CreateUserContext) (*models.User, error) {
	if err := s.validateBusinessRules(ctx.User); err != nil {
		return nil, err
	}

	cognitoUser, err := s.cognitoService.GetUserInCognito(ctx.Claims.Username)
	if err != nil {
		return nil, err
	}
	ctx.User.Name, err = utils.GetUserAttribute(cognitoUser, "name")
	if err != nil {
		return nil, err
	}
	ctx.User.Email, err = utils.GetUserAttribute(cognitoUser, "email")
	if err != nil {
		return nil, err
	}

	if err := s.validateUserIsNew(ctx.User); err != nil {
		return nil, err
	}

	createdUser, err := s.userRepo.Create(ctx.User)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return createdUser, nil
}

func (s *UserService) validateUserIsNew(user *models.User) error {
	emailExists, err := s.userRepo.IsEmailExists(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error checking email existence: %w", err)
	}

	if emailExists {
		return domain.ErrEmailAlreadyExists
	}

	return nil
}

func (s *UserService) validateBusinessRules(user *models.User) error {
	if user.Type == models.UserTypeExternal && user.AgeGroup == models.UserAgeGroupChild {
		return domain.ErrExternalUsersCannotBeChild // example
	}

	if user.AgeGroup == models.UserAgeGroupChild && ((user.Role == models.UserRoleAdvisor) || (user.Role == models.UserRoleManager)) {
		return domain.ErrChildrenCannotBeManagerOrAdvisors
	}

	if !isValidUserAgeGroup(user.AgeGroup) {
		return domain.ErrAgeGroupInvalid
	}

	if !isValidUserType(user.Type) {
		return domain.ErrTypeInvalid
	}

	return nil
}

func isValidUserType(userType models.UserType) bool {
	return userType == models.UserTypeResident || userType == models.UserTypeExternal
}

func isValidUserAgeGroup(userAgeGroup models.UserAgeGroup) bool {
	return (userAgeGroup == models.UserAgeGroupAdult) || (userAgeGroup == models.UserAgeGroupChild)
}
