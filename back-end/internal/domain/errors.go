package domain

import "errors"

var (
	ErrUnauthorizedEmailCreation         = errors.New("cannot create user with different email from initial registration")
	ErrEmailAlreadyExists                = errors.New("email already exists")
	ErrExternalUsersCannotBeChild        = errors.New("external users cannot be child")
	ErrChildrenCannotBeManagerOrAdvisors = errors.New("children cannot be managers or advisors") // example
	ErrAgeGroupInvalid                   = errors.New("ageGroup is invalid")
	ErrTypeInvalid                       = errors.New("type is invalid")
)
