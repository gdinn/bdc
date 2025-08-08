package models

import "github.com/google/uuid"

type Apartment struct {
	BaseModel
	Number   string `json:"number" gorm:"not null;size:10" validate:"required,max=10"`
	Building string `json:"building" gorm:"not null;size:10" validate:"required,max=10"`

	// Chave estrangeira para representante legal
	LegalRepresentativeID *uuid.UUID `json:"legal_representative_id,omitempty" gorm:"type:uuid"`
	LegalRepresentative   *User      `json:"legal_representative,omitempty" gorm:"foreignKey:LegalRepresentativeID"`

	// Relacionamentos
	Users    []User    `json:"users,omitempty" gorm:"many2many:user_apartments;"`
	Vehicles []Vehicle `json:"vehicles,omitempty" gorm:"foreignKey:ApartmentID"`
	Pets     []Pet     `json:"pets,omitempty" gorm:"foreignKey:ApartmentID"`
	Bicycles []Bicycle `json:"bicycles,omitempty" gorm:"foreignKey:ApartmentID"`
}

// AddUser adiciona um usuário ao apartamento
func (a *Apartment) AddUser(user *User) {
	a.Users = append(a.Users, *user)
}

// RemoveUser remove um usuário do apartamento
func (a *Apartment) RemoveUser(userID uuid.UUID) {
	for i, user := range a.Users {
		if user.ID == userID {
			a.Users = append(a.Users[:i], a.Users[i+1:]...)
			break
		}
	}
}

// AssignLegalRep atribui um representante legal ao apartamento
func (a *Apartment) AssignLegalRep(user *User) {
	a.LegalRepresentativeID = &user.ID
	a.LegalRepresentative = user
}
