package testutils

import (
	"fmt"
	"testing"

	"gorm.io/gorm"
)

type MockAuthData struct {
	AWS_REGION               string
	AWS_COGNITO_USER_POOL_ID string
}

func GetMockAuthData() *MockAuthData {
	return &MockAuthData{
		AWS_REGION:               "us-east-1",
		AWS_COGNITO_USER_POOL_ID: "us-east-1_test123",
	}
}

// setupAuthEnvironment configura variáveis de ambiente para testes
func SetupAuthEnvironment(t *testing.T, ma *MockAuthData) {
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)
}

// setupCognitoUser simula usuário no Cognito (para testes)
func SetupCognitoUser(t *testing.T, db *gorm.DB, email, name string) {
	// Em testes E2E reais, você usaria um mock do serviço Cognito
	// Por ora, apenas garantimos que não haja erro de contexto
	t.Logf("Mock Cognito user created: %s (%s)", email, name)
}

func GetJwksUrl(ma *MockAuthData) string {
	return fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", ma.AWS_REGION, ma.AWS_COGNITO_USER_POOL_ID)
}

func GetIssuerUrl(ma *MockAuthData) string {
	return fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", ma.AWS_REGION, ma.AWS_COGNITO_USER_POOL_ID)
}
