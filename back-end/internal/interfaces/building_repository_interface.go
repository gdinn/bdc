package interfaces

import (
	"bdc/internal/models"
)

type BuildingRepositoryInterface interface {
	Create(building *models.Building) (*models.Building, error)
	IsBuildingExists(building *models.Building) (bool, error)
}
