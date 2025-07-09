package models

import (
	"time"
)

type UserPermission string

const (
	UserPermissionRead      UserPermission = "READ"
	UserPermissionWrite     UserPermission = "WRITE"
	UserPermissionReadWrite UserPermission = "READ_WRITE"
	UserPermissionNone      UserPermission = "NO_PERMISSION"
)

type UserType string

const (
	UserTypeResident UserType = "RESIDENT"
	UserTypeExternal UserType = "EXTERNAL"
)

type UserAgeGroup string

const (
	UserAgeGroupAdult UserAgeGroup = "ADULT"
	UserAgeGroupChild UserAgeGroup = "CHILD"
)

type User struct {
	BaseModel
	Name      string       `json:"name" gorm:"not null;size:255" validate:"required,min=2,max=255"`
	Email     string       `json:"email" gorm:"uniqueIndex;not null;size:255" validate:"required,email"`
	Phone     string       `json:"phone" gorm:"size:20" validate:"omitempty,min=10,max=20"`
	BirthDate *time.Time   `json:"birth_date,omitempty"`
	Type      UserType     `json:"type" gorm:"not null;default:'EXTERNAL'" validate:"required"`
	AgeGroup  UserAgeGroup `json:"age_group" gorm:"not null;default:'ADULT'" validate:"required"`
	IsManager bool         `json:"is_manager" gorm:"default:false"`
	IsAdvisor bool         `json:"is_advisor" gorm:"default:false"`

	// Relacionamentos
	Apartments         []Apartment `json:"apartments,omitempty" gorm:"many2many:user_apartments;"`
	LegalRepApartments []Apartment `json:"legal_rep_apartments,omitempty" gorm:"foreignKey:LegalRepresentativeID"`
}

// CheckPermissions retorna as permissões do usuário
func (u *User) CheckPermissions() UserPermission {
	if u.IsManager {
		return UserPermissionReadWrite
	}
	if u.IsAdvisor {
		return UserPermissionRead
	}
	return UserPermissionReadWrite // Para residentes em suas próprias unidades
}

// ListAssociatedApartments lista os apartamentos associados ao usuário
func (u *User) ListAssociatedApartments() []Apartment {
	// Esta implementação seria feita na camada de serviço
	// aqui apenas retornamos os apartamentos já carregados
	return u.Apartments
}
