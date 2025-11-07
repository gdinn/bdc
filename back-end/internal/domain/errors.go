package domain

import "errors"

var (
	// Validation errors
	ErrUnauthorizedEmailCreation = errors.New("cannot create user with different email from initial registration")
	ErrEmailAlreadyExists        = errors.New("email already exists")
	ErrApartmentAlreadyExists    = errors.New("apartment already exists")
	ErrBuildingAlreadyExists     = errors.New("building already exists")
	ErrApartmentViewUnauth       = errors.New("apartment view unauthorized")

	// Input errors
	ErrEmailEmpty                        = errors.New("email cannot be empty")
	ErrExternalUsersCannotBeChild        = errors.New("external users cannot be child")
	ErrChildrenCannotBeManagerOrAdvisors = errors.New("children cannot be managers or advisors")
	ErrAgeGroupInvalid                   = errors.New("ageGroup is invalid")
	ErrTypeInvalid                       = errors.New("type is invalid")
	ErrManagerRoleRequired               = errors.New("manager role is required")
	ErrAdvisorRoleRequired               = errors.New("advisor role is required")
	ErrManagerOrAdvisorRoleRequired      = errors.New("special role is required")
)
