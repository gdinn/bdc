package interfaces

import (
	"bdc/internal/models"
)

type ApartmentRepositoryInterface interface {
	Create(apartment *models.Apartment) (*models.Apartment, error)
	IsApartmentExists(apartment *models.Apartment) (bool, error)
}
