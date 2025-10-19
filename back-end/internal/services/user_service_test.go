package services

import (
	"testing"
	"time"

	"bdc/internal/domain"
	"bdc/internal/models"
	"bdc/internal/repositories"
)

func TestUserService_CreateUserWithContext_Success(t *testing.T) {
	// Arrange
	mockUserRepo := repositories.NewMockUserRepository()
	mockCognitoRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockCognitoRepo)
	userService := NewUserService(mockUserRepo, cognitoService)

	// Adicionar usuário no Cognito mock
	mockCognitoRepo.AddTestUser("test@example.com", "Test User", models.UserRoleCommon)

	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "test@example.com",
		Role:     "COMMON",
	}

	user := &models.User{
		Type:     models.UserTypeResident,
		AgeGroup: models.UserAgeGroupAdult,
		Role:     models.UserRoleCommon,
		Phone:    "11999999999",
	}

	ctx := domain.NewCreateUserContext(claims, user)

	// Act
	createdUser, err := userService.CreateUserWithContext(ctx)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if createdUser == nil {
		t.Error("Expected created user to be returned, but got nil")
	}

	if createdUser.Email != "test@example.com" {
		t.Errorf("Expected email to be test@example.com, but got %s", createdUser.Email)
	}

	if createdUser.Name != "Test User" {
		t.Errorf("Expected name to be Test User, but got %s", createdUser.Name)
	}

	if mockUserRepo.CreateCallCount != 1 {
		t.Errorf("Expected Create to be called once, but was called %d times", mockUserRepo.CreateCallCount)
	}

	if mockCognitoRepo.GetCallCount != 1 {
		t.Errorf("Expected GetUserFromCognito to be called once, but was called %d times", mockCognitoRepo.GetCallCount)
	}
}

func TestUserService_CreateUserWithContext_ExternalChildError(t *testing.T) {
	// Arrange
	mockUserRepo := repositories.NewMockUserRepository()
	mockCognitoRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockCognitoRepo)
	userService := NewUserService(mockUserRepo, cognitoService)

	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "test@example.com",
		Role:     "COMMON",
	}

	user := &models.User{
		Type:     models.UserTypeExternal,  // EXTERNAL
		AgeGroup: models.UserAgeGroupChild, // CHILD - Invalid combination
		Role:     models.UserRoleCommon,
	}

	ctx := domain.NewCreateUserContext(claims, user)

	// Act
	_, err := userService.CreateUserWithContext(ctx)

	// Assert
	if err == nil {
		t.Error("Expected error for external child user, but got nil")
	}

	if err != nil && err.Error() != "user.Type is EXTERNAL and user.AgeGroup is CHILD: external users cannot be child" {
		t.Errorf("Expected specific error message, but got: %v", err)
	}

	if mockUserRepo.CreateCallCount != 0 {
		t.Errorf("Expected Create not to be called, but was called %d times", mockUserRepo.CreateCallCount)
	}
}

func TestUserService_CreateUserWithContext_ChildManagerError(t *testing.T) {
	// Arrange
	mockUserRepo := repositories.NewMockUserRepository()
	mockCognitoRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockCognitoRepo)
	userService := NewUserService(mockUserRepo, cognitoService)

	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "test@example.com",
		Role:     "MANAGER",
	}

	user := &models.User{
		Type:     models.UserTypeResident,
		AgeGroup: models.UserAgeGroupChild, // CHILD
		Role:     models.UserRoleManager,   // MANAGER - Invalid combination
	}

	ctx := domain.NewCreateUserContext(claims, user)

	// Act
	_, err := userService.CreateUserWithContext(ctx)

	// Assert
	if err == nil {
		t.Error("Expected error for child manager user, but got nil")
	}

	if err != nil && err.Error() != "user.AgeGroup is CHILD and user.Role is MANAGER: children cannot be managers or advisors" {
		t.Errorf("Expected specific error message, but got: %v", err)
	}
}

func TestUserService_CreateUserWithContext_EmailAlreadyExists(t *testing.T) {
	// Arrange
	mockUserRepo := repositories.NewMockUserRepository()
	mockCognitoRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockCognitoRepo)
	userService := NewUserService(mockUserRepo, cognitoService)

	// Adicionar usuário existente no mock repo
	mockUserRepo.AddTestUser("test@example.com", "Existing User", models.UserTypeResident, models.UserAgeGroupAdult, models.UserRoleCommon)

	// Adicionar usuário no Cognito mock
	mockCognitoRepo.AddTestUser("test@example.com", "Test User", models.UserRoleCommon)

	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "test@example.com",
		Role:     "COMMON",
	}

	user := &models.User{
		Type:     models.UserTypeResident,
		AgeGroup: models.UserAgeGroupAdult,
		Role:     models.UserRoleCommon,
	}

	ctx := domain.NewCreateUserContext(claims, user)

	// Act
	_, err := userService.CreateUserWithContext(ctx)

	// Assert
	if err == nil {
		t.Error("Expected error for existing email, but got nil")
	}

	if mockUserRepo.CreateCallCount != 0 {
		t.Errorf("Expected Create not to be called, but was called %d times", mockUserRepo.CreateCallCount)
	}

	if mockUserRepo.IsEmailExistsCallCount != 1 {
		t.Errorf("Expected IsEmailExists to be called once, but was called %d times", mockUserRepo.IsEmailExistsCallCount)
	}
}

func TestUserService_CreateUserWithContext_CognitoError(t *testing.T) {
	// Arrange
	mockUserRepo := repositories.NewMockUserRepository()
	mockCognitoRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockCognitoRepo)
	userService := NewUserService(mockUserRepo, cognitoService)

	// Configurar o mock do Cognito para falhar
	mockCognitoRepo.ShouldFailGet = true

	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "test@example.com",
		Role:     "COMMON",
	}

	user := &models.User{
		Type:     models.UserTypeResident,
		AgeGroup: models.UserAgeGroupAdult,
		Role:     models.UserRoleCommon,
	}

	ctx := domain.NewCreateUserContext(claims, user)

	// Act
	_, err := userService.CreateUserWithContext(ctx)

	// Assert
	if err == nil {
		t.Error("Expected error when Cognito fails, but got nil")
	}

	if mockUserRepo.CreateCallCount != 0 {
		t.Errorf("Expected Create not to be called, but was called %d times", mockUserRepo.CreateCallCount)
	}
}

func TestUserService_CreateUserWithContext_RepositoryError(t *testing.T) {
	// Arrange
	mockUserRepo := repositories.NewMockUserRepository()
	mockCognitoRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockCognitoRepo)
	userService := NewUserService(mockUserRepo, cognitoService)

	// Adicionar usuário no Cognito mock
	mockCognitoRepo.AddTestUser("test@example.com", "Test User", models.UserRoleCommon)

	// Configurar o mock do repositório para falhar na criação
	mockUserRepo.ShouldFailCreate = true

	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "test@example.com",
		Role:     "COMMON",
	}

	user := &models.User{
		Type:     models.UserTypeResident,
		AgeGroup: models.UserAgeGroupAdult,
		Role:     models.UserRoleCommon,
	}

	ctx := domain.NewCreateUserContext(claims, user)

	// Act
	_, err := userService.CreateUserWithContext(ctx)

	// Assert
	if err == nil {
		t.Error("Expected error when repository fails, but got nil")
	}

	if mockUserRepo.CreateCallCount != 1 {
		t.Errorf("Expected Create to be called once, but was called %d times", mockUserRepo.CreateCallCount)
	}
}

func TestUserService_ValidateBusinessRules_InvalidUserType(t *testing.T) {
	// Arrange
	mockUserRepo := repositories.NewMockUserRepository()
	mockCognitoRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockCognitoRepo)
	userService := NewUserService(mockUserRepo, cognitoService)

	user := &models.User{
		Type:     "INVALID_TYPE", // Tipo inválido
		AgeGroup: models.UserAgeGroupAdult,
		Role:     models.UserRoleCommon,
	}

	// Act
	err := userService.validateBusinessRules(user)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid user type, but got nil")
	}
}

func TestUserService_ValidateBusinessRules_InvalidAgeGroup(t *testing.T) {
	// Arrange
	mockUserRepo := repositories.NewMockUserRepository()
	mockCognitoRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockCognitoRepo)
	userService := NewUserService(mockUserRepo, cognitoService)

	user := &models.User{
		Type:     models.UserTypeResident,
		AgeGroup: "INVALID_AGE_GROUP", // Grupo de idade inválido
		Role:     models.UserRoleCommon,
	}

	// Act
	err := userService.validateBusinessRules(user)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid age group, but got nil")
	}
}

func TestUserService_ValidateBusinessRules_ValidCombinations(t *testing.T) {
	// Arrange
	mockUserRepo := repositories.NewMockUserRepository()
	mockCognitoRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockCognitoRepo)
	userService := NewUserService(mockUserRepo, cognitoService)

	validCombinations := []struct {
		userType models.UserType
		ageGroup models.UserAgeGroup
		role     models.UserRole
	}{
		{models.UserTypeResident, models.UserAgeGroupAdult, models.UserRoleCommon},
		{models.UserTypeResident, models.UserAgeGroupAdult, models.UserRoleManager},
		{models.UserTypeResident, models.UserAgeGroupAdult, models.UserRoleAdvisor},
		{models.UserTypeResident, models.UserAgeGroupChild, models.UserRoleCommon},
		{models.UserTypeExternal, models.UserAgeGroupAdult, models.UserRoleCommon},
		{models.UserTypeExternal, models.UserAgeGroupAdult, models.UserRoleManager},
		{models.UserTypeExternal, models.UserAgeGroupAdult, models.UserRoleAdvisor},
	}

	for _, combo := range validCombinations {
		user := &models.User{
			Type:     combo.userType,
			AgeGroup: combo.ageGroup,
			Role:     combo.role,
		}

		// Act
		err := userService.validateBusinessRules(user)

		// Assert
		if err != nil {
			t.Errorf("Expected no error for valid combination %+v, but got: %v", combo, err)
		}
	}
}

func TestUserService_ValidateUserIsNew_RepositoryError(t *testing.T) {
	// Arrange
	mockUserRepo := repositories.NewMockUserRepository()
	mockCognitoRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockCognitoRepo)
	userService := NewUserService(mockUserRepo, cognitoService)

	// Configurar o mock para falhar na verificação de email
	mockUserRepo.ShouldFailIsEmailExists = true

	user := &models.User{
		Email: "test@example.com",
	}

	// Act
	err := userService.validateUserIsNew(user)

	// Assert
	if err == nil {
		t.Error("Expected error when repository fails, but got nil")
	}
}

func TestUserService_CreateUserWithContext_WithBirthDate(t *testing.T) {
	// Arrange
	mockUserRepo := repositories.NewMockUserRepository()
	mockCognitoRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockCognitoRepo)
	userService := NewUserService(mockUserRepo, cognitoService)

	// Adicionar usuário no Cognito mock
	mockCognitoRepo.AddTestUser("test@example.com", "Test User", models.UserRoleCommon)

	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "test@example.com",
		Role:     "COMMON",
	}

	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	user := &models.User{
		Type:      models.UserTypeResident,
		AgeGroup:  models.UserAgeGroupAdult,
		Role:      models.UserRoleCommon,
		Phone:     "11999999999",
		BirthDate: &birthDate,
	}

	ctx := domain.NewCreateUserContext(claims, user)

	// Act
	createdUser, err := userService.CreateUserWithContext(ctx)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if createdUser == nil {
		t.Error("Expected created user to be returned, but got nil")
	}

	if createdUser.BirthDate == nil {
		t.Error("Expected birth date to be preserved, but got nil")
	}

	if !createdUser.BirthDate.Equal(birthDate) {
		t.Errorf("Expected birth date to be %v, but got %v", birthDate, *createdUser.BirthDate)
	}
}

func TestUserService_Setup(t *testing.T) {
	// Test que verifica se o service é criado corretamente
	mockUserRepo := repositories.NewMockUserRepository()
	mockCognitoRepo := repositories.NewMockCognitoRepository()
	cognitoService := NewCognitoService(mockCognitoRepo)
	userService := NewUserService(mockUserRepo, cognitoService)

	if userService == nil {
		t.Error("Expected UserService to be created, but got nil")
	}

	if userService.userRepo == nil {
		t.Error("Expected userRepo to be set, but got nil")
	}

	if userService.cognitoService == nil {
		t.Error("Expected cognitoService to be set, but got nil")
	}
}
