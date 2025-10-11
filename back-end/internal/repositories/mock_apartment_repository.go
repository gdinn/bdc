package repositories

import (
	"bdc/internal/interfaces"
	"bdc/internal/models"
	"fmt"
)

// MockApartmentRepository implementa ApartmentRepositoryInterface para testes
type MockApartmentRepository struct {
	// Maps para simular o estado do banco de dados
	apartments map[string]*models.Apartment // Key: %s:%s (building:number)
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
	if apartment.Number == "" {
		return nil, fmt.Errorf("number cannot be empty")
	}
	if apartment.Building == "" {
		return nil, fmt.Errorf("building cannot be empty")
	}

	// Replica validação gorm de existência

	/*
		TODO: Inserir composite primary key p/ validar existência
		TODO: Atualizar diagrama UML
	*/

}

func (r *MockApartmentRepository) IsApartmentExists(apartment *models.Apartment) (bool, error) {
}
