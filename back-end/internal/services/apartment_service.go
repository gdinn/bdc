package services

import (
	"bdc/internal/domain"
	"bdc/internal/interfaces"
	"bdc/internal/models"
	"bdc/internal/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ApartmentService struct {
	apartmentRepo interfaces.ApartmentRepositoryInterface
}

const (
	ErrCreatingApartment = "failed to create apartment"
)

func NewApartmentService(apartmentRepo interfaces.ApartmentRepositoryInterface) *ApartmentService {
	return &ApartmentService{
		apartmentRepo: apartmentRepo,
	}
}

func (s *ApartmentService) CreateApartment(apartment *models.Apartment, claims *models.UserClaims) (*models.Apartment, error) {
	if err := utils.ValidateManagerRole(claims); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCreatingApartment, err)
	}

	if err := s.validateApartmentIsNew(apartment); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCreatingApartment, err)
	}

	createdApartment, err := s.apartmentRepo.Create(apartment)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCreatingApartment, err)
	}

	return createdApartment, nil
}

func (s *ApartmentService) validateApartmentIsNew(apartment *models.Apartment) error {
	apartmentExists, err := s.apartmentRepo.IsApartmentExists(apartment)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error checking apartment existence: %w", err)
	}

	if apartmentExists {
		return fmt.Errorf("apartment.Number is %s and apartment.Building is %s: %w", apartment.Number, apartment.Building, domain.ErrApartmentAlreadyExists)
	}

	return nil
}
