package models

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
	Species PetSpecies `json:"species" gorm:"not null" validate:"required"`
	Breed   string     `json:"breed" gorm:"size:100" validate:"max=100"`
	Size    string     `json:"size" gorm:"size:20" validate:"max=20"`

	// Chave estrangeira
	ApartmentID uint      `json:"apartment_id" gorm:"not null" validate:"required"`
	Apartment   Apartment `json:"apartment,omitempty" gorm:"foreignKey:ApartmentID"`
}
