package models

import (
	"time"
)

type UserPermission string

const (
	READ          UserPermission = "READ"
	WRITE         UserPermission = "WRITE"
	READ_WRITE    UserPermission = "READ_WRITE"
	NO_PERMISSION UserPermission = "NO_PERMISSION"
)

type UserType string

const (
	RESIDENT UserType = "RESIDENT"
	EXTERNAL UserType = "EXTERNAL"
)

type UserAgeGroup string

const (
	ADULT UserAgeGroup = "ADULT"
	CHILD UserAgeGroup = "CHILD"
)

type User struct {
	ID         int          `json:"id,omitempty"`
	Name       string       `json:"name"`
	Email      string       `json:"email"`
	Phone      string       `json:"phone,omitempty"`
	BirthDate  time.Time    `json:"birthDate"`
	Type       UserType     `json:"type"`
	AgeGroup   UserAgeGroup `json:"ageGroup"`
	IsManager  bool         `json:"isManager"`
	IsAdvisor  bool         `json:"isAdvisor"`
	IsLegalRep bool         `json:"isLeglRep"`
}
