package services

import (
	"bdc/internal/domain"
	"bdc/internal/interfaces"
	"bdc/internal/models"
	"bdc/internal/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type BuildingService struct {
	buildingRepo interfaces.BuildingRepositoryInterface
}

const (
	ErrCreatingBuilding = "failed to create building"
)

func NewBuildingService(buildingRepo interfaces.BuildingRepositoryInterface) *BuildingService {
	return &BuildingService{
		buildingRepo: buildingRepo,
	}
}

func (s *BuildingService) CreateBuilding(building *models.Building, claims *models.UserClaims) (*models.Building, error) {
	if err := utils.ValidateManagerRole(claims); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCreatingBuilding, err)
	}

	if err := s.validateBuildingIsNew(building); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCreatingBuilding, err)
	}

	createdBuilding, err := s.buildingRepo.Create(building)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCreatingBuilding, err)
	}

	return createdBuilding, nil

}

func (s *BuildingService) validateBuildingIsNew(building *models.Building) error {
	buildingExists, err := s.buildingRepo.IsBuildingExists(building)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error checking building existence: %w", err)
	}

	if buildingExists {
		return fmt.Errorf("building.Name is %s: %w", building.Name, domain.ErrBuildingAlreadyExists)
	}

	return nil
}
