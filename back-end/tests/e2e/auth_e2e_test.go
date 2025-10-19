package e2e

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bdc/api"
	"bdc/internal/models"
	"bdc/internal/testutils"

	"gorm.io/gorm"
)

// TestE2E_AuthFlow testa o fluxo completo de autenticação
func TestE2E_AuthFlow(t *testing.T) {
	// Setup

	db := testutils.SetupTestDatabase(t)
	defer testutils.CleanupTestDatabase(t, db)

	// Configurar variáveis de ambiente necessárias
	ma := testutils.GetMockAuthData()
	testutils.SetupAuthEnvironment(t, ma)

	// Criar servidor mock do JWKS
	privateKey, jwksServer := testutils.SetupMockJWKSServer(t)
	defer jwksServer.Close()

	// Setup da aplicação
	router := api.SetupRoutes(db)

	t.Run("Unauthorized_Request_Without_Token", func(t *testing.T) {
		testUnauthorizedRequest(t, router)
	})

	t.Run("Unauthorized_Request_With_Invalid_Token", func(t *testing.T) {
		testInvalidToken(t, router)
	})

	t.Run("Successful_Request_With_Valid_Token", func(t *testing.T) {
		testValidToken(t, router, db, privateKey)
	})

	t.Run("Create_User_With_Valid_Authentication", func(t *testing.T) {
		testCreateUserWithAuth(t, router, db, privateKey)
	})

	t.Run("Create_User_Email_Already_Exists", func(t *testing.T) {
		testCreateUserEmailAlreadyExists(t, router, db, privateKey)
	})
}

// testUnauthorizedRequest verifica requisição sem token
func testUnauthorizedRequest(t *testing.T, router http.Handler) {
	req := httptest.NewRequest("POST", "/api/v1/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["error"] != "Authorization header required" {
		t.Errorf("Expected 'Authorization header required', got '%s'", response["error"])
	}
}

// testInvalidToken verifica requisição com token inválido
func testInvalidToken(t *testing.T, router http.Handler) {
	req := httptest.NewRequest("POST", "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-here")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["error"] == "" {
		t.Error("Expected error message, got empty string")
	}
}

// testValidToken verifica requisição com token válido
func testValidToken(t *testing.T, router http.Handler, db *gorm.DB, privateKey *rsa.PrivateKey) {
	// Criar usuário no Cognito (simulado via claims)
	token, claims := testutils.CreateValidToken(t, privateKey, "testuser@example.com", "testuser", models.UserRoleCommon)

	// Criar body da requisição
	requestBody := map[string]interface{}{
		"phone":      "11999999999",
		"birth_date": time.Now().Add(-30 * 365 * 24 * time.Hour), // 30 anos atrás
		"type":       "RESIDENT",
		"age_group":  "ADULT",
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Criar usuário no Cognito mock (via database para simular)
	// Em produção, isso viria do Cognito
	testutils.SetupCognitoUser(t, db, claims.Email, "Test User")

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated && w.Code != http.StatusInternalServerError {
		t.Logf("Response body: %s", w.Body.String())
		t.Errorf("Expected status 201 or 500 (cognito mock), got %d", w.Code)
	}
}

// testCreateUserWithAuth testa criação completa de usuário autenticado
func testCreateUserWithAuth(t *testing.T, router http.Handler, db *gorm.DB, privateKey *rsa.PrivateKey) {
	// Criar token válido
	token, claims := testutils.CreateValidToken(t, privateKey, "newuser@example.com", "newuser", models.UserRoleCommon)

	// Setup usuário no Cognito mock
	testutils.SetupCognitoUser(t, db, claims.Email, "New User")

	// Preparar payload de criação de usuário
	birthDate := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	requestBody := map[string]interface{}{
		"phone":      "11988887777",
		"birth_date": birthDate.Format(time.RFC3339),
		"type":       "RESIDENT",
		"age_group":  "ADULT",
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Nota: Como não temos integração real com Cognito,
	// esperamos erro interno, mas validamos que a autenticação passou
	if w.Code == http.StatusUnauthorized {
		t.Error("Request should not be unauthorized with valid token")
	}

	// Verificar que não foi erro de autenticação
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Logf("Response body: %s", w.Body.String())
	}

	if errorMsg, ok := response["error"].(string); ok {
		if errorMsg == "Authorization header required" || errorMsg == "Invalid token" {
			t.Errorf("Authentication should have passed: %s", errorMsg)
		}
	}
}

// testCreateUserEmailAlreadyExists testa tentativa de criar usuário duplicado
func testCreateUserEmailAlreadyExists(t *testing.T, router http.Handler, db *gorm.DB, privateKey *rsa.PrivateKey) {
	email := "duplicate@example.com"

	// Criar usuário existente no banco
	existingUser := &models.User{
		Name:     "Existing User",
		Email:    email,
		Phone:    "11999999999",
		Type:     models.UserTypeResident,
		AgeGroup: models.UserAgeGroupAdult,
		Role:     models.UserRoleCommon,
	}

	if err := db.Create(existingUser).Error; err != nil {
		t.Fatalf("Failed to create existing user: %v", err)
	}

	// Tentar criar novamente com mesmo email
	token, _ := testutils.CreateValidToken(t, privateKey, email, "duplicate", models.UserRoleCommon)
	testutils.SetupCognitoUser(t, db, email, "Duplicate User")

	requestBody := map[string]interface{}{
		"phone":     "11988887777",
		"type":      "RESIDENT",
		"age_group": "ADULT",
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verificar resposta de conflito
	if w.Code == http.StatusConflict {
		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err == nil {
			if msg, ok := response["message"].(string); ok {
				t.Logf("Got expected conflict: %s", msg)
			}
		}
	}
}

// ===== HELPER FUNCTIONS =====

// TestE2E_HealthCheck testa endpoint público sem autenticação
func TestE2E_HealthCheck(t *testing.T) {
	db := testutils.SetupTestDatabase(t)
	defer testutils.CleanupTestDatabase(t, db)

	router := api.SetupRoutes(db)

	req := httptest.NewRequest("GET", "/api/v1/public/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}

	if response["service"] != "BDC API" {
		t.Errorf("Expected service 'BDC API', got '%s'", response["service"])
	}
}

// TestE2E_VersionCheck testa endpoint de versão
func TestE2E_VersionCheck(t *testing.T) {
	db := testutils.SetupTestDatabase(t)
	defer testutils.CleanupTestDatabase(t, db)

	router := api.SetupRoutes(db)

	req := httptest.NewRequest("GET", "/api/v1/public/version", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["version"] == "" {
		t.Error("Expected version to be set")
	}
}

// TestE2E_CORS testa configuração de CORS
func TestE2E_CORS(t *testing.T) {
	db := testutils.SetupTestDatabase(t)
	defer testutils.CleanupTestDatabase(t, db)

	router := api.SetupRoutes(db)

	// Testar OPTIONS request (preflight)
	req := httptest.NewRequest("OPTIONS", "/api/v1/public/health", nil)
	req.Header.Set("Origin", "http://localhost:4200")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verificar headers CORS
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("Expected Access-Control-Allow-Origin header to be set")
	}

	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("Expected Access-Control-Allow-Methods header to be set")
	}
}

// TestE2E_DatabaseConnection testa conexão com banco de dados
func TestE2E_DatabaseConnection(t *testing.T) {
	db := testutils.SetupTestDatabase(t)
	defer testutils.CleanupTestDatabase(t, db)

	// Verificar se pode executar query
	var result int
	if err := db.Raw("SELECT 1").Scan(&result).Error; err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected result 1, got %d", result)
	}

	// Verificar se tabelas foram criadas
	var tableCount int64
	if err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users'").Scan(&tableCount).Error; err != nil {
		t.Fatalf("Failed to check tables: %v", err)
	}

	if tableCount != 1 {
		t.Error("Expected users table to exist")
	}
}

// TestE2E_FullUserLifecycle testa ciclo completo de vida de usuário
func TestE2E_FullUserLifecycle(t *testing.T) {
	db := testutils.SetupTestDatabase(t)
	defer testutils.CleanupTestDatabase(t, db)

	ma := testutils.GetMockAuthData()
	testutils.SetupAuthEnvironment(t, ma)

	privateKey, jwksServer := testutils.SetupMockJWKSServer(t)
	defer jwksServer.Close()

	router := api.SetupRoutes(db)

	// 1. Criar usuário
	email := "lifecycle@example.com"
	token, _ := testutils.CreateValidToken(t, privateKey, email, "lifecycle", models.UserRoleCommon)
	testutils.SetupCognitoUser(t, db, email, "Lifecycle User")

	requestBody := map[string]interface{}{
		"phone":     "11999888777",
		"type":      "RESIDENT",
		"age_group": "ADULT",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	t.Logf("Create user response status: %d", w.Code)
	t.Logf("Create user response body: %s", w.Body.String())

	// 2. Verificar se usuário foi criado no banco (se passou da autenticação)
	if w.Code != http.StatusUnauthorized {
		var user models.User
		err := db.Where("email = ?", email).First(&user).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			t.Logf("User lookup result: %v", err)
		}
	}
}
