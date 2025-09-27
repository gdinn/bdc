package repositories

import (
	"bdc/internal/interfaces"
	"bdc/internal/models"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// MockCognitoRepository implementa CognitoRepositoryInterface para testes
type MockCognitoRepository struct {
	// Maps para simular o estado do Cognito
	users         map[string]*cognitoidentityprovider.AdminGetUserOutput
	disabledUsers map[string]bool

	// Flags para simular erros
	ShouldFailUpdate  bool
	ShouldFailGet     bool
	ShouldFailDelete  bool
	ShouldFailDisable bool
	ShouldFailEnable  bool

	// Contadores para verificar quantas vezes os métodos foram chamados
	UpdateCallCount  int
	GetCallCount     int
	DeleteCallCount  int
	DisableCallCount int
	EnableCallCount  int
}

// Verificação em tempo de compilação se MockCognitoRepository implementa a interface
var _ interfaces.CognitoRepositoryInterface = (*MockCognitoRepository)(nil)

// NewMockCognitoRepository cria uma nova instância do mock repository
func NewMockCognitoRepository() *MockCognitoRepository {
	return &MockCognitoRepository{
		users:         make(map[string]*cognitoidentityprovider.AdminGetUserOutput),
		disabledUsers: make(map[string]bool),
	}
}

// AddTestUser adiciona um usuário de teste ao mock
func (m *MockCognitoRepository) AddTestUser(email, name string, role models.UserRole) {
	m.users[email] = &cognitoidentityprovider.AdminGetUserOutput{
		Username: aws.String(email),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("name"),
				Value: aws.String(name),
			},
			{
				Name:  aws.String("custom:role"),
				Value: aws.String(string(role)),
			},
		},
		Enabled: true,
	}
}

// UpdateUserInCognito simula atualização de usuário no Cognito
func (m *MockCognitoRepository) UpdateUserInCognito(user *models.User) error {
	m.UpdateCallCount++

	if m.ShouldFailUpdate {
		return fmt.Errorf("mock error: failed to update user")
	}

	// Verificar se o usuário existe
	cognitoUser, exists := m.users[user.Email]
	if !exists {
		return fmt.Errorf("user not found in cognito")
	}

	// Atualizar os atributos do usuário
	for i := range cognitoUser.UserAttributes {
		attr := &cognitoUser.UserAttributes[i]
		switch *attr.Name {
		case "name":
			attr.Value = aws.String(user.Name)
		case "custom:role":
			attr.Value = aws.String(string(user.Role))
		}
	}

	return nil
}

// GetUserFromCognito simula busca de usuário no Cognito
func (m *MockCognitoRepository) GetUserFromCognito(username string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	m.GetCallCount++

	if m.ShouldFailGet {
		return nil, fmt.Errorf("mock error: failed to get user")
	}

	user, exists := m.users[username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Criar uma cópia para evitar modificações acidentais
	userCopy := *user
	return &userCopy, nil
}

// DeleteUser simula exclusão de usuário do Cognito
func (m *MockCognitoRepository) DeleteUser(email string) error {
	m.DeleteCallCount++

	if m.ShouldFailDelete {
		return fmt.Errorf("mock error: failed to delete user")
	}

	if _, exists := m.users[email]; !exists {
		return fmt.Errorf("user not found")
	}

	delete(m.users, email)
	delete(m.disabledUsers, email)
	return nil
}

// DisableUser simula desabilitação de usuário no Cognito
func (m *MockCognitoRepository) DisableUser(email string) error {
	m.DisableCallCount++

	if m.ShouldFailDisable {
		return fmt.Errorf("mock error: failed to disable user")
	}

	if _, exists := m.users[email]; !exists {
		return fmt.Errorf("user not found")
	}

	m.disabledUsers[email] = true
	// Atualizar o enabled status
	if user, exists := m.users[email]; exists {
		user.Enabled = false
	}

	return nil
}

// EnableUser simula habilitação de usuário no Cognito
func (m *MockCognitoRepository) EnableUser(email string) error {
	m.EnableCallCount++

	if m.ShouldFailEnable {
		return fmt.Errorf("mock error: failed to enable user")
	}

	if _, exists := m.users[email]; !exists {
		return fmt.Errorf("user not found")
	}

	delete(m.disabledUsers, email)
	// Atualizar o enabled status
	if user, exists := m.users[email]; exists {
		user.Enabled = true
	}

	return nil
}

// IsUserDisabled verifica se um usuário está desabilitado (método auxiliar para testes)
func (m *MockCognitoRepository) IsUserDisabled(email string) bool {
	return m.disabledUsers[email]
}

// Reset reseta o estado do mock para um novo teste
func (m *MockCognitoRepository) Reset() {
	m.users = make(map[string]*cognitoidentityprovider.AdminGetUserOutput)
	m.disabledUsers = make(map[string]bool)
	m.ShouldFailUpdate = false
	m.ShouldFailGet = false
	m.ShouldFailDelete = false
	m.ShouldFailDisable = false
	m.ShouldFailEnable = false
	m.UpdateCallCount = 0
	m.GetCallCount = 0
	m.DeleteCallCount = 0
	m.DisableCallCount = 0
	m.EnableCallCount = 0
}
