package models

import "github.com/google/uuid"

type Bicycle struct {
	BaseModel
	Color             string `json:"color" gorm:"not null;size:50" validate:"required,max=50"`
	Model             string `json:"model" gorm:"size:100" validate:"max=100"`
	IdentificationTag string `json:"identification_tag" gorm:"size:50" validate:"max=50"`

	// Chave estrangeira
	ApartmentID uuid.UUID `json:"apartment_id" gorm:"not null" validate:"required"`
	Apartment   Apartment `json:"apartment,omitempty" gorm:"foreignKey:ApartmentID"`
}
