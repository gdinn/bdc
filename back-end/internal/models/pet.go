package models

import "github.com/google/uuid"

type PetSpecies string

const (
	PetSpeciesDog    PetSpecies = "DOG"
	PetSpeciesCat    PetSpecies = "CAT"
	PetSpeciesBird   PetSpecies = "BIRD"
	PetSpeciesRabbit PetSpecies = "RABBIT"
	PetSpeciesOther  PetSpecies = "OTHER"
)

type Pet struct {
	BaseModel
	Name    string     `json:"name" gorm:"not null;size:100" validate:"required,max=100"`
	Species PetSpecies `json:"species" gorm:"not null;-:migration" validate:"required"`
	Breed   string     `json:"breed" gorm:"size:100" validate:"max=100"`
	Size    string     `json:"size" gorm:"size:20" validate:"max=20"`

	// Chave estrangeira
	ApartmentID uuid.UUID `json:"apartment_id" gorm:"not null type:uuid" validate:"required"`
	Apartment   Apartment `json:"apartment,omitempty" gorm:"foreignKey:ApartmentID"`
}
