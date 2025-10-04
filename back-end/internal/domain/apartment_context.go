package domain

import "bdc/internal/models"

type CreateApartmentRequest struct {
	Number   string `json:"number" validate:"required"`
	Building string `json:"building" validate:"required"`
}

func CreateApartmentData(req *CreateApartmentRequest) *models.Apartment {
	return &models.Apartment{
		Number:   req.Number,
		Building: req.Building,
	}
}

type CreateApartmentResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *models.Apartment `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
}
