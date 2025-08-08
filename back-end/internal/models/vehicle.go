package models

import "github.com/google/uuid"

type VehicleType string

const (
	VehicleTypeCar        VehicleType = "CAR"
	VehicleTypeMotorcycle VehicleType = "MOTORCYCLE"
)

type Vehicle struct {
	BaseModel
	Plate         string      `json:"plate" gorm:"uniqueIndex;not null;size:10" validate:"required,max=10"`
	Model         string      `json:"model" gorm:"not null;size:100" validate:"required,max=100"`
	Color         string      `json:"color" gorm:"not null;size:50" validate:"required,max=50"`
	Year          int         `json:"year" validate:"min=1900,max=2030"`
	ParkingNumber string      `json:"parking_number" gorm:"size:10" validate:"max=10"`
	Type          VehicleType `json:"type" gorm:"not null;-:migration" validate:"required"`

	// Chave estrangeira
	ApartmentID uuid.UUID `json:"apartment_id" gorm:"not null; type:uuid" validate:"required"`
	Apartment   Apartment `json:"apartment,omitempty" gorm:"foreignKey:ApartmentID"`
}
