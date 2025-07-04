package models

type PetSpecies string

const (
	CACHORRO PetSpecies = "DOG"
	GATO     PetSpecies = "CAT"
	AVE      PetSpecies = "BIRD"
	COELHO   PetSpecies = "RABBIT"
)

// Pet struct
type Pet struct {
	ID      int
	Name    string
	Species PetSpecies
	Breed   string
	Size    string
}
