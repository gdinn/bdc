package handlers

import (
	"bdc/internal/domain"
	"bdc/internal/middleware"
	"bdc/internal/models"
	"bdc/internal/services"
	"bdc/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type BuildingHandler struct {
	buildingService *services.BuildingService
	validator       *validator.Validate
}

func NewBuildingHandler(buildingService *services.BuildingService) *BuildingHandler {
	return &BuildingHandler{
		buildingService: buildingService,
		validator:       validator.New(),
	}
}

type CreateBuildingResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *models.Building `json:"data,omitempty"`
	Error   string           `json:"error,omitempty"`
}

func (h *BuildingHandler) CreateBuilding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims, ok := middleware.GetUserClaimsFromContext(r.Context())

	if !ok {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Authentication context not found", fmt.Errorf("missing authentication context"))
		return
	}

	var req domain.CreateBuildingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	buidlingData := domain.CreateBuildingData(&req)

	createdBuilding, err := h.buildingService.CreateBuilding(buidlingData, userClaims)

	if err != nil {
		if errors.Is(err, domain.ErrBuildingAlreadyExists) {
			utils.SendErrorResponse(w, http.StatusConflict, "Building already exists", err)
			return
		} else if errors.Is(err, domain.ErrManagerRoleRequired) {
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Insuficient permissions", err)
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create building", err)
		return
	}

	response := domain.CreateBuildingResponse{
		Success: true,
		Message: fmt.Sprintf("building created successfully by %s", userClaims.Email),
		Data:    createdBuilding,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
