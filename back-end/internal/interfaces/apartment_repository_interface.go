package interfaces

import (
	"bdc/internal/models"

	"github.com/google/uuid"
)

type ApartmentRepositoryInterface interface {
	Create(apartment *models.Apartment) (*models.Apartment, error)
	Get(apartmentID uuid.UUID) (*models.Apartment, error)
	IsApartmentExists(apartment *models.Apartment) (bool, error)
}
