package repositories

import (
	"bdc/internal/interfaces"
	"bdc/internal/models"
	"fmt"

	"github.com/google/uuid"
)

// MockApartmentRepository implementa ApartmentRepositoryInterface para testes
type MockApartmentRepository struct {
	// Maps para simular o estado do banco de dados
	apartments map[string]*models.Apartment // Key: %s:%s (building:name)
	nextID     uint

	// Flags para simular erros
	ShouldFailCreate            bool
	ShouldFailIsApartmentExists bool

	// Contadores para verificar quantas vezes os métodos foram chamados
	CreateCallCount            int
	IsApartmentExistsCallCount int
}

// Verificação em tempo de compilação se MockUserRepository implementa a interface
var _ interfaces.ApartmentRepositoryInterface = (*MockApartmentRepository)(nil)

func NewMockApartmentRepository() *MockApartmentRepository {
	return &MockApartmentRepository{
		apartments: make(map[string]*models.Apartment),
		nextID:     1,
	}
}

func (m *MockApartmentRepository) Create(apartment *models.Apartment) (*models.Apartment, error) {
	m.CreateCallCount++

	if m.ShouldFailCreate {
		return nil, fmt.Errorf("mock error: failed to create apartment")
	}

	// Replica validação gorm de não nulidade
	if apartment.Name == "" {
		return nil, fmt.Errorf("mock gorm: name cannot be empty")
	}

	apartment.BaseModel.ID = uuid.New()
	apartmentKey := m.GetApartmentKey(apartment)
	m.apartments[apartmentKey] = apartment

	return apartment, nil
}

func (m *MockApartmentRepository) IsApartmentExists(apartment *models.Apartment) (bool, error) {
	apartmentKey := m.GetApartmentKey(apartment)
	_, exists := m.apartments[apartmentKey]
	return exists, nil
}

func (m *MockApartmentRepository) GetApartmentKey(apartment *models.Apartment) string {
	return fmt.Sprintf("%v:%s", apartment.Building, apartment.Name)
}
