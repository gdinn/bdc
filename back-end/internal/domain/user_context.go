package domain

import (
	"bdc/internal/middleware"
	"bdc/internal/models"
	"time"
)

type CreateUserRequest struct {
	Phone     string              `json:"phone" validate:"omitempty,min=10,max=20"`
	BirthDate *time.Time          `json:"birth_date,omitempty"`
	Type      models.UserType     `json:"type" validate:"required"`
	AgeGroup  models.UserAgeGroup `json:"age_group" validate:"required"`
}

type CreateUserContext struct {
	Claims *middleware.UserClaims
	User   *models.User
}

func CreateUserData(req *CreateUserRequest) *models.User {
	return &models.User{
		Phone:     req.Phone,
		BirthDate: req.BirthDate,
		Type:      req.Type,
		AgeGroup:  req.AgeGroup,
		Role:      models.UserRoleCommon, // Default role
	}
}

func NewCreateUserContext(claims *middleware.UserClaims, user *models.User) *CreateUserContext {
	return &CreateUserContext{
		Claims: claims,
		User:   user,
	}
}
