package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"bdc/internal/domain"
	"bdc/internal/middleware"
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

type CreateUserResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *models.User `json:"data,omitempty"`
	Error   string       `json:"error,omitempty"`
}

// CreateUser handles POST /api/v1/users (with mandatory auth)
// Will create a user in BDC to the corresponding user in Cognito
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	userClaims, ok := middleware.GetUserClaimsFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Authentication context not found", fmt.Errorf("missing authentication context"))
		return
	}

	// Parse request body
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}

	// Validate request structure (tags)
	if err := h.validator.Struct(req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	ctx := h.createUserContext(userClaims, &req)

	createdUser, err := h.userService.CreateUserWithContext(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
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

// createUserContext factory method para criar o contexto completo
func (h *UserHandler) createUserContext(claims *middleware.UserClaims, req *domain.CreateUserRequest) *domain.CreateUserContext {
	// Converter request HTTP para domain data
	userData := domain.CreateUserData(req)

	ctx := domain.NewCreateUserContext(claims, userData)

	return ctx
}
