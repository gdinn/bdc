package repositories

import (
	"bdc/internal/interfaces"
	"bdc/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type BuildingRepository struct {
	db *gorm.DB
}

const (
	ErrCheckingBuildingExistence = "error checking building existence"
)

func NewBuildingRepository(db *gorm.DB) interfaces.BuildingRepositoryInterface {
	return &BuildingRepository{
		db: db,
	}
}

func (r *BuildingRepository) Create(building *models.Building) (*models.Building, error) {
	if err := r.db.Create(building).Error; err != nil {
		return nil, fmt.Errorf("error creating building: %w", err)
	}
	return building, nil
}

func (r *BuildingRepository) IsBuildingExists(building *models.Building) (bool, error) {
	var count int64
	err := r.db.Model(&models.Building{}).
		Where("name = ?  AND deleted_at IS NULL", building.Name).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", ErrCheckingBuildingExistence, err)
	}

	return count > 0, nil
}
