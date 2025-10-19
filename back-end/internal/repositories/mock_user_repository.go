package repositories

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bdc/internal/domain"
	"bdc/internal/interfaces"
	"bdc/internal/models"
)

// MockUserRepository implementa UserRepositoryInterface para testes
type MockUserRepository struct {
	// Maps para simular o estado do banco de dados
	users     map[string]*models.User // Key: email
	usersByID map[uint]*models.User   // Key: ID
	nextID    uint

	// Flags para simular erros
	ShouldFailCreate        bool
	ShouldFailGetByEmail    bool
	ShouldFailIsEmailExists bool
	ShouldFailGetByID       bool

	// Contadores para verificar quantas vezes os métodos foram chamados
	CreateCallCount        int
	GetByEmailCallCount    int
	IsEmailExistsCallCount int
	GetByIDCallCount       int
}

// Verificação em tempo de compilação se MockUserRepository implementa a interface
var _ interfaces.UserRepositoryInterface = (*MockUserRepository)(nil)

// NewMockUserRepository cria uma nova instância do mock repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:     make(map[string]*models.User),
		usersByID: make(map[uint]*models.User),
		nextID:    1,
	}
}

// AddTestUser adiciona um usuário de teste ao mock
func (m *MockUserRepository) AddTestUser(email, name string, userType models.UserType, ageGroup models.UserAgeGroup, role models.UserRole) *models.User {
	user := &models.User{
		Name:     name,
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Type:     userType,
		AgeGroup: ageGroup,
		Role:     role,
	}

	// Simular ID gerado pelo banco
	user.BaseModel.ID = uuid.New()

	m.users[user.Email] = user
	id := m.nextID
	m.nextID++
	m.usersByID[id] = user

	return user
}

// Create simula criação de usuário no banco
func (m *MockUserRepository) Create(user *models.User) (*models.User, error) {
	m.CreateCallCount++

	if m.ShouldFailCreate {
		return nil, fmt.Errorf("mock error: failed to create user")
	}

	// Replica validação gorm de não nulidade
	if user.Email == "" {
		return nil, fmt.Errorf("mock gorm: email cannot be empty")
	}

	// Replica validação gorm de existência
	normalizedEmail := strings.ToLower(strings.TrimSpace(user.Email))
	if _, exists := m.users[normalizedEmail]; exists {
		return nil, fmt.Errorf("mock gorm: email already exists")
	}

	// Simular ID gerado pelo banco
	user.BaseModel.ID = uuid.New()

	// Adicionar ao mock storage
	m.users[user.Email] = user
	id := m.nextID
	m.nextID++
	m.usersByID[id] = user

	return user, nil
}

// GetByEmail simula busca de usuário por email
func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	m.GetByEmailCallCount++

	if m.ShouldFailGetByEmail {
		return nil, fmt.Errorf("mock error: failed to get user by email")
	}

	email = strings.ToLower(strings.TrimSpace(email))

	user, exists := m.users[email]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	// Criar uma cópia para evitar modificações acidentais
	userCopy := *user
	return &userCopy, nil
}

// IsEmailExists simula verificação de existência de email
func (m *MockUserRepository) IsEmailExists(email string) (bool, error) {
	m.IsEmailExistsCallCount++

	if m.ShouldFailIsEmailExists {
		return false, fmt.Errorf("mock error: failed to check email existence")
	}

	if email == "" {
		return false, fmt.Errorf("%s: %w", "error checking email existence", domain.ErrEmailEmpty)
	}

	email = strings.ToLower(strings.TrimSpace(email))

	_, exists := m.users[email]
	return exists, nil
}

// GetByID simula busca de usuário por ID
func (m *MockUserRepository) GetByID(id uint) (*models.User, error) {
	m.GetByIDCallCount++

	if m.ShouldFailGetByID {
		return nil, fmt.Errorf("mock error: failed to get user by ID")
	}

	user, exists := m.usersByID[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	// Criar uma cópia para evitar modificações acidentais
	userCopy := *user
	return &userCopy, nil
}

// Reset reseta o estado do mock para um novo teste
func (m *MockUserRepository) Reset() {
	m.users = make(map[string]*models.User)
	m.usersByID = make(map[uint]*models.User)
	m.nextID = 1
	m.ShouldFailCreate = false
	m.ShouldFailGetByEmail = false
	m.ShouldFailIsEmailExists = false
	m.ShouldFailGetByID = false
	m.CreateCallCount = 0
	m.GetByEmailCallCount = 0
	m.IsEmailExistsCallCount = 0
	m.GetByIDCallCount = 0
}

// GetUserCount retorna o número de usuários no mock (método auxiliar para testes)
func (m *MockUserRepository) GetUserCount() int {
	return len(m.users)
}

// UserExists verifica se um usuário existe pelo email (método auxiliar para testes)
func (m *MockUserRepository) UserExists(email string) bool {
	email = strings.ToLower(strings.TrimSpace(email))
	_, exists := m.users[email]
	return exists
}
