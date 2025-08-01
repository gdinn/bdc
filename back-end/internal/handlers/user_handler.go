package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"bdc/internal/middleware"
	"bdc/internal/models"
	"bdc/internal/services"
	"bdc/internal/utils"
)

type UserHandler struct {
	userService    *services.UserService
	validator      *validator.Validate
	authMiddleware *middleware.AuthMiddleware
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService:    userService,
		validator:      validator.New(),
		authMiddleware: middleware.NewAuthMiddleware(),
	}
}

type CreateUserRequest struct {
	Name      string              `json:"name" validate:"required,min=2,max=255"`
	Email     string              `json:"email" validate:"required,email"`
	Phone     string              `json:"phone" validate:"omitempty,min=10,max=20"`
	BirthDate *time.Time          `json:"birth_date,omitempty"`
	Type      models.UserType     `json:"type" validate:"required"`
	AgeGroup  models.UserAgeGroup `json:"age_group" validate:"required"`
}

type CreateUserResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *models.User `json:"data,omitempty"`
	Error   string       `json:"error,omitempty"`
}

// CreateUser handles POST /api/users (with mandatory auth)
// After registering in cognito with only email and pwd, users must provide extra info here
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Authorization header required", fmt.Errorf("missing authorization header"))
		return
	}

	userClaims, err := h.authMiddleware.ValidateToken(authHeader)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token", err)
		return
	}

	// Parse request body
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}

	// Validate request structure
	if err := h.validator.Struct(req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	if !isValidUserType(req.Type) {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid user type", fmt.Errorf("type must be RESIDENT or EXTERNAL"))
		return
	}

	if !isValidUserAgeGroup(req.AgeGroup) {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid age group", fmt.Errorf("age_group must be ADULT or CHILD"))
		return
	}

	if err := h.validateBusinessRulesForUserCreation(userClaims, &req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Business rule validation failed", err)
		return
	}

	user := &models.User{
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		BirthDate: req.BirthDate,
		Type:      req.Type,
		AgeGroup:  req.AgeGroup,
		Role:      models.UserRoleCommon,
	}

	createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		if utils.IsEmailAlreadyExistsError(err) {
			utils.SendErrorResponse(w, http.StatusConflict, "Email already exists", err)
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	response := CreateUserResponse{
		Success: true,
		Message: fmt.Sprintf("User created successfully by %s", userClaims.Email),
		Data:    createdUser,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) validateBusinessRulesForUserCreation(claims *middleware.UserClaims, req *CreateUserRequest) error {

	// Validate if the received email matches the claims
	if req.Email != claims.Email {
		return fmt.Errorf("cannot create user with same email as authenticated user")
	}

	return nil
}

func isValidUserType(userType models.UserType) bool {
	return userType == models.UserTypeResident || userType == models.UserTypeExternal
}

func isValidUserAgeGroup(userAgeGroup models.UserAgeGroup) bool {
	return userAgeGroup == models.UserAgeGroupAdult || userAgeGroup == models.UserAgeGroupChild
}

// CreateUserWithAuth is a wrapper that applies authentication middleware
func (h *UserHandler) CreateUserWithAuth() http.HandlerFunc {
	return h.authMiddleware.RequireAuth(h.CreateUser)
}
