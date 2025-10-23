package domain

import "bdc/internal/models"

type CreateApartmentRequest struct {
	Name     string `json:"name" validate:"required"`
	Building string `json:"building" validate:"required"`
}

func CreateApartmentData(req *CreateApartmentRequest) *models.Apartment {
	return &models.Apartment{
		Name:     req.Name,
		Building: req.Building,
	}
}

type CreateApartmentResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *models.Apartment `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
}
