package domain

import (
	"bdc/internal/models"

	"github.com/google/uuid"
)

type CreateApartmentRequest struct {
	Name       string    `json:"name" validate:"required"`
	BuildingID uuid.UUID `json:"building_id" validate:"required"`
}

func CreateApartmentData(req *CreateApartmentRequest) *models.Apartment {
	return &models.Apartment{
		Name:       req.Name,
		BuildingID: req.BuildingID,
	}
}

type CreateApartmentResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *models.Apartment `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
}

type GetApartmentResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *models.Apartment `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
}
