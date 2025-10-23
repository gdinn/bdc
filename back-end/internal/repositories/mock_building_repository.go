package repositories

import (
	"bdc/internal/interfaces"
	"bdc/internal/models"
	"fmt"

	"github.com/google/uuid"
)

// MockBuildingRepository implementa BuildingRepositoryInterface para testes
type MockBuildingRepository struct {
	// Maps para simular o estado do banco de dados
	buildings map[uuid.UUID]*models.Building
	nextID    uint

	// Flags para simular erros
	ShouldFailCreate           bool
	ShouldFailIsBuildingExists bool

	// Contadores para verificar quantas vezes os métodos foram chamados
	CreateCallCount           int
	IsBuildingExistsCallCount int
}

// Verificação em tempo de compilação se MockBuildingRepository implementa a interface
var _ interfaces.BuildingRepositoryInterface = (*MockBuildingRepository)(nil)

func NewMockBuildingRepository() *MockBuildingRepository {
	return &MockBuildingRepository{
		buildings: make(map[uuid.UUID]*models.Building), // building.BaseModel.ID as key
		nextID:    1,
	}
}

func (m *MockBuildingRepository) Create(building *models.Building) (*models.Building, error) {
	m.CreateCallCount++

	if m.ShouldFailCreate {
		return nil, fmt.Errorf("mock error: failed to create building")
	}

	// Replica validação gorm de não nulidade
	if building.Name == "" {
		return nil, fmt.Errorf("mock gorm: name cannot be empty")
	}

	building.BaseModel.ID = uuid.New()
	m.buildings[building.BaseModel.ID] = building

	return building, nil
}

func (m *MockBuildingRepository) IsBuildingExists(building *models.Building) (bool, error) {
	for _, b := range m.buildings {
		if b.Name == building.Name {
			return true, nil
		}
	}
	return false, nil
}
