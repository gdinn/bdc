package repositories

import (
	"bdc/internal/interfaces"
	"bdc/internal/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApartmentRepository struct {
	db *gorm.DB
}

const (
	ErrCheckingApartmentExistence = "error checking apartment existence"
	ErrGettingApartmentById       = "error obtaing apartment"
)

func NewApartmentRepository(db *gorm.DB) interfaces.ApartmentRepositoryInterface {
	return &ApartmentRepository{
		db: db,
	}
}

func (r *ApartmentRepository) Create(apartment *models.Apartment) (*models.Apartment, error) {
	if err := r.db.Create(apartment).Error; err != nil {
		return nil, fmt.Errorf("error creating apartment: %w", err)
	}
	return apartment, nil
}

func (r *ApartmentRepository) IsApartmentExists(apartment *models.Apartment) (bool, error) {
	var count int64
	err := r.db.Model(&models.Apartment{}).
		Where("name = ? AND building_id = ? AND deleted_at IS NULL", apartment.Name, apartment.BuildingID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", ErrCheckingApartmentExistence, err)
	}

	return count > 0, nil
}

func (r *ApartmentRepository) Get(apartmentID uuid.UUID) (*models.Apartment, error) {
	var apartment models.Apartment
	err := r.db.
		Preload("Building").
		Preload("LegalRepresentative").
		Preload("Users").
		Preload("Vehicles").
		Preload("Pets").
		Preload("Bicycles").
		Where("id = ? AND deleted_at IS NULL", apartmentID).First(&apartment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: apartment not found", ErrGettingApartmentById)
		}
		return nil, fmt.Errorf("%s: %w", ErrGettingApartmentById, err)
	}

	return &apartment, nil
}
