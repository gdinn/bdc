package domain

import (
	"testing"
	"time"

	"bdc/internal/models"
)

func TestCreateUserData_Success(t *testing.T) {
	// Arrange
	birthDate := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	req := &CreateUserRequest{
		Phone:     "11999999999",
		BirthDate: &birthDate,
		Type:      models.UserTypeResident,
		AgeGroup:  models.UserAgeGroupAdult,
	}

	// Act
	user := CreateUserData(req)

	// Assert
	if user == nil {
		t.Error("Expected user to be created, but got nil")
	}

	if user.Phone != "11999999999" {
		t.Errorf("Expected phone to be '11999999999', but got '%s'", user.Phone)
	}

	if user.BirthDate == nil {
		t.Error("Expected birth date to be set, but got nil")
	} else if !user.BirthDate.Equal(birthDate) {
		t.Errorf("Expected birth date to be %v, but got %v", birthDate, *user.BirthDate)
	}

	if user.Type != models.UserTypeResident {
		t.Errorf("Expected type to be %s, but got %s", models.UserTypeResident, user.Type)
	}

	if user.AgeGroup != models.UserAgeGroupAdult {
		t.Errorf("Expected age group to be %s, but got %s", models.UserAgeGroupAdult, user.AgeGroup)
	}

	if user.Role != models.UserRoleCommon {
		t.Errorf("Expected role to be %s (default), but got %s", models.UserRoleCommon, user.Role)
	}
}

func TestCreateUserData_EmptyPhone(t *testing.T) {
	// Arrange
	req := &CreateUserRequest{
		Phone:    "", // Empty phone
		Type:     models.UserTypeExternal,
		AgeGroup: models.UserAgeGroupAdult,
	}

	// Act
	user := CreateUserData(req)

	// Assert
	if user == nil {
		t.Error("Expected user to be created, but got nil")
	}

	if user.Phone != "" {
		t.Errorf("Expected phone to be empty, but got '%s'", user.Phone)
	}

	if user.Type != models.UserTypeExternal {
		t.Errorf("Expected type to be %s, but got %s", models.UserTypeExternal, user.Type)
	}

	if user.Role != models.UserRoleCommon {
		t.Errorf("Expected role to be %s (default), but got %s", models.UserRoleCommon, user.Role)
	}
}

func TestCreateUserData_NilBirthDate(t *testing.T) {
	// Arrange
	req := &CreateUserRequest{
		Phone:     "11988887777",
		BirthDate: nil, // Nil birth date
		Type:      models.UserTypeResident,
		AgeGroup:  models.UserAgeGroupChild,
	}

	// Act
	user := CreateUserData(req)

	// Assert
	if user == nil {
		t.Error("Expected user to be created, but got nil")
	}

	if user.BirthDate != nil {
		t.Errorf("Expected birth date to be nil, but got %v", *user.BirthDate)
	}

	if user.Type != models.UserTypeResident {
		t.Errorf("Expected type to be %s, but got %s", models.UserTypeResident, user.Type)
	}

	if user.AgeGroup != models.UserAgeGroupChild {
		t.Errorf("Expected age group to be %s, but got %s", models.UserAgeGroupChild, user.AgeGroup)
	}
}

func TestCreateUserData_AllUserTypesAndAgeGroups(t *testing.T) {
	// Test cases para diferentes combinações de Type e AgeGroup
	testCases := []struct {
		name     string
		userType models.UserType
		ageGroup models.UserAgeGroup
	}{
		{
			name:     "Resident Adult",
			userType: models.UserTypeResident,
			ageGroup: models.UserAgeGroupAdult,
		},
		{
			name:     "Resident Child",
			userType: models.UserTypeResident,
			ageGroup: models.UserAgeGroupChild,
		},
		{
			name:     "External Adult",
			userType: models.UserTypeExternal,
			ageGroup: models.UserAgeGroupAdult,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := &CreateUserRequest{
				Phone:    "11999888777",
				Type:     tc.userType,
				AgeGroup: tc.ageGroup,
			}

			// Act
			user := CreateUserData(req)

			// Assert
			if user == nil {
				t.Error("Expected user to be created, but got nil")
			}

			if user.Type != tc.userType {
				t.Errorf("Expected type to be %s, but got %s", tc.userType, user.Type)
			}

			if user.AgeGroup != tc.ageGroup {
				t.Errorf("Expected age group to be %s, but got %s", tc.ageGroup, user.AgeGroup)
			}

			// Sempre deve ter role padrão
			if user.Role != models.UserRoleCommon {
				t.Errorf("Expected role to be %s (default), but got %s", models.UserRoleCommon, user.Role)
			}
		})
	}
}

func TestCreateUserData_DefaultRoleIsAlwaysCommon(t *testing.T) {
	// Arrange
	req := &CreateUserRequest{
		Phone:    "11987654321",
		Type:     models.UserTypeExternal,
		AgeGroup: models.UserAgeGroupAdult,
	}

	// Act
	user := CreateUserData(req)

	// Assert
	if user.Role != models.UserRoleCommon {
		t.Errorf("Expected default role to always be %s, but got %s", models.UserRoleCommon, user.Role)
	}

	// Verificar que outros campos não relacionados ao role estão vazios (não preenchidos)
	if user.Name != "" {
		t.Errorf("Expected name to be empty (not set by CreateUserData), but got '%s'", user.Name)
	}

	if user.Email != "" {
		t.Errorf("Expected email to be empty (not set by CreateUserData), but got '%s'", user.Email)
	}
}

func TestNewCreateUserContext_Success(t *testing.T) {
	// Arrange
	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "testuser",
		Role:     "COMMON",
		Sub:      "user-sub-123",
		TokenUse: "access",
	}

	user := &models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Phone:    "11999999999",
		Type:     models.UserTypeResident,
		AgeGroup: models.UserAgeGroupAdult,
		Role:     models.UserRoleCommon,
	}

	// Act
	ctx := NewCreateUserContext(claims, user)

	// Assert
	if ctx == nil {
		t.Error("Expected context to be created, but got nil")
	}

	if ctx.Claims == nil {
		t.Error("Expected claims to be set, but got nil")
	}

	if ctx.User == nil {
		t.Error("Expected user to be set, but got nil")
	}

	if ctx.Claims.Email != "test@example.com" {
		t.Errorf("Expected claims email to be 'test@example.com', but got '%s'", ctx.Claims.Email)
	}

	if ctx.User.Email != "test@example.com" {
		t.Errorf("Expected user email to be 'test@example.com', but got '%s'", ctx.User.Email)
	}
}

func TestNewCreateUserContext_NilClaims(t *testing.T) {
	// Arrange
	user := &models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Type:     models.UserTypeResident,
		AgeGroup: models.UserAgeGroupAdult,
		Role:     models.UserRoleCommon,
	}

	// Act
	ctx := NewCreateUserContext(nil, user)

	// Assert
	if ctx == nil {
		t.Error("Expected context to be created, but got nil")
	}

	if ctx.Claims != nil {
		t.Error("Expected claims to be nil, but got a value")
	}

	if ctx.User == nil {
		t.Error("Expected user to be set, but got nil")
	}

	if ctx.User.Email != "test@example.com" {
		t.Errorf("Expected user email to be 'test@example.com', but got '%s'", ctx.User.Email)
	}
}

func TestNewCreateUserContext_NilUser(t *testing.T) {
	// Arrange
	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "testuser",
		Role:     "COMMON",
	}

	// Act
	ctx := NewCreateUserContext(claims, nil)

	// Assert
	if ctx == nil {
		t.Error("Expected context to be created, but got nil")
	}

	if ctx.Claims == nil {
		t.Error("Expected claims to be set, but got nil")
	}

	if ctx.User != nil {
		t.Error("Expected user to be nil, but got a value")
	}

	if ctx.Claims.Email != "test@example.com" {
		t.Errorf("Expected claims email to be 'test@example.com', but got '%s'", ctx.Claims.Email)
	}
}

func TestNewCreateUserContext_BothNil(t *testing.T) {
	// Arrange & Act
	ctx := NewCreateUserContext(nil, nil)

	// Assert
	if ctx == nil {
		t.Error("Expected context to be created, but got nil")
	}

	if ctx.Claims != nil {
		t.Error("Expected claims to be nil, but got a value")
	}

	if ctx.User != nil {
		t.Error("Expected user to be nil, but got a value")
	}
}

func TestCreateUserRequest_Struct(t *testing.T) {
	// Test para verificar que a struct CreateUserRequest pode ser criada corretamente
	birthDate := time.Date(1985, 12, 25, 0, 0, 0, 0, time.UTC)

	req := CreateUserRequest{
		Phone:     "11987654321",
		BirthDate: &birthDate,
		Type:      models.UserTypeExternal,
		AgeGroup:  models.UserAgeGroupAdult,
	}

	// Verificar que todos os campos foram definidos corretamente
	if req.Phone != "11987654321" {
		t.Errorf("Expected phone to be '11987654321', but got '%s'", req.Phone)
	}

	if req.BirthDate == nil {
		t.Error("Expected birth date to be set, but got nil")
	} else if !req.BirthDate.Equal(birthDate) {
		t.Errorf("Expected birth date to be %v, but got %v", birthDate, *req.BirthDate)
	}

	if req.Type != models.UserTypeExternal {
		t.Errorf("Expected type to be %s, but got %s", models.UserTypeExternal, req.Type)
	}

	if req.AgeGroup != models.UserAgeGroupAdult {
		t.Errorf("Expected age group to be %s, but got %s", models.UserAgeGroupAdult, req.AgeGroup)
	}
}

func TestCreateUserContext_Struct(t *testing.T) {
	// Test para verificar que a struct CreateUserContext pode ser criada corretamente
	claims := &models.UserClaims{
		Email:    "admin@example.com",
		Username: "admin",
		Role:     "MANAGER",
	}

	user := &models.User{
		Name:     "Admin User",
		Email:    "admin@example.com",
		Type:     models.UserTypeResident,
		AgeGroup: models.UserAgeGroupAdult,
		Role:     models.UserRoleManager,
	}

	ctx := CreateUserContext{
		Claims: claims,
		User:   user,
	}

	// Verificar que todos os campos foram definidos corretamente
	if ctx.Claims == nil {
		t.Error("Expected claims to be set, but got nil")
	}

	if ctx.User == nil {
		t.Error("Expected user to be set, but got nil")
	}

	if ctx.Claims.Role != "MANAGER" {
		t.Errorf("Expected claims role to be 'MANAGER', but got '%s'", ctx.Claims.Role)
	}

	if ctx.User.Role != models.UserRoleManager {
		t.Errorf("Expected user role to be %s, but got %s", models.UserRoleManager, ctx.User.Role)
	}
}

func TestCreateUserData_CompleteWorkflow(t *testing.T) {
	// Test que simula um workflow completo de criação de usuário

	// Arrange - Simular dados vindos de uma requisição HTTP
	birthDate := time.Date(1992, 8, 10, 0, 0, 0, 0, time.UTC)
	req := &CreateUserRequest{
		Phone:     "11955443322",
		BirthDate: &birthDate,
		Type:      models.UserTypeResident,
		AgeGroup:  models.UserAgeGroupAdult,
	}

	claims := &models.UserClaims{
		Email:    "user@example.com",
		Username: "user123",
		Role:     "COMMON",
		Sub:      "user-123-456",
	}

	// Act - Converter request para user model
	user := CreateUserData(req)

	// Act - Criar contexto completo
	ctx := NewCreateUserContext(claims, user)

	// Assert - Verificar workflow completo
	if ctx == nil {
		t.Error("Expected context to be created, but got nil")
	}

	if ctx.User.Phone != req.Phone {
		t.Errorf("Expected user phone to match request phone '%s', but got '%s'", req.Phone, ctx.User.Phone)
	}

	if ctx.User.Type != req.Type {
		t.Errorf("Expected user type to match request type %s, but got %s", req.Type, ctx.User.Type)
	}

	if ctx.User.AgeGroup != req.AgeGroup {
		t.Errorf("Expected user age group to match request age group %s, but got %s", req.AgeGroup, ctx.User.AgeGroup)
	}

	if ctx.Claims.Email != "user@example.com" {
		t.Errorf("Expected claims email to be 'user@example.com', but got '%s'", ctx.Claims.Email)
	}

	// Role deve sempre ser Common ao criar usuário via CreateUserData
	if ctx.User.Role != models.UserRoleCommon {
		t.Errorf("Expected user role to be %s (default), but got %s", models.UserRoleCommon, ctx.User.Role)
	}
}
