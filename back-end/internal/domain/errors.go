package domain

import "errors"

var (
	// Validation errors
	ErrUnauthorizedEmailCreation = errors.New("cannot create user with different email from initial registration")
	ErrEmailAlreadyExists        = errors.New("email already exists")

	// Input errors
	ErrEmailEmpty                        = errors.New("email cannot be empty")
	ErrExternalUsersCannotBeChild        = errors.New("external users cannot be child")
	ErrChildrenCannotBeManagerOrAdvisors = errors.New("children cannot be managers or advisors")
	ErrAgeGroupInvalid                   = errors.New("ageGroup is invalid")
	ErrTypeInvalid                       = errors.New("type is invalid")
)
