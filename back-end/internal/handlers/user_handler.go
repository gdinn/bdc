package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"bdc/internal/models"
	"bdc/internal/services"
	"bdc/internal/utils"
)

type UserHandler struct {
	userService *services.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.New(),
	}
}

// CreateUserRequest representa a estrutura de dados para criação de usuário
type CreateUserRequest struct {
	Name      string              `json:"name" validate:"required,min=2,max=255"`
	Email     string              `json:"email" validate:"required,email"`
	Phone     string              `json:"phone" validate:"omitempty,min=10,max=20"`
	BirthDate *time.Time          `json:"birth_date,omitempty"`
	Type      models.UserType     `json:"type" validate:"required"`
	AgeGroup  models.UserAgeGroup `json:"age_group" validate:"required"`
	IsManager bool                `json:"is_manager,omitempty"`
	IsAdvisor bool                `json:"is_advisor,omitempty"`
}

// CreateUserResponse representa a resposta da criação de usuário
type CreateUserResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *models.User `json:"data,omitempty"`
	Error   string       `json:"error,omitempty"`
}

// CreateUser handles POST /api/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Validate user type and age group enums
	if !isValidUserType(req.Type) {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid user type", fmt.Errorf("type must be RESIDENT or EXTERNAL"))
		return
	}

	if !isValidUserAgeGroup(req.AgeGroup) {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid age group", fmt.Errorf("age_group must be ADULT or CHILD"))
		return
	}

	// Create user model
	user := &models.User{
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		BirthDate: req.BirthDate,
		Type:      req.Type,
		AgeGroup:  req.AgeGroup,
		IsManager: req.IsManager,
		IsAdvisor: req.IsAdvisor,
	}

	// Call service to create user
	createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		if utils.IsEmailAlreadyExistsError(err) {
			utils.SendErrorResponse(w, http.StatusConflict, "Email already exists", err)
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	// Send success response
	response := CreateUserResponse{
		Success: true,
		Message: "User created successfully",
		Data:    createdUser,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Helper functions for validation
func isValidUserType(userType models.UserType) bool {
	return userType == models.UserTypeResident || userType == models.UserTypeExternal
}

func isValidUserAgeGroup(userAgeGroup models.UserAgeGroup) bool {
	return userAgeGroup == models.UserAgeGroupAdult || userAgeGroup == models.UserAgeGroupChild
}
