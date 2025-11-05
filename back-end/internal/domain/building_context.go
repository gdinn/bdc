package domain

import (
	"bdc/internal/models"
)

type CreateBuildingRequest struct {
	Name string `json:"name" validate:"required"`
}

func CreateBuildingData(req *CreateBuildingRequest) *models.Building {
	return &models.Building{
		Name: req.Name,
	}
}

type CreateBuildingResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *models.Building `json:"data,omitempty"`
	Error   string           `json:"error,omitempty"`
}
