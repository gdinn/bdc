package repositories

import (
	"bdc/internal/interfaces"
	"bdc/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type ApartmentRepository struct {
	db *gorm.DB
}

const (
	ErrCheckingApartmentExistence = "error checking apartment existence"
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
		Where("number = ? AND building = ? AND deleted_at IS NULL", apartment.Number, apartment.Building).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", ErrCheckingApartmentExistence, err)
	}

	return count > 0, nil
}
