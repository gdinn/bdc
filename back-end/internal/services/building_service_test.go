package services

import (
	"bdc/internal/domain"
	"bdc/internal/models"
	"bdc/internal/repositories"
	"errors"
	"testing"
)

func TestBuildingService_CreateBuilding_Success(t *testing.T) {
	// Arrange
	mockBuildingRepo := repositories.NewMockBuildingRepository()
	buildingService := NewBuildingService(mockBuildingRepo)

	building := &models.Building{
		Name: "North",
	}

	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "test@example.com",
		Role:     "MANAGER",
	}

	// Act
	createdBuilding, err := buildingService.CreateBuilding(building, claims)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if createdBuilding == nil {
		t.Error("Expected created building to be returned, but got nil")
	}

	if createdBuilding.Name != "North" {
		t.Errorf("Expected name to be North, but got %s", createdBuilding.Name)
	}

	if mockBuildingRepo.CreateCallCount != 1 {
		t.Errorf("Expected Create to be called once, but was called %d times", mockBuildingRepo.CreateCallCount)
	}
}

func TestBuildingService_CreateBuilding_AlreadyExistsError(t *testing.T) {
	// Arrange
	mockBuildingRepo := repositories.NewMockBuildingRepository()
	mockBuildingRepo.AddTestBuilding("North")

	buildingService := NewBuildingService(mockBuildingRepo)
	building := &models.Building{
		Name: "North",
	}

	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "test@example.com",
		Role:     "MANAGER",
	}

	// Act
	_, err := buildingService.CreateBuilding(building, claims)

	// Assert
	if err == nil {
		t.Error("Expected error for existing building, but got nil")
	}

	if mockBuildingRepo.CreateCallCount != 0 {
		t.Errorf("Expected Create not to be called, but was called %d times", mockBuildingRepo.CreateCallCount)
	}

	if mockBuildingRepo.IsBuildingExistsCallCount != 1 {
		t.Errorf("Expected IsBuildingExists to be called once, but was called %d times", mockBuildingRepo.IsBuildingExistsCallCount)
	}
}

func TestBuildingService_CreateBuilding_InvalidClaims(t *testing.T) {
	// Arrange
	mockBuildingRepo := repositories.NewMockBuildingRepository()

	buildingService := NewBuildingService(mockBuildingRepo)
	building := &models.Building{
		Name: "North",
	}

	claims := &models.UserClaims{
		Email: "test@example.com",
	}

	// Act
	_, err := buildingService.CreateBuilding(building, claims)
	// Assert
	if err == nil {
		t.Error("Expected error to be ErrManagerRoleRequired, but got nil")
	}

	if !errors.Is(err, domain.ErrManagerRoleRequired) {
		t.Errorf("Expected error to be ErrManagerRoleRequired, but got %s", err)
	}

	if mockBuildingRepo.CreateCallCount != 0 {
		t.Errorf("Expected Create not to be called, but was called %d times", mockBuildingRepo.CreateCallCount)
	}

	if mockBuildingRepo.IsBuildingExistsCallCount != 0 {
		t.Errorf("Expected IsBuildingExists not to be called, but was called %d times", mockBuildingRepo.IsBuildingExistsCallCount)
	}
}

func TestBuildingService_CreateBuilding_AdvisorRole(t *testing.T) {
	// Arrange
	mockBuildingRepo := repositories.NewMockBuildingRepository()

	buildingService := NewBuildingService(mockBuildingRepo)
	building := &models.Building{
		Name: "North",
	}

	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "test@example.com",
		Role:     "ADVISOR",
	}

	// Act
	_, err := buildingService.CreateBuilding(building, claims)
	// Assert
	if err == nil {
		t.Error("Expected error to be ErrManagerRoleRequired, but got nil")
	}

	if !errors.Is(err, domain.ErrManagerRoleRequired) {
		t.Errorf("Expected error to be ErrManagerRoleRequired, but got %s", err)
	}

	if mockBuildingRepo.CreateCallCount != 0 {
		t.Errorf("Expected Create not to be called, but was called %d times", mockBuildingRepo.CreateCallCount)
	}

	if mockBuildingRepo.IsBuildingExistsCallCount != 0 {
		t.Errorf("Expected IsBuildingExists not to be called, but was called %d times", mockBuildingRepo.IsBuildingExistsCallCount)
	}
}

func TestBuildingService_CreateBuilding_CommonRole(t *testing.T) {
	// Arrange
	mockBuildingRepo := repositories.NewMockBuildingRepository()

	buildingService := NewBuildingService(mockBuildingRepo)
	building := &models.Building{
		Name: "North",
	}

	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "test@example.com",
		Role:     "COMMON",
	}

	// Act
	_, err := buildingService.CreateBuilding(building, claims)
	// Assert
	if err == nil {
		t.Error("Expected error to be ErrManagerRoleRequired, but got nil")
	}

	if !errors.Is(err, domain.ErrManagerRoleRequired) {
		t.Errorf("Expected error to be ErrManagerRoleRequired, but got %s", err)
	}

	if mockBuildingRepo.CreateCallCount != 0 {
		t.Errorf("Expected Create not to be called, but was called %d times", mockBuildingRepo.CreateCallCount)
	}

	if mockBuildingRepo.IsBuildingExistsCallCount != 0 {
		t.Errorf("Expected IsBuildingExists not to be called, but was called %d times", mockBuildingRepo.IsBuildingExistsCallCount)
	}
}

func TestBuildingService_Setup(t *testing.T) {
	mockBuildingRepo := repositories.NewMockBuildingRepository()
	buildingService := NewBuildingService(mockBuildingRepo)

	if buildingService == nil {
		t.Error("Expected BuildingService to be created, but got nil")
	}
}
