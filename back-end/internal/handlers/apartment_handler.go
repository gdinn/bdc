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
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ApartmentHandler struct {
	apartmentService *services.ApartmentService
	validator        *validator.Validate
}

func NewApartmentHandler(apartmentService *services.ApartmentService) *ApartmentHandler {
	return &ApartmentHandler{
		apartmentService: apartmentService,
		validator:        validator.New(),
	}
}

type CreateApartmentResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *models.Apartment `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
}

// CreateApartment handles POST /api/v1/apartments (with mandatory auth)
// Will create a apartment in BDC
func (h *ApartmentHandler) CreateApartment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims, ok := middleware.GetUserClaimsFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Authentication context not found", fmt.Errorf("missing authentication context"))
		return
	}

	var req domain.CreateApartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	apartmentData := domain.CreateApartmentData(&req)

	createdApartment, err := h.apartmentService.CreateApartment(apartmentData, userClaims)
	if err != nil {
		if errors.Is(err, domain.ErrApartmentAlreadyExists) {
			utils.SendErrorResponse(w, http.StatusConflict, "Apartment already exists", err)
			return
		} else if errors.Is(err, domain.ErrManagerRoleRequired) {
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Insuficient permissions", err)
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create apartment", err)
		return
	}

	response := domain.CreateApartmentResponse{
		Success: true,
		Message: fmt.Sprintf("Apartment created successfully by %s", userClaims.Email),
		Data:    createdApartment,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *ApartmentHandler) GetApartmentByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims, ok := middleware.GetUserClaimsFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Authentication context not found", fmt.Errorf("missing authentication context"))
		return
	}

	vars := mux.Vars(r)
	apartmentIDStr := vars["id"]
	apartmentID, err := uuid.Parse(apartmentIDStr)

	apartment, err := h.apartmentService.GetApartment(apartmentID, userClaims)
	if err != nil {
		if errors.Is(err, domain.ErrManagerRoleRequired) {
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Insuficient permissions", err)
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get apartment", err)
		return
	}

	response := domain.GetApartmentResponse{
		Success: true,
		Message: fmt.Sprintf("Apartment obtained successfully by %s", userClaims.Email),
		Data:    apartment,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
