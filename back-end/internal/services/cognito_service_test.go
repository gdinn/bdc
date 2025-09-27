package services

import (
	"testing"

	"bdc/internal/models"
	"bdc/internal/repositories"
)

func TestCognitoService_UpdateUserInCognito_Success(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockRepo)

	// Adicionar usuário de teste no mock
	mockRepo.AddTestUser("test@example.com", "Test User", models.UserRoleCommon)

	user := &models.User{
		Name:  "Updated Name",
		Email: "test@example.com",
		Role:  models.UserRoleAdvisor,
	}

	// Act
	err := cognitoService.UpdateUserInCognito(user)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if mockRepo.UpdateCallCount != 1 {
		t.Errorf("Expected UpdateUserInCognito to be called once, but was called %d times", mockRepo.UpdateCallCount)
	}

	if mockRepo.GetCallCount != 1 {
		t.Errorf("Expected GetUserFromCognito to be called once, but was called %d times", mockRepo.GetCallCount)
	}
}

func TestCognitoService_UpdateUserInCognito_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockRepo)

	user := &models.User{
		Name:  "Nonexistent User",
		Email: "nonexistent@example.com",
		Role:  models.UserRoleCommon,
	}

	// Act
	err := cognitoService.UpdateUserInCognito(user)

	// Assert
	if err == nil {
		t.Error("Expected an error for nonexistent user, but got nil")
	}

	if mockRepo.UpdateCallCount != 0 {
		t.Errorf("Expected UpdateUserInCognito not to be called, but was called %d times", mockRepo.UpdateCallCount)
	}

	if mockRepo.GetCallCount != 1 {
		t.Errorf("Expected GetUserFromCognito to be called once, but was called %d times", mockRepo.GetCallCount)
	}
}

func TestCognitoService_DeleteUserFromCognito_Success(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockRepo)

	email := "test@example.com"
	mockRepo.AddTestUser(email, "Test User", models.UserRoleCommon)

	// Act
	err := cognitoService.DeleteUserFromCognito(email)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if mockRepo.DeleteCallCount != 1 {
		t.Errorf("Expected DeleteUser to be called once, but was called %d times", mockRepo.DeleteCallCount)
	}

	if mockRepo.GetCallCount != 1 {
		t.Errorf("Expected GetUserFromCognito to be called once, but was called %d times", mockRepo.GetCallCount)
	}
}

func TestCognitoService_DeleteUserFromCognito_EmptyEmail(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockRepo)

	// Act
	err := cognitoService.DeleteUserFromCognito("")

	// Assert
	if err == nil {
		t.Error("Expected an error for empty email, but got nil")
	}

	if mockRepo.DeleteCallCount != 0 {
		t.Errorf("Expected DeleteUser not to be called, but was called %d times", mockRepo.DeleteCallCount)
	}
}

func TestCognitoService_DisableUserInCognito_Success(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockRepo)

	email := "test@example.com"
	mockRepo.AddTestUser(email, "Test User", models.UserRoleCommon)

	// Act
	err := cognitoService.DisableUserInCognito(email)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if mockRepo.DisableCallCount != 1 {
		t.Errorf("Expected DisableUser to be called once, but was called %d times", mockRepo.DisableCallCount)
	}

	if !mockRepo.IsUserDisabled(email) {
		t.Error("Expected user to be disabled, but user is still enabled")
	}
}

func TestCognitoService_EnableUserInCognito_Success(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockRepo)

	email := "test@example.com"
	mockRepo.AddTestUser(email, "Test User", models.UserRoleCommon)

	// Primeiro desabilitar o usuário
	_ = mockRepo.DisableUser(email)

	// Act
	err := cognitoService.EnableUserInCognito(email)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if mockRepo.EnableCallCount != 1 {
		t.Errorf("Expected EnableUser to be called once, but was called %d times", mockRepo.EnableCallCount)
	}

	if mockRepo.IsUserDisabled(email) {
		t.Error("Expected user to be enabled, but user is still disabled")
	}
}

func TestCognitoService_GetUserInCognito_Success(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockRepo)

	username := "test@example.com"
	mockRepo.AddTestUser(username, "Test User", models.UserRoleCommon)

	// Act
	user, err := cognitoService.GetUserInCognito(username)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if user == nil {
		t.Error("Expected user to be returned, but got nil")
	}

	if user != nil && *user.Username != username {
		t.Errorf("Expected username to be %s, but got %s", username, *user.Username)
	}

	if mockRepo.GetCallCount != 1 {
		t.Errorf("Expected GetUserFromCognito to be called once, but was called %d times", mockRepo.GetCallCount)
	}
}

func TestCognitoService_WithRepositoryErrors(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockRepo)

	// Configurar o mock para falhar
	mockRepo.ShouldFailGet = true

	user := &models.User{
		Name:  "Test User",
		Email: "test@example.com",
		Role:  models.UserRoleCommon,
	}

	// Act & Assert - UpdateUserInCognito
	err := cognitoService.UpdateUserInCognito(user)
	if err == nil {
		t.Error("Expected error when repository fails, but got nil")
	}

	// Act & Assert - DeleteUserFromCognito
	err = cognitoService.DeleteUserFromCognito("test@example.com")
	if err == nil {
		t.Error("Expected error when repository fails, but got nil")
	}

	// Act & Assert - GetUserInCognito
	_, err = cognitoService.GetUserInCognito("test@example.com")
	if err == nil {
		t.Error("Expected error when repository fails, but got nil")
	}
}
